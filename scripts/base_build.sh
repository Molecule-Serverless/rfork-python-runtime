#!/bin/bash
BUNDLE_PATH=~/.base/baseline

cd ../baseline-image
docker build -t baseline-image .
sudo rm -rf $BUNDLE_PATH/rootfs
sudo mkdir -p $BUNDLE_PATH/rootfs
if [[ ! -f "$BUNDLE_PATH/config-base.json" ]]; then
    echo "Cannot find config.json. Paste a new one"
    cp config.json $BUNDLE_PATH/config-base.json
    cp config-loop.json $BUNDLE_PATH
fi
sudo docker export `docker create baseline-image` | sudo tar -C $BUNDLE_PATH/rootfs -xf -
