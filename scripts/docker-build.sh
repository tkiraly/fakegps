#!/bin/bash

docker build --build-arg=GOPROXY=$ATHENS --build-arg=GONOSUMDB=$GONOSUMDB -t registry.gitlab.com/loranna/fakegps:latest .
