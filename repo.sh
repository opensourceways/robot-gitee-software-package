#!/bin/bash

set -euo pipefail


init() {
    if [ -d $repo ]; then
       return
    fi

    git clone $repo_url
    cd $repo

    git config user.name $git_user
    git config user.password $git_password

    git remote add upstream ${upstream}

}

new_branch() {
    cd  $repo

    git checkout master

    git fetch upstream
    git rebase upstream/master

    git push origin master

    git checkout -b $branch_name

}

commit() {
    cd $repo

    git add .

    git commit -m 'apply new package'

    git push origin $branch_name

    git checkout master
}




cmd=$1
git_user=$2
git_password=$3
branch_name=$4

repo=community
upstream=https://gitee.com/openeuler/${repo}.git
repo_url=https://gitee.com/${git_user}/${repo}.git


case $cmd in
    "init")
        init
        ;;
    "new")
        new_branch
        ;;
    "commit")
        commit
        ;;
esac