#!/bin/sh
source ./scripts/menuselect.sh

bold=$(tput bold)
normal=$(tput sgr0)

# Verify if Github CLI is installed
if ! command -v gh &> /dev/null
then
    echo "Github CLI could not be found. Please install it before running this script."
    echo "👉 https://cli.github.com/"
    exit
fi

# Jump to repository root
cd "$(git rev-parse --show-toplevel)"

# Get the description given by parameter, if any
description="$@"

# Deletes the temporary file, if it exists here by an error
template=".github/pull_request_template.md"

echo "✨  ✨  Pull Request ✨  ✨ "
echo "\n"

echo "✨ Context: 📱"
options=("ANY" "API" "WORKER" "INFRA" "FRONTEND")
select_option "${options[@]}"
app="${options[$?]}" 

# # Pull request type
# Types: "Feature", "Change, "Fix", "Removed", "Refactor", "Docs", "Hotfix"
# get from brach name, ex: feature/update_tracking_view_infos then type is Feature,
branch=$(git rev-parse --abbrev-ref HEAD)
if [[ $branch == *"feature"* ]]; then
    type="Feature"
elif [[ $branch == *"change"* ]]; then
    type="Change"
elif [[ $branch == *"hotfix"* ]]; then
    type="Hotfix"
elif [[ $branch == *"removed"* ]]; then
    type="Removed"
elif [[ $branch == *"refactor"* ]]; then
    type="Refactor"
elif [[ $branch == *"doc"* ]]; then
    type="Docs"
elif [[ $branch == *"fix"* ]]; then
    type="Fix"
fi

# If no type is given, ask for it
if [[ -z $type ]]; then
    echo "✨ Type: 📱"
    options=("Feature" "Change" "Fix" "Removed" "Refactor" "Docs" "Hotfix")
    select_option "${options[@]}"
    type="${options[$?]}" 
fi

# Description
if [[ -z $description ]]; then
    # description=$(git show -s --format=%s) # old, get the last commit message

    # new: extract from branch name:
    # ex branch "feature/update_tracking_view_infos", 
    # the description should be: "update tracking view infos"
    # Step 1: remove the first part of the branch name
    branch=$(git rev-parse --abbrev-ref HEAD)
    description=$(echo $branch | sed -r 's/^.*\///g')
    # Step 2: replace _ with space
    description=$(echo $description | sed -r 's/_/ /g')
    # Step 3: make the first letter uppercase
    description=$(echo "$description" | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) substr($i,2)}1')
fi

# Join title
title="[$app] $type: $description"
echo "\n${bold}$title${normal}"

# label is the type in lowercase
label=$(echo $type | awk '{print tolower($0)}')

# Open pull request with gh CLI
gh pr create --title "$title" --base main --body-file $template --web --assignee @me --label $label
