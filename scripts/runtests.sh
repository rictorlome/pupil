#!/bin/bash

mv tests/*_test.go .

if [ -z "$1" ]
then
    go test -v
elif [ $1 == 'bench' ]
then
    go test -v -bench=.
else
    go test -v -run $1
fi

mv *_test.go tests/
