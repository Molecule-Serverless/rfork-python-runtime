#!/bin/bash
cd ../python-base-image 
docker build -t python-base-image .
sudo rm -rf ~/.base/container0/rootfs
sudo mkdir ~/.base/container0/rootfs
sudo docker export `docker create python-base-image` | sudo tar -C ~/.base/container0/rootfs -xf -