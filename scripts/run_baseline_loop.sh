#!/bin/bash
source ./config
BUNDLE_PATH=/run/.base/baseline
cd $BUNDLE_PATH
cp config-loop.json config.json
sudo $RUNC run -d baseline-$1 > /dev/null 2>&1
