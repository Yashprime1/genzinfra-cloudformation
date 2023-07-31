# x1-Mongo Stacks Dev Deploy
## Prerequisites:
To deploy the Mongo Stacks to development environment, you need to deploy firstly the foundation stacks (dev.foundation.json).
https://bitbucket.clevertap.net/projects/INFRA/repos/cloudformation/browse/README_local.md

You will need to build the GO templates to deploy the stack.
With the following commands: 
  go build -o dist/cloudformation
  chmod +x dist/cloudformation && ./dist/cloudformation  
## Deploy the Mongo Stacks
You have to deploy the the dev.manifest.json.
It will deploy the go lang templates.
cfstack deploy --manifest dev.manifest.json --role arn:aws:iam::726622043621:role/cfstack-Init-CloudFormationServiceIamRole-1OQP94KQAFTTS

## Deploy the single Mongo Stack
You have to replace the stack name with one that you want to deploy.In the example below , the deployed stack name is x1-Mongo-Delivery-1.

cfstack deploy --manifest dev.manifest.json stack -n x1-Mongo-Delivery-1 -r us-east-1 --role arn:aws:iam::726622043621:role/cfstack-Init-CloudFormationServiceIamRole-1OQP94KQAFTTS

