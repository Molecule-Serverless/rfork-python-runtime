#!/bin/bash
cd ~/.base/container0
sudo runc run -d python-test
echo "run python-test complete"
cd ~/.base/spin0
sudo runc run -d app-test
echo "run app-test complete"
echo "ready to fork..."
sleep 1s
sudo ~/molecule/runc/runc fork2container --zygote python-test --target app-test 