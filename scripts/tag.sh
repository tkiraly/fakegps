#!/bin/bash

COMMITCOUNT=$(git rev-list --count HEAD)
TAG=$(git tag -l --sort=-v:refname | head -n 1 | cut -c 2-)
IFS='.' read -ra vers <<< "$TAG"
MAJOR="${vers[0]}"
MINOR="${vers[1]}"

if [ $# -eq 1 ]; then
    if [ $1 = "mi" ]; then
        MINOR=$(($MINOR+1))
    fi

    if [ $1 = "ma" ]; then
        MAJOR=$(($MAJOR+1))
        MINOR="0"
    fi
fi

git tag "v$MAJOR.$MINOR.$COMMITCOUNT"
git tag -l --sort=-v:refname | head -n 1