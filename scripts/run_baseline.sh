#!/bin/bash
source ./config
BUNDLE_PATH=/run/.base/baseline
cd $BUNDLE_PATH
sudo $RUNC run baseline
