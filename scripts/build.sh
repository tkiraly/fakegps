#!/bin/bash

if [ -z ${COMMITCOUNT+x} ]; then
    #this happens when no env variable is set
    TAG=$(git tag -l --sort=-v:refname | head -n 1 | cut -c 2-)
    IFS='.' read -ra vers <<< "$TAG"
    MAJOR="${vers[0]}"
    MINOR="${vers[1]}"
    COMMITCOUNT=$(git rev-parse --short HEAD)
fi

go build -ldflags "-linkmode external -extldflags -static" -ldflags \
"-X 'github.com/tkiraly/fakegps/cmd.version=$MAJOR.$MINOR.$COMMITCOUNT'" \
-o fakegps


