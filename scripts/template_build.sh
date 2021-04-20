#!/bin/bash
BUNDLE_PATH=/run/.base/container0

cd ../python-base-image 
docker build -t python-base-image .
sudo rm -rf $BUNDLE_PATH/rootfs
sudo mkdir -p $BUNDLE_PATH/rootfs
if [[ ! -f "$BUNDLE_PATH/config.json" ]]; then
    echo "Cannot find config.json. Paste a new one"
    cp config.json $BUNDLE_PATH
fi
sudo docker export `docker create python-base-image` | sudo tar -C $BUNDLE_PATH/rootfs -xf -
