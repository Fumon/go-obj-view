#!/bin/sh
export GOPATH=`pwd`:$GOROOT;
go install monkey
if [ $? -eq 0 ]; then
cp src/monkey/*.glsl bin/
fi
