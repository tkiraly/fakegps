#!/bin/bash

TAG=$(git tag -l --sort=-v:refname | head -n 1 | cut -c 2-)
IFS='.' read -ra vers <<< "$TAG"
MAJOR="${vers[0]}"
MINOR="${vers[1]}"
COMMITCOUNT="${vers[2]}"

docker push registry.github.com/tkiraly/fakegps:$MAJOR.$MINOR.$COMMITCOUNT