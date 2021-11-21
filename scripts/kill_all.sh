#!/bin/bash

source ./config
sudo ${RUNC} delete -f new-python-test > /dev/null 2>&1
sudo ${RUNC} delete -f python-test > /dev/null 2>&1

sudo ${RUNC} delete -f new-python-test > /dev/null 2>&1
sudo ${RUNC} delete -f python-test > /dev/null 2>&1

sudo ${RUNC} delete -f app-test > /dev/null 2>&1
