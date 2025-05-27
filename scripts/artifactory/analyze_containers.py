import requests
import json
import os
from datetime import datetime, timedelta, timezone
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
    
    def get_docker_images(self, repo_name,component_name):
        """Get all Docker images in a repository"""
        images = []
        
        # Method: Docker Registry API v2
        try:
            url = f"{self.base_url}/artifactory/api/docker/{repo_name}/v2/_catalog"
            response = self.session.get(url)
            if response.status_code == 200:
                data = response.json()
                images = data.get('repositories', [])
                if images:
                    print(f"    Found {len(images)} images using Docker Registry API v2")
                    return images
        except Exception as e:
            print(f"    Docker Registry API v2 failed: {e}")
        
        print(f"    No Docker images found in repository {repo_name}")
        return []
    
    def get_image_tags(self, repo_name, image_name):
        """Get all tags for a Docker image using multiple approaches"""
        tags = []
        
        # Method 1: Try Docker Registry API v2
        try:
            url = f"{self.base_url}/artifactory/api/docker/{repo_name}/v2/{image_name}/tags/list"
            response = self.session.get(url)
            if response.status_code == 200:
                data = response.json()
                tags = data.get('tags', [])
                if tags:
                    return tags
        except Exception as e:
            print(f"      Docker API tags failed: {e}")

        return list(set(tags)) if tags else ['latest']  # Default to 'latest' if no tags found
    
    def get_artifact_info(self, repo_name, image_name, tag):
        """Get detailed information about an artifact"""
        try:
            # Try multiple path structures for different Artifactory setups
            paths_to_try = [
                f"{image_name}/{tag}/manifest.json",
                f"{image_name}/{tag}",
                f"{image_name}/{tag}/latest",
                f"docker/{image_name}/{tag}",
                f"{image_name}/{tag}.json",
                f"{image_name}/manifest-{tag}.json",
                # Add more specific paths for your setup
                f"{image_name}/{tag}/sha256",
                f"{image_name}/tags/{tag}",
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
            
            # Last resort: try to get any information about the image directory
            image_url = f"{self.base_url}/artifactory/api/storage/{repo_name}/{image_name}"
            response = self.session.get(image_url)
            if response.status_code == 200:
                image_info = response.json()
                return {
                    'created': image_info.get('created', image_info.get('lastModified')),
                    'lastModified': image_info.get('lastModified'),
                    'size': image_info.get('size', 0),
                    'note': 'Image-level info (tag-specific info not available)'
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
    
    def create_cleanup_artifacts(self, cleanup_data):
        """Create cleanup artifacts organized by repository and image"""
        timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
        
        # Create detailed cleanup report
        cleanup_report = {
            'generated_at': datetime.now(timezone.utc).isoformat(),
            'summary': {
                'total_repositories': len(cleanup_data),
                'total_images': sum(len(repo_data['images']) for repo_data in cleanup_data.values()),
                'total_cleanup_candidates': sum(
                    len(image_data['cleanup_candidates']) 
                    for repo_data in cleanup_data.values() 
                    for image_data in repo_data['images'].values()
                )
            },
            'repositories': cleanup_data
        }
        
        # Save main cleanup report
        cleanup_report_file = f'cleanup_report_{timestamp}.json'
        with open(cleanup_report_file, 'w') as f:
            json.dump(cleanup_report, f, indent=2)
        print(f"Detailed cleanup report saved to: {cleanup_report_file}")
        
        # Create repository-specific artifacts
        for repo_name, repo_data in cleanup_data.items():
            if repo_data['images']:
                repo_file = f'cleanup_{repo_name}_{timestamp}.json'
                with open(repo_file, 'w') as f:
                    json.dump({
                        'repository': repo_name,
                        'generated_at': datetime.now(timezone.utc).isoformat(),
                        'summary': {
                            'total_images': len(repo_data['images']),
                            'total_cleanup_candidates': sum(
                                len(image_data['cleanup_candidates']) 
                                for image_data in repo_data['images'].values()
                            )
                        },
                        'images': repo_data['images']
                    }, f, indent=2)
                print(f"Repository-specific cleanup file saved to: {repo_file}")
        
        # Create cleanup commands script
        commands_file = f'cleanup_commands_{timestamp}.sh'
        with open(commands_file, 'w') as f:
            f.write("#!/bin/bash\n")
            f.write("# Artifactory Docker Image Cleanup Commands\n")
            f.write(f"# Generated on: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
            f.write("# WARNING: Review carefully before executing!\n\n")
            
            for repo_name, repo_data in cleanup_data.items():
                if any(image_data['cleanup_candidates'] for image_data in repo_data['images'].values()):
                    f.write(f"\n# Repository: {repo_name}\n")
                    f.write("-" * 50 + "\n")
                    
                    for image_name, image_data in repo_data['images'].items():
                        if image_data['cleanup_candidates']:
                            f.write(f"\n# Image: {image_name}\n")
                            for candidate in image_data['cleanup_candidates']:
                                artifact_path = f"{repo_name}/{image_name}/{candidate['tag']}"
                                f.write(f"# Delete {candidate['full_name']} (created: {candidate['created_date'][:10]}, {candidate['age_days']} days old)\n")
                                f.write(f"curl -X DELETE -u $ARTIFACTORY_USERNAME:$ARTIFACTORY_PASSWORD \\\n")
                                f.write(f"  '$ARTIFACTORY_URL/artifactory/{artifact_path}'\n")
                                f.write(f"echo 'Deleted: {candidate['full_name']}'\n\n")
        
        os.chmod(commands_file, 0o755)  # Make script executable
        print(f"Cleanup commands script saved to: {commands_file}")
        
        # Create summary CSV for easy viewing
        csv_file = f'cleanup_summary_{timestamp}.csv'
        with open(csv_file, 'w') as f:
            f.write("Repository,Image,Tag,Full Name,Created Date,Age Days,Status\n")
            for repo_name, repo_data in cleanup_data.items():
                for image_name, image_data in repo_data['images'].items():
                    for candidate in image_data['cleanup_candidates']:
                        f.write(f"{repo_name},{image_name},{candidate['tag']},{candidate['full_name']},{candidate['created_date'][:10]},{candidate['age_days']},CLEANUP_CANDIDATE\n")
                    for artifact in image_data.get('all_artifacts', []):
                        if artifact not in image_data['cleanup_candidates']:
                            status = "REFERENCED" if artifact.get('is_referenced') else "RECENT"
                            f.write(f"{repo_name},{image_name},{artifact.get('tag', 'unknown')},{artifact.get('full_name', 'unknown')},{artifact.get('created_date', 'unknown')[:10]},{artifact.get('age_days', 'unknown')},{status}\n")
        
        print(f"Summary CSV saved to: {csv_file}")
        
        return {
            'cleanup_report': cleanup_report_file,
            'commands_script': commands_file,
            'summary_csv': csv_file
        }
    
    def analyze_containers(self, values_json_path,component_name, months_threshold=3):
        """Analyze containers and identify candidates for cleanup"""
        print("Starting Artifactory container analysis...")
        print("=" * 60)
        
        # Load referenced images from values.json
        referenced_images = self.load_values_json(values_json_path)
        
        # Calculate cutoff date (make it timezone-aware to match parsed dates)
        from datetime import timezone
        cutoff_date = datetime.now(timezone.utc) - timedelta(days=months_threshold * 30)
        print(f"Analyzing containers older than: {cutoff_date.strftime('%Y-%m-%d')}")
        print(f"Referenced images in values.json: {len(referenced_images)}")
        print("=" * 60)
        
        cleanup_data = {}
        total_analyzed = 0
        
        # Get all Docker repositories
        repositories = self.get_repositories()
        print(f"Found {len(repositories)} Docker repositories")
        for repo in repositories:
            repo_name = repo['key']
            print(f"\nAnalyzing repository: {repo_name}")
            print("-" * 40)
            if not repo_name.equals("clevertap"):
                continue
            # Initialize repository data structure
            cleanup_data[repo_name] = {
                'repository_info': repo,
                'images': {},
                'summary': {
                    'total_images': 0,
                    'total_tags': 0,
                    'cleanup_candidates_count': 0
                }
            }
            
            # Get all images in repository
            images = self.get_docker_images(repo_name,component_name)
            
            if not images:
                print(f"  No Docker images found in {repo_name}")
                continue
            
            cleanup_data[repo_name]['summary']['total_images'] = len(images)
            
            for image_name in images:
                if not image_name.equals(component_name):
                    continue
                print(f"  Image: {image_name}")
            
                # Initialize image data structure
                cleanup_data[repo_name]['images'][image_name] = {
                    'all_artifacts': [],
                    'cleanup_candidates': [],
                    'referenced_artifacts': [],
                    'recent_artifacts': [],
                    'summary': {
                        'total_tags': 0,
                        'cleanup_candidates_count': 0,
                        'referenced_count': 0,
                        'recent_count': 0
                    }
                }
                
                # Get all tags for this image
                tags = self.get_image_tags(repo_name, image_name)
                cleanup_data[repo_name]['images'][image_name]['summary']['total_tags'] = len(tags)
                cleanup_data[repo_name]['summary']['total_tags'] += len(tags)
                
                for tag in tags:
                    total_analyzed += 1
                    full_image_name = f"{image_name}:{tag}"
                    
                    # Get artifact information
                    artifact_info = self.get_artifact_info(repo_name, image_name, tag)
                    
                    artifact_data = {
                        'repository': repo_name,
                        'image': image_name,
                        'tag': tag,
                        'full_name': full_image_name,
                        'artifact_info': artifact_info
                    }
                    
                    if artifact_info:
                        # Try multiple date fields and formats
                        date_fields = ['created', 'lastModified', 'lastUpdated']
                        created_date = None
                        date_source = None
                        
                        for field in date_fields:
                            date_str = artifact_info.get(field)
                            if date_str:
                                try:
                                    # Try different date parsing approaches
                                    if isinstance(date_str, str):
                                        # Handle ISO format with timezone (most common)
                                        created_date = parser.parse(date_str)
                                        # Ensure timezone awareness for comparison
                                        if created_date.tzinfo is None:
                                            created_date = created_date.replace(tzinfo=timezone.utc)
                                    elif isinstance(date_str, (int, float)):
                                        # Handle timestamp as number  
                                        created_date = datetime.fromtimestamp(date_str / 1000 if date_str > 1e10 else date_str, tz=timezone.utc)
                                    
                                    if created_date:
                                        date_source = field
                                        break
                                        
                                except Exception as date_error:
                                    continue
                        
                        if created_date:
                            is_old = created_date < cutoff_date
                            is_referenced = any(ref_img in full_image_name or full_image_name in ref_img 
                                                for ref_img in referenced_images)
                            
                            # Add date information to artifact data
                            artifact_data.update({
                                'created_date': created_date.isoformat(),
                                'age_days': (datetime.now(timezone.utc) - created_date).days,
                                'date_source': date_source,
                                'is_old': is_old,
                                'is_referenced': is_referenced
                            })
                            
                            status = []
                            if is_old:
                                status.append("OLD")
                            if is_referenced:
                                status.append("REFERENCED")
                            if not status:
                                status.append("RECENT")
                            
                            print(f"    Tag: {tag:20} | Date: {created_date.strftime('%Y-%m-%d')} ({date_source}) | Status: {', '.join(status)}")
                            
                            # Categorize artifacts
                            if is_old and not is_referenced:
                                cleanup_data[repo_name]['images'][image_name]['cleanup_candidates'].append(artifact_data)
                                cleanup_data[repo_name]['images'][image_name]['summary']['cleanup_candidates_count'] += 1
                                cleanup_data[repo_name]['summary']['cleanup_candidates_count'] += 1
                            elif is_referenced:
                                cleanup_data[repo_name]['images'][image_name]['referenced_artifacts'].append(artifact_data)
                                cleanup_data[repo_name]['images'][image_name]['summary']['referenced_count'] += 1
                            else:
                                cleanup_data[repo_name]['images'][image_name]['recent_artifacts'].append(artifact_data)
                                cleanup_data[repo_name]['images'][image_name]['summary']['recent_count'] += 1
                            
                            cleanup_data[repo_name]['images'][image_name]['all_artifacts'].append(artifact_data)
                            
                        else:
                            print(f"    Tag: {tag:20} | Date: NO PARSEABLE DATE | Available fields: {list(artifact_info.keys())}")
                            # Show date field values for debugging
                            for field in date_fields:
                                value = artifact_info.get(field, 'NOT FOUND')
                            
                            # Add to all artifacts even without date
                            artifact_data.update({
                                'created_date': None,
                                'age_days': None,
                                'date_source': None,
                                'is_old': False,
                                'is_referenced': False,
                                'status': 'NO_DATE_INFO'
                            })
                            cleanup_data[repo_name]['images'][image_name]['all_artifacts'].append(artifact_data)
                    else:
                        print(f"    Tag: {tag:20} | Date: INFO NOT AVAILABLE | Status: UNKNOWN")
                        
                        # Add to all artifacts even without info
                        artifact_data.update({
                            'created_date': None,
                            'age_days': None,
                            'date_source': None,
                            'is_old': False,
                            'is_referenced': False,
                            'status': 'NO_ARTIFACT_INFO'
                        })
                        cleanup_data[repo_name]['images'][image_name]['all_artifacts'].append(artifact_data)
        
        # Create cleanup artifacts
        print("\n" + "=" * 60)
        print("CREATING CLEANUP ARTIFACTS")
        print("=" * 60)
        
        artifacts_created = self.create_cleanup_artifacts(cleanup_data)
        
        # Print summary
        print("\n" + "=" * 60)
        print("ANALYSIS SUMMARY")
        print("=" * 60)
        print(f"Total containers analyzed: {total_analyzed}")
        
        total_cleanup_candidates = sum(
            len(image_data['cleanup_candidates']) 
            for repo_data in cleanup_data.values() 
            for image_data in repo_data['images'].values()
        )
        print(f"Cleanup candidates (old + not referenced): {total_cleanup_candidates}")
        
        # Repository-wise summary
        print(f"\nREPOSITORY BREAKDOWN:")
        print("-" * 60)
        for repo_name, repo_data in cleanup_data.items():
            if repo_data['images']:
                print(f"Repository: {repo_name}")
                print(f"  Total Images: {repo_data['summary']['total_images']}")
                print(f"  Total Tags: {repo_data['summary']['total_tags']}")
                print(f"  Cleanup Candidates: {repo_data['summary']['cleanup_candidates_count']}")
                
                # Image-wise breakdown
                for image_name, image_data in repo_data['images'].items():
                    if image_data['cleanup_candidates']:
                        print(f"    {image_name}: {len(image_data['cleanup_candidates'])} cleanup candidates")
                print()
        
        if total_cleanup_candidates > 0:
            print(f"\nTop 10 oldest cleanup candidates:")
            print("-" * 60)
            all_candidates = []
            for repo_data in cleanup_data.values():
                for image_data in repo_data['images'].values():
                    all_candidates.extend(image_data['cleanup_candidates'])
            
            # Sort by age and show top 10
            sorted_candidates = sorted(all_candidates, key=lambda x: x.get('age_days', 0), reverse=True)[:10]
            for i, candidate in enumerate(sorted_candidates, 1):
                print(f"{i:2d}. {candidate['full_name']:40} | {candidate.get('age_days', 'N/A'):4} days | {candidate['created_date'][:10] if candidate.get('created_date') else 'No date'}")
        
        print(f"\nARTIFACTS CREATED:")
        print("-" * 60)
        for artifact_type, filename in artifacts_created.items():
            print(f"{artifact_type.replace('_', ' ').title()}: {filename}")
        
        # Extract simple cleanup candidates list for backward compatibility
        cleanup_candidates = []
        for repo_data in cleanup_data.values():
            for image_data in repo_data['images'].values():
                cleanup_candidates.extend(image_data['cleanup_candidates'])
        
        return {
            'cleanup_candidates': cleanup_candidates,
            'detailed_data': cleanup_data,
            'artifacts_created': artifacts_created,
            'summary': {
                'total_analyzed': total_analyzed,
                'total_cleanup_candidates': total_cleanup_candidates,
                'repositories_analyzed': len(cleanup_data)
            }
        }

def main():
    # Get environment variables
    artifactory_url = os.getenv('ARTIFACTORY_URL')
    username = os.getenv('ARTIFACTORY_USERNAME')
    password = os.getenv('ARTIFACTORY_PASSWORD')
    component_name = os.getenv('COMPONENT_NAME')
    values_json_path = os.getenv('VALUES_JSON_PATH', 'values.json')
    
    if not all([artifactory_url, username, password]):
        print("Error: Missing required environment variables")
        print("Required: ARTIFACTORY_URL, ARTIFACTORY_USERNAME, ARTIFACTORY_PASSWORD")
        sys.exit(1)
    
    analyzer = ArtifactoryAnalyzer(artifactory_url, username, password)
    result = analyzer.analyze_containers(values_json_path,component_name)
    
    # Set output for GitHub Actions
    if 'GITHUB_OUTPUT' in os.environ:
        with open(os.environ['GITHUB_OUTPUT'], 'a') as f:
            f.write(f"cleanup_count={result['summary']['total_cleanup_candidates']}\n")
            f.write(f"repositories_analyzed={result['summary']['repositories_analyzed']}\n")
            f.write(f"total_analyzed={result['summary']['total_analyzed']}\n")

if __name__ == "__main__":
    main()