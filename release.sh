#!/bin/bash

# stop on error
set -e

# log error and exit script
function fail {
  echo
  echo "ERROR -- $1"
  exit -1
}

function step {
  echo
  echo $1
  echo "---"
}

# check arguments
if [ $# != 3 ]; then
  echo "Usage: $0 <release version> <username> <container_registry_token>"
  fail "Error: no <release version> <username> <container_registry_token> specified"
fi

cd $(dirname $0)
BASEDIR=$(pwd)

ORIGINAL_BRANCH=$(git rev-parse --symbolic-full-name --abbrev-ref HEAD)
VERSION=$1
GH_USERNAME=$2
GH_CR_TOKEN=$3
TAG=release-$VERSION

step "Preparing to release version '$VERSION' from branch '$ORIGINAL_BRANCH'"

# update the codebase
git pull --rebase
git fetch --tags

# preflight checks
if [[ $(git diff --shortstat 2> /dev/null | tail -n1) != "" ]]; then
  fail "Uncommitted changes in git repository - aborting the release"
fi
if [[ $(git status --porcelain 2>/dev/null| grep "^??" | wc -l) > 0 ]]; then
  fail "Untracked files in git repository - aborting the release"
fi
if [[ $(git tag --list $TAG ) ]]; then
  fail "Tag '$TAG' already exists - aborting the release"
fi


RELEASE_BRANCH="release-$VERSION-from-$ORIGINAL_BRANCH"

step "Creating release branch $RELEASE_BRANCH"
if [ $(git branch --list $RELEASE_BRANCH) ]; then
   echo "Branch name $RELEASE_BRANCH already exists - deleting the existing branch"
   git branch -D $RELEASE_BRANCH
fi
git checkout -b $RELEASE_BRANCH $ORIGINAL_BRANCH

git tag $TAG -m "Release $VERSION"

step "Building and publishing Docker image"
docker build --tag ghcr.io/guanaco-io/notifications:$VERSION .
echo $GH_CR_TOKEN | docker login ghcr.io --username $GH_USERNAME --password-stdin
docker push ghcr.io/guanaco-io/notifications:$VERSION

step "Pushing tag '$TAG' to remote 'origin' and deleting release branch"
git push origin $TAG
git checkout $ORIGINAL_BRANCH
git branch -D $RELEASE_BRANCH

step "Done - tag $TAG has been created"
echo "What's next?"
echo "- check out the tag with 'git checkout $TAG'"
echo "- run the build and test the release locally"
echo "- update the release notes at https://github.com/guanaco-io/notification/releases/tag/$TAG"
