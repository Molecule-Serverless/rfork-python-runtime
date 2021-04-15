# README

## Scripts
* Prepare environment
``` bash
cd scripts
./kill_containers.sh # make sure that no old container exists

./base_build.sh # build baseline container's bundle

./template_build.sh # build template container's bundle
./endpoint_build.sh # build endpoint container's bundle

# test baseline
./run_baseline.sh

# test cfork
./run_fork.sh
```
* run tests
``` bash
cd scripts/tests
# usage: python3 test_baseline.py [test]
# test can be baseline or fork
# if no test is specified, it runs all tests by default
# Caution: if the test is "fork", please make sure that you have run ./run_fork.sh successfully to warm up the environment
python3 test_baseline.py
```

## Old scripts (outdated)
### Make base fs for the runc container

```bash
make base-image # build the base image
make base-fs
```

### Compile our own version of runc

```bash
git clone -b add-fork-command https://ipads.se.sjtu.edu.cn:1312/xcontainer/runc.git
cd runc
make static
export RUNC=${PWD}/runc
```

### Launching the zygote container

```bash
sudo -E RUNC=${RUNC} sh ./scripts/launch_zygote_container.sh
sudo ${RUNC} list # will show the "python-test" container
```

### Forking the zygote container

```bash
sudo -E RUNC=${RUNC} sh ./scripts/fork_zygote_container.sh
sudo ${RUNC} list # will show the "new-python-test" container
```

