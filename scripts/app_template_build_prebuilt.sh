#!/bin/bash
BUNDLE_PATH=~/.base/container0

if [ $# -eq 0 ]
then
	echo "Please pass the app (baseline) name first!"
	exit -1
fi

docker pull $1 >/dev/null 2>&1
sudo rm -rf $BUNDLE_PATH/rootfs
sudo mkdir -p $BUNDLE_PATH/rootfs
if [[ ! -f "$BUNDLE_PATH/config-base.json" ]]; then
    echo "Not find config.json. Paste a new one"
    cp configs/container0/config.json $BUNDLE_PATH/config-base.json
    cp configs/container0/config-loop.json $BUNDLE_PATH/config-loop.json
fi
sudo docker export `docker create $1` | sudo tar -C $BUNDLE_PATH/rootfs -xf -
