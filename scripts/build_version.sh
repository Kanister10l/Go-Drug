#!/bin/bash

version=$(git rev-parse HEAD)
time=$(date +"%d/%m/%Y %T")

echo "${time}"

echo "package main" > version.go
echo "" >> version.go
echo "var Version = \"${version}\"" >> version.go
echo "" >> version.go
echo "var Time = \"$time\"" >> version.go