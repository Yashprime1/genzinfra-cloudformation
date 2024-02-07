import json
import os

def replace_container_image(stack_data, stackname, new_image):
    if stackname in stack_data["Stacks"]:
        stack = stack_data["Stacks"][stackname]
        for component in stack:
            if "ContainerImage" in stack[component]["Parameters"]:
                stack[component]["Parameters"]["ContainerImage"] = new_image
    else:
        print(f"Stack with name '{stackname}' not found.")

# Read JSON data from file
def read_json_from_file(filename):
    with open(filename, 'r') as file:
        return json.load(file)

# Sample usage
if __name__ == "__main__":
    filename = '../prod.json' 
    stack_data = read_json_from_file(filename)

    region = "ap-south-1" 
    stackname = os.getenv('bamboo_deploy_environment')

    print(stack_data)