#!/bin/bash
echo $bamboo_deploy_environment
echo $bamboo_buildResultKey
python3 scripts/replace_tags.py
git --version
git config --global user.name "Yashdeep Shetty"
git config --global user.email yashdeep@clevertap.com
git clone https://x-access-token:$bamboo_gittoken@github.com/Yashprime1/genzinfra-cloudformation.git
git commit -am "Updating tags to deploy"
git push --set-upstream origin master