#!/bin/bash
echo $bamboo_deploy_environment
echo $bamboo_buildResultKey
python3 genzinfra-cloudformation/scripts/replace_tags.py
git --version
git config --global user.name "Yashdeep Shetty"
git config --global user.email yashdeep@clevertap.com
cd genzinfra-cloudformation
go run main.go
pip3 install cfstack tabulate
python3 execute.py
# git commit -am "Updating tags to deploy"
# git push --set-upstream origin master