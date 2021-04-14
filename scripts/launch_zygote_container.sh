${RUNC} delete -f python-test
${RUNC} run --bundle .base/fs -d python-test
