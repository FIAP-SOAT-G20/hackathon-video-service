#!/bin/sh
source ./scripts/menuselect.sh

TAG="[MAKE-NEWBRANCH]"

usage() {
    echo "Usage: new-branch [ -d | --dev ]
                  [ -c | --current ]
                  [ -m | --main ] 
                  [ -a | --alphafix ] 
                  [ -h | --hotfix ]
                  [ -f your_source_branch | --from your_source_branch ]"
    exit 2
}

# Get arguments passed 
PARSED_ARGUMENTS=$(getopt -a new_branch -o dcmahf: --long dev,current,main,alphafix,hotfix,from: -- "$@")

# echo $PARSED_ARGUMENTS

# Verify command format if source branch param exists
VALID_ARGUMENTS=$?
if [ "$VALID_ARGUMENTS" != "0" ]; then
    usage
fi

# Checkout and push branch to remote
function create_branch {
    current_branch=$(git rev-parse --abbrev-ref HEAD)

    source_branch=$2
    target_branch=$1

    if [[ $current_branch != $source_branch ]]; then
        echo "$TAG git checkout $source_branch && git pull origin $source_branch"
        git checkout $source_branch && git pull origin $source_branch || exit 0 #in case of error, maybe it's still missing to commit staged changes
    fi

    echo "$TAG git checkout -b $target_branch && git push --set-upstream origin $target_branch"
    git checkout -b $target_branch && git push --set-upstream origin $target_branch || echo "An error occured at create or push new branch"
}

# Get the source branch given by parameter, if any
function get_source_branch_by_param {
    source_branch=""

    eval set -- "$PARSED_ARGUMENTS"

    case "$1" in
        -d | --dev)   source_branch="dev"   ; break ;;
        -c | --current)   source_branch=$(git rev-parse --abbrev-ref HEAD)    ; break ;;
        -m | --main)   source_branch="main"   ; break ;;
        -a | --alphafix)   source_branch="main"   ; break ;;
        -h | --hotfix)   source_branch="main"   ; break ;;
        -f | --from)   shift; source_branch=$1   ; break ;;
        --) shift; break ;;
        *) echo "Invalid argument: $1";
            usage ;;
    esac

    echo $source_branch
}

# Jump to repository root
cd "$(git rev-parse --show-toplevel)"

echo "✨  ✨  New Branch ✨  ✨ "
echo "\n"

# Source branch from parameter
source_branch="$(get_source_branch_by_param)"

# Which Feature?
echo "✨ Branch Type: 📱"
options=("feature" "change" "fix" "removed" "refactor" "doc" "hotfix")
select_option "${options[@]}"
featureType="${options[$?]}"

# Description
read -p "✨ Description: 📝 "
description="${REPLY// /_}"

if [[ ! -z $featureType ]]; then
    branch="$featureType/"
fi
if [[ ! -z $description ]]; then
    branch="$branch$description/"
fi

# Remove last
length=${#branch}
endindex=$(expr $length - 1)
branch=${branch:0:$endindex}
echo "\n$branch\n"

if [[ $featureType = "hotfix" && $source_branch != "main" ]]; then
    echo "Hotfix will be created from main branch. Do you want to continue?"
    options=("Yes" "No")
    select_option "${options[@]}"
    answer="${options[$?]}"
    
    if [[ $answer = "Yes" ]]; then
        create_branch $branch "main"
    else
        echo "Hotfix can only be created from main branch, finishing script..."
    fi
elif [[ ! -z $source_branch ]]; then
    create_branch $branch $source_branch
else
    current_branch=$(git rev-parse --abbrev-ref HEAD)
    if [[ $current_branch != "dev" && $current_branch != "main" ]]; then
        branch_options+=("$current_branch")
    fi
    branch_options+=("main")
    echo "✨ Source Branch: 🚀"
    select_option "${branch_options[@]}"
    source_branch="${branch_options[$?]}"
    create_branch $branch $source_branch
fi
