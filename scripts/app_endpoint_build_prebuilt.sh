#!/bin/bash
BUNDLE_PATH=~/.base/spin0

if [ $# -eq 0 ]
then
	echo "Please pass the app (baseline) name first!"
	exit -1
fi

docker pull $1  >/dev/null 2>&1
sudo rm -rf $BUNDLE_PATH/rootfs
sudo mkdir -p $BUNDLE_PATH/rootfs

if [[ ! -f "$BUNDLE_PATH/config.json" ]]; then
    echo "Not find config.json. Paste a new one"
    cp configs/spin0/config.json $BUNDLE_PATH
fi
sudo docker export `docker create $1` | sudo tar -C $BUNDLE_PATH/rootfs -xf -
