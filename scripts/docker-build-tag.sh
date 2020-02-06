#!/bin/bash

TAG=$(git tag -l --sort=-v:refname | head -n 1 | cut -c 2-)
IFS='.' read -ra vers <<< "$TAG"
MAJOR="${vers[0]}"
MINOR="${vers[1]}"
COMMITCOUNT="${vers[2]}"

docker build --build-arg MINOR=$MINOR --build-arg MAJOR=$MAJOR --build-arg COMMITCOUNT=$COMMITCOUNT\
  -t registry.gitlab.com/loranna/fakegps:$MAJOR.$MINOR.$COMMITCOUNT \
  -t registry.gitlab.com/loranna/fakegps:latest .
