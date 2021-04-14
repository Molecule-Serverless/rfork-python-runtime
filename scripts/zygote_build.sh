#!/bin/bash
cd ../spin-base-image 
gcc spin.c -o spin -O2 -static && docker build -t spin-base-image .
sudo rm -rf ~/.base/spin0/rootfs
sudo mkdir ~/.base/spin0/rootfs
sudo docker export `docker create spin-base-image` | sudo tar -C ~/.base/spin0/rootfs -xf -
