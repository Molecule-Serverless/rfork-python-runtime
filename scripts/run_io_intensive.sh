#!/bin/bash

## Note(DD): This script is used for performance isolation test case!
## 	     Please carefully read the code before you run!
source ./config
BUNDLE_PATH=~/.base/iotensive
cd $BUNDLE_PATH
cp config-base.json config.json

## A loop to continuously run I/O intensive tasks
for (( i=0; i<10000; i++ ))
do
	sudo $RUNC run baseline
done
