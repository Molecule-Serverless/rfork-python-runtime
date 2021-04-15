#!/bin/bash
source ./config
BUNDLE_PATH=~/.base/baseline
cd $BUNDLE_PATH
sudo $RUNC run baseline
