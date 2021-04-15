#!/bin/bash
RUNNING_CONTAINERS=`expr \`sudo runc list | grep -c ""\` - 1`
# if no container is running, run the template and the endpoint container
if [[ $RUNNING_CONTAINERS = 0 ]]; then
    cd ~/.base/container0
    sudo runc run -d python-test
    echo "run python-test complete"
    cd ~/.base/spin0
    sudo runc run -d app-test
    echo "run app-test complete"
    echo "ready to fork..."
    sleep 1s # wait for containers to complete startup
fi
sudo ~/molecule/runc/runc fork2container --zygote python-test --target app-test 