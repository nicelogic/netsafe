#!/bin/bash

version=$1

git add .
git commit
git tag $version
echo "update version: $version"
git push origin $version
git push origin main