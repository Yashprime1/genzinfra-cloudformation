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


# Write JSON data to file
def write_json_to_file(filename, data):
    with open(filename, 'w') as file:
        json.dump(data, file, indent=4)

# Sample usage
if __name__ == "__main__":
    filename = 'prod.json' 
    stack_data = read_json_from_file(filename)

    region = "ap-south-1" 
    stackname = os.getenv('bamboo_deploy_environment')

    if "dash" in stackname:
        deployimage = "yashprime07/dashboard:"+os.getenv('bamboo_planKey')+"-"+os.getenv('bamboo_buildNumber')
    else:
        deployimage = "yashprime07/notificationbackend:"+os.getenv('bamboo_planKey')+"-"+os.getenv('bamboo_buildNumber')
    stack_data["Stacks"][region][stackname]["Parameters"]["ContainerImage"]=deployimage


    write_json_to_file(filename,stack_data)
    stack_data = read_json_from_file(filename)
    print(stack_data)
    