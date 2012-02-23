#!/bin/sh
GOPATH=`pwd`; go install monkey
cp src/monkey/*.glsl bin/
