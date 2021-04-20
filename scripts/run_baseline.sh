#!/bin/bash
source ./config
BUNDLE_PATH=~/.base/baseline
cd $BUNDLE_PATH
cp config-base.json config.json
sudo $RUNC run baseline
