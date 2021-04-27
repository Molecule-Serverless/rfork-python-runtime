#!/bin/bash
BUNDLE_PATH=~/.base/spin0

cd ../spin-base-image 
gcc spin.c -o spin -O2 -static && docker build -t spin-base-image -f Dockerfile .
docker build -t python-endpoint-base -f python_endpoint.Dockerfile .
sudo rm -rf $BUNDLE_PATH/rootfs
sudo mkdir -p $BUNDLE_PATH/rootfs

if [[ ! -f "$BUNDLE_PATH/config.json" ]]; then
    echo "Cannot find config.json. Paste a new one"
    cp config.json $BUNDLE_PATH
fi
sudo docker export `docker create python-endpoint-base` | sudo tar -C $BUNDLE_PATH/rootfs -xf -
