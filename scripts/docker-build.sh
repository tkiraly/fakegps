#!/bin/bash

docker build --build-arg=GOPROXY=$ATHENS --build-arg=GONOSUMDB=$GONOSUMDB -t registry.github.com/tkiraly/fakegps:latest .
