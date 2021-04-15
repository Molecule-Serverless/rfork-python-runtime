#!/bin/bash
BUNDLE_PATH=~/.base/baseline
RUNC=~/molecule/runc/runc
cd $BUNDLE_PATH
sudo $RUNC run baseline
