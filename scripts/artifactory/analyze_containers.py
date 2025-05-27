import requests
import json
import os
from datetime import datetime, timedelta
from dateutil import parser
import sys
from urllib.parse import urljoin
import base64

class ArtifactoryAnalyzer:
    def __init__(self, base_url, username, password):
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        self.session.auth = (username, password)
        self.session.headers.update({
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        })
        
    def get_repositories(self):
        """Get all Docker repositories"""
        try:
            url = f"{self.base_url}/artifactory/api/repositories"
            response = self.session.get(url)
            response.raise_for_status()
            
            # Filter for Docker repositories
            repos = response.json()
            docker_repos = [repo for repo in repos if repo.get('packageType') == 'Docker']
            return docker_repos
        except Exception as e:
            print(f"Error getting repositories: {e}")
            return []
    
    def get_docker_images(self, repo_name):
        """Get all Docker images in a repository"""
        try:
            url = f"{self.base_url}/artifactory/api/docker/{repo_name}/v2/_catalog"
            response = self.session.get(url)
            response.raise_for_status()
            
            data = response.json()
            return data.get('repositories', [])
        except Exception as e:
            print(f"Error getting images from {repo_name}: {e}")
            return []
    
    def get_image_tags(self, repo_name, image_name):
        """Get all tags for a Docker image"""
        try:
            url = f"{self.base_url}/artifactory/api/docker/{repo_name}/v2/{image_name}/tags/list"
            response = self.session.get(url)
            response.raise_for_status()
            
            data = response.json()
            return data.get('tags', [])
        except Exception as e:
            print(f"Error getting tags for {image_name}: {e}")
            return []
    
    def get_artifact_info(self, repo_name, image_name, tag):
        """Get detailed information about an artifact"""
        try:
            # Try multiple path structures for different Artifactory setups
            paths_to_try = [
                f"{image_name}/{tag}/manifest.json",
                f"{image_name}/{tag}",
                f"{image_name}/{tag}/latest",
                f"docker/{image_name}/{tag}",
            ]
            
            for artifact_path in paths_to_try:
                url = f"{self.base_url}/artifactory/api/storage/{repo_name}/{artifact_path}"
                response = self.session.get(url)
                
                if response.status_code == 200:
                    artifact_info = response.json()
                    # Debug: print the raw response to understand the structure
                    if not artifact_info.get('created'):
                        print(f"    DEBUG - No 'created' field. Available fields: {list(artifact_info.keys())}")
                    return artifact_info
            
            # If no path worked, try getting folder info
            folder_url = f"{self.base_url}/artifactory/api/storage/{repo_name}/{image_name}"
            response = self.session.get(folder_url)
            if response.status_code == 200:
                folder_info = response.json()
                # Look for children that match our tag
                children = folder_info.get('children', [])
                for child in children:
                    if child.get('uri', '').strip('/') == tag:
                        return {
                            'created': child.get('lastModified', folder_info.get('lastModified')),
                            'lastModified': child.get('lastModified'),
                            'size': child.get('size', 0)
                        }
            
            print(f"    Could not find artifact info for {image_name}:{tag} (tried {len(paths_to_try)} paths)")
            return None
                
        except Exception as e:
            print(f"    Error getting artifact info for {image_name}:{tag}: {e}")
            return None
    
    def load_values_json(self, file_path):
        """Load and parse values.json to extract referenced images"""
        referenced_images = set()
        try:
            if os.path.exists(file_path):
                with open(file_path, 'r') as f:
                    values = json.load(f)
                    
                # Extract image references (adjust based on your values.json structure)
                def extract_images(obj, path=""):
                    if isinstance(obj, dict):
                        for key, value in obj.items():
                            current_path = f"{path}.{key}" if path else key
                            if key in ['image', 'repository', 'tag'] and isinstance(value, str):
                                if ':' in value:
                                    referenced_images.add(value)
                                elif key == 'repository' and isinstance(obj.get('tag'), str):
                                    referenced_images.add(f"{value}:{obj['tag']}")
                            extract_images(value, current_path)
                    elif isinstance(obj, list):
                        for i, item in enumerate(obj):
                            extract_images(item, f"{path}[{i}]")
                
                extract_images(values)
                print(f"Found {len(referenced_images)} referenced images in values.json")
                
            else:
                print(f"Warning: {file_path} not found")
                
        except Exception as e:
            print(f"Error loading values.json: {e}")
        
        return referenced_images
    
    def analyze_containers(self, values_json_path, months_threshold=3):
        """Analyze containers and identify candidates for cleanup"""
        print("Starting Artifactory container analysis...")
        print("=" * 60)
        
        # Load referenced images from values.json
        referenced_images = self.load_values_json(values_json_path)
        
        # Calculate cutoff date
        cutoff_date = datetime.now() - timedelta(days=months_threshold * 30)
        print(f"Analyzing containers older than: {cutoff_date.strftime('%Y-%m-%d')}")
        print(f"Referenced images in values.json: {len(referenced_images)}")
        print("=" * 60)
        
        cleanup_candidates = []
        total_analyzed = 0
        
        # Get all Docker repositories
        repositories = self.get_repositories()
        print(f"Found {len(repositories)} Docker repositories")
        
        for repo in repositories:
            repo_name = repo['key']
            print(f"\nAnalyzing repository: {repo_name}")
            print("-" * 40)
            
            # Get all images in repository
            images = self.get_docker_images(repo_name)
            
            for image_name in images:
                print(f"  Image: {image_name}")
                
                # Get all tags for this image
                tags = self.get_image_tags(repo_name, image_name)
                
                for tag in tags:
                    total_analyzed += 1
                    full_image_name = f"{image_name}:{tag}"
                    
                    # Get artifact information
                    artifact_info = self.get_artifact_info(repo_name, image_name, tag)
                    
                    if artifact_info:
                        # Try multiple date fields and formats
                        date_fields = ['created', 'lastModified', 'lastUpdated']
                        created_date = None
                        date_source = None
                        
                        print(f"    DEBUG - Artifact info keys: {list(artifact_info.keys())}")
                        
                        for field in date_fields:
                            date_str = artifact_info.get(field)
                            if date_str:
                                print(f"    DEBUG - Trying to parse {field}='{date_str}'")
                                try:
                                    # Try different date parsing approaches
                                    if isinstance(date_str, str):
                                        # Handle ISO format with timezone (most common)
                                        created_date = parser.parse(date_str)
                                    elif isinstance(date_str, (int, float)):
                                        # Handle timestamp as number
                                        created_date = datetime.fromtimestamp(date_str / 1000 if date_str > 1e10 else date_str)
                                    
                                    if created_date:
                                        date_source = field
                                        print(f"    DEBUG - Successfully parsed {field} to {created_date}")
                                        break
                                        
                                except Exception as date_error:
                                    print(f"    DEBUG - Failed to parse {field}='{date_str}': {date_error}")
                                    continue
                        
                        if created_date:
                            is_old = created_date < cutoff_date
                            is_referenced = any(ref_img in full_image_name or full_image_name in ref_img 
                                                for ref_img in referenced_images)
                            
                            status = []
                            if is_old:
                                status.append("OLD")
                            if is_referenced:
                                status.append("REFERENCED")
                            if not status:
                                status.append("RECENT")
                            
                            print(f"    Tag: {tag:20} | Date: {created_date.strftime('%Y-%m-%d')} ({date_source}) | Status: {', '.join(status)}")
                            
                            # Add to cleanup candidates if old and not referenced
                            if is_old and not is_referenced:
                                cleanup_candidates.append({
                                    'repository': repo_name,
                                    'image': image_name,
                                    'tag': tag,
                                    'full_name': full_image_name,
                                    'created_date': created_date.isoformat(),
                                    'age_days': (datetime.now() - created_date).days,
                                    'date_source': date_source
                                })
                        else:
                            print(f"    Tag: {tag:20} | Date: NO PARSEABLE DATE | Available fields: {list(artifact_info.keys())}")
                            # Show date field values for debugging
                            for field in date_fields:
                                value = artifact_info.get(field, 'NOT FOUND')
                                print(f"    DEBUG - {field}: {value}")
                    else:
                        print(f"    Tag: {tag:20} | Date: INFO NOT AVAILABLE | Status: UNKNOWN")
        
        # Print summary
        print("\n" + "=" * 60)
        print("ANALYSIS SUMMARY")
        print("=" * 60)
        print(f"Total containers analyzed: {total_analyzed}")
        print(f"Cleanup candidates (old + not referenced): {len(cleanup_candidates)}")
        
        if cleanup_candidates:
            print(f"\nCLEANUP CANDIDATES:")
            print("-" * 60)
            for candidate in sorted(cleanup_candidates, key=lambda x: x['age_days'], reverse=True):
                print(f"Repository: {candidate['repository']}")
                print(f"Image: {candidate['full_name']}")
                print(f"Created: {candidate['created_date'][:10]} ({candidate['age_days']} days ago)")
                print("-" * 40)
            
            # Export to JSON for potential automation
            with open('cleanup_candidates.json', 'w') as f:
                json.dump(cleanup_candidates, f, indent=2)
            print(f"\nCleanup candidates exported to: cleanup_candidates.json")
        
        return cleanup_candidates

def main():
    # Get environment variables
    artifactory_url = os.getenv('ARTIFACTORY_URL')
    username = os.getenv('ARTIFACTORY_USERNAME')
    password = os.getenv('ARTIFACTORY_PASSWORD')
    values_json_path = os.getenv('VALUES_JSON_PATH', 'values.json')
    
    if not all([artifactory_url, username, password]):
        print("Error: Missing required environment variables")
        print("Required: ARTIFACTORY_URL, ARTIFACTORY_USERNAME, ARTIFACTORY_PASSWORD")
        sys.exit(1)
    
    analyzer = ArtifactoryAnalyzer(artifactory_url, username, password)
    cleanup_candidates = analyzer.analyze_containers(values_json_path)
    
    # Set output for GitHub Actions
    if 'GITHUB_OUTPUT' in os.environ:
        with open(os.environ['GITHUB_OUTPUT'], 'a') as f:
            f.write(f"cleanup_count={len(cleanup_candidates)}\n")

if __name__ == "__main__":
    main()