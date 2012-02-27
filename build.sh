#!/bin/sh
GOPATH=`pwd`:$GOROOT; go install monkey
cp src/monkey/*.glsl bin/
