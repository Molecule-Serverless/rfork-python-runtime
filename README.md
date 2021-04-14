# README

## Make base fs for the runc container

```bash
make base-image # build the base image
make base-fs
```

## Compile our own version of runc

```bash
git clone -b add-fork-command https://ipads.se.sjtu.edu.cn:1312/xcontainer/runc.git
cd runc
make static
export RUNC=${PWD}/runc
```

## Launching the zygote container

```bash
sudo -E RUNC=${RUNC} sh ./scripts/launch_zygote_container.sh
sudo ${RUNC} list # will show the "python-test" container
```

## Forking the zygote container

```bash
sudo -E RUNC=${RUNC} sh ./scripts/fork_zygote_container.sh
sudo ${RUNC} list # will show the "new-python-test" container
```

## Scripts
* How to run a test
``` bash
cd scripts
./kill_containers.sh
./template_build.sh
./zygote_build.sh
./run_fork.sh
```