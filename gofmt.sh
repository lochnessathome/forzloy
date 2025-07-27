#!/usr/bin/env bash

for f in `git status | grep ".go$" | awk '{print $NF}'` ; do
    go fmt $f
done

