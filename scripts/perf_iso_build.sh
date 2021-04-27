#!/bin/bash
IOBUNDLE_PATH=~/.base/iotensive
CPUBUNDLE_PATH=~/.base/baseline # Reuse baseline dir, so we can reuse test scripts now

if [ $# -eq 0 ]
then
	echo "Please pass the io-intensive app dir first!"
	exit -1
fi

if [ $# -eq 1 ]
then
	echo "Please pass the cpu/mem-intensive app dir second!"
	exit -1
fi

cd $1
docker build -t iotensive-image .
sudo rm -rf $IOBUNDLE_PATH/rootfs
sudo mkdir -p $IOBUNDLE_PATH/rootfs
if [[ ! -f "$IOBUNDLE_PATH/config-base.json" ]]; then
    echo "Cannot find config.json. Paste a new one"
    cp config.json $IOBUNDLE_PATH/config-base.json
    #cp config-loop.json $BUNDLE_PATH
fi
sudo docker export `docker create iotensive-image` | sudo tar -C $IOBUNDLE_PATH/rootfs -xf -

cd $2
docker build -t cputensive-image .
sudo rm -rf $CPUBUNDLE_PATH/rootfs
sudo mkdir -p $CPUBUNDLE_PATH/rootfs
if [[ ! -f "$CPUBUNDLE_PATH/config-base.json" ]]; then
    echo "Cannot find config.json. Paste a new one"
    cp config.json $CPUBUNDLE_PATH/config-base.json
    #cp config-loop.json $BUNDLE_PATH
fi
sudo docker export `docker create cputensive-image` | sudo tar -C $CPUBUNDLE_PATH/rootfs -xf -
