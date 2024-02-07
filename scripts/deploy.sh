#!/bin/bash
echo $bamboo_deploy_environment
echo $bamboo_buildResultKey
python3 scripts/replace_tags.py
# git commit -am "Updating tags to deploy"
# git push --set-upstream origin master