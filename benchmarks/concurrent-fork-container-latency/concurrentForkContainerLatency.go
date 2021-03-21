package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var configJSON = `
{
	"ociVersion": "1.0.1-dev",
	"process": {
		"terminal": false,
		"user": {
			"uid": 0,
			"gid": 0
		},
		"args": [
			"python", "daemon.py"
		],
		"env": [
			"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
			"TERM=xterm"
		],
		"cwd": "/",
		"capabilities": {
			"bounding": [
				"CAP_CHOWN",
				"CAP_DAC_OVERRIDE",
				"CAP_FSETID",
				"CAP_FOWNER",
				"CAP_MKNOD",
				"CAP_NET_RAW",
				"CAP_SETGID",
				"CAP_SETUID",
				"CAP_SETFCAP",
				"CAP_SETPCAP",
				"CAP_NET_BIND_SERVICE",
				"CAP_SYS_CHROOT",
				"CAP_KILL",
				"CAP_AUDIT_WRITE",
				"CAP_SYS_ADMIN"
			],
			"effective": [
				"CAP_CHOWN",
				"CAP_DAC_OVERRIDE",
				"CAP_FSETID",
				"CAP_FOWNER",
				"CAP_MKNOD",
				"CAP_NET_RAW",
				"CAP_SETGID",
				"CAP_SETUID",
				"CAP_SETFCAP",
				"CAP_SETPCAP",
				"CAP_NET_BIND_SERVICE",
				"CAP_SYS_CHROOT",
				"CAP_KILL",
				"CAP_AUDIT_WRITE",
				"CAP_SYS_ADMIN"
			],
			"inheritable": [
				"CAP_CHOWN",
				"CAP_DAC_OVERRIDE",
				"CAP_FSETID",
				"CAP_FOWNER",
				"CAP_MKNOD",
				"CAP_NET_RAW",
				"CAP_SETGID",
				"CAP_SETUID",
				"CAP_SETFCAP",
				"CAP_SETPCAP",
				"CAP_NET_BIND_SERVICE",
				"CAP_SYS_CHROOT",
				"CAP_KILL",
				"CAP_AUDIT_WRITE",
				"CAP_SYS_ADMIN"
			],
			"permitted": [
				"CAP_CHOWN",
				"CAP_DAC_OVERRIDE",
				"CAP_FSETID",
				"CAP_FOWNER",
				"CAP_MKNOD",
				"CAP_NET_RAW",
				"CAP_SETGID",
				"CAP_SETUID",
				"CAP_SETFCAP",
				"CAP_SETPCAP",
				"CAP_NET_BIND_SERVICE",
				"CAP_SYS_CHROOT",
				"CAP_KILL",
				"CAP_AUDIT_WRITE",
				"CAP_SYS_ADMIN"
			],
			"ambient": [
				"CAP_CHOWN",
				"CAP_DAC_OVERRIDE",
				"CAP_FSETID",
				"CAP_FOWNER",
				"CAP_MKNOD",
				"CAP_NET_RAW",
				"CAP_SETGID",
				"CAP_SETUID",
				"CAP_SETFCAP",
				"CAP_SETPCAP",
				"CAP_NET_BIND_SERVICE",
				"CAP_SYS_CHROOT",
				"CAP_KILL",
				"CAP_AUDIT_WRITE",
				"CAP_SYS_ADMIN"
			]
		},
		"rlimits": [
			{
				"type": "RLIMIT_NOFILE",
				"hard": 1024,
				"soft": 1024
			}
		],
		"noNewPrivileges": true
	},
	"root": {
		"path": "rootfs",
		"readonly": false
	},
	"hostname": "runc",
	"mounts": [
		{
			"destination": "/proc",
			"type": "proc",
			"source": "proc"
		},
		{
			"destination": "/dev",
			"type": "tmpfs",
			"source": "tmpfs",
			"options": [
				"nosuid",
				"strictatime",
				"mode=755",
				"size=65536k"
			]
		},
		{
			"destination": "/dev/pts",
			"type": "devpts",
			"source": "devpts",
			"options": [
				"nosuid",
				"noexec",
				"newinstance",
				"ptmxmode=0666",
				"mode=0620",
				"gid=5"
			]
		},
		{
			"destination": "/dev/shm",
			"type": "tmpfs",
			"source": "shm",
			"options": [
				"nosuid",
				"noexec",
				"nodev",
				"mode=1777",
				"size=65536k"
			]
		},
		{
			"destination": "/dev/mqueue",
			"type": "mqueue",
			"source": "mqueue",
			"options": [
				"nosuid",
				"noexec",
				"nodev"
			]
		},
		{
			"destination": "/sys",
			"type": "sysfs",
			"source": "sysfs",
			"options": [
				"nosuid",
				"noexec",
				"nodev",
				"ro"
			]
		},
		{
			"destination": "/sys/fs/cgroup",
			"type": "cgroup",
			"source": "cgroup",
			"options": [
				"nosuid",
				"noexec",
				"nodev",
				"relatime",
				"ro"
			]
		}
	],
	"linux": {
		"resources": {
			"devices": [
				{
					"allow": false,
					"access": "rwm"
				}
			]
		},
		"namespaces": [
			{
				"type": "network"
			},
			{
				"type": "ipc"
			},
			{
				"type": "uts"
			},
			{
				"type": "mount"
			}
		],
		"maskedPaths": [
			"/proc/acpi",
			"/proc/asound",
			"/proc/kcore",
			"/proc/keys",
			"/proc/latency_stats",
			"/proc/timer_list",
			"/proc/timer_stats",
			"/proc/sched_debug",
			"/sys/firmware",
			"/proc/scsi"
		],
		"readonlyPaths": [
			"/proc/bus",
			"/proc/fs",
			"/proc/irq",
			"/proc/sys",
			"/proc/sysrq-trigger"
		]
	}
}
`

var containerName = "container%d"
var containerBase = ".base/container%d"
var rootfs = ".base/container%d/rootfs"
var configJSONPath = ".base/container%d/config.json"
var socketPath = ".base/container%d/rootfs/fork.sock"
var runc = "runc"
var zygoteContainerName = "python-test"

func main() {
	parallelCount, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	timeSpan, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	initEnviron(parallelCount)
	/*for i := 0; i < 500; i++ {
		span, err := routine(0)
		if err == nil {
			_ = span
			// fmt.Println(i)
			// fmt.Println(span)
		} else {
			fmt.Println(err.Error())
		}
	}*/
	_ = timeSpan
	stopChanArr, resultsChanArr := makeChannels(parallelCount)
	for i := 0; i < parallelCount; i++ {
		go benchmark(stopChanArr[i], resultsChanArr[i], i)
	}
	time.Sleep(time.Second * time.Duration(timeSpan))
	for i := 0; i < parallelCount; i++ {
		stopChanArr[i] <- struct{}{}
	}
	var results []int64
	for i := 0; i < parallelCount; i++ {
		result := <-resultsChanArr[i]
		results = append(results, result...)
	}
	jsonResults, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonResults))
	fmt.Println(len(results))
	fmt.Println(avg(results))
}

func benchmark(stop chan struct{}, resultsChan chan []int64, count int) {
	results := []int64{}
	socketName := fmt.Sprintf(socketPath, count)
	newContainerName := fmt.Sprintf(containerName, count)
	newContainerBase := fmt.Sprintf(containerBase, count)
	for {
		select {
		case <-stop:
			resultsChan <- results
			return
		default:
		}
		result, err := routine(count, socketName, newContainerName, newContainerBase)
		if err != nil {
			panic(err)
		}
		results = append(results, result)
	}
}

func makeChannels(parallelCount int) (stopChanArr []chan struct{}, resultsChanArr []chan []int64) {
	stopChanArr = make([]chan struct{}, parallelCount)
	resultsChanArr = make([]chan []int64, parallelCount)
	for i := 0; i < parallelCount; i++ {
		stopChanArr[i] = make(chan struct{})
		resultsChanArr[i] = make(chan []int64)
	}
	return stopChanArr, resultsChanArr
}

func routine(count int, socketName string, newContainerName string, newContainerBase string) (int64, error) {

	// 0. delete the container, delete the socket file
	// socketName := fmt.Sprintf(socketPath, count)
	// newContainerName := fmt.Sprintf(containerName, count)
	// newContainerBase := fmt.Sprintf(containerBase, count)
	var start, end int64
	deleteContainer(newContainerName)
	err := removeExistingSocket(socketName)
	if err != nil {
		return -1, err
	}

	// choice1: start a container by forking from a zygote
	// startCmd := exec.Command(runc, "fork", "--bundle", newContainerBase, zygoteContainerName, "rootfs/fork.sock", newContainerName)

	// choice2: start a container from scratch
	startCmd := exec.Command(runc, "run", "--bundle", newContainerBase, "-d", newContainerName)

	// startChpt!
	start = time.Now().UnixNano()

	// 1. fork a container
	startCmd.Start()

	// 2. check for the socket existence
	pollForExistence(socketName)

	// endChpt!
	end = time.Now().UnixNano()

	// fmt.Println(end - start) // TODO
	return end - start, nil
}

func initEnviron(parallelCount int) error {
	if runcPath, ok := os.LookupEnv("RUNC"); ok {
		runc = runcPath
	}
	for i := 0; i < parallelCount; i++ {
		rootfsString := fmt.Sprintf(rootfs, i)
		configJSONPathString := fmt.Sprintf(configJSONPath, i)
		os.MkdirAll(rootfsString, os.ModePerm)
		f, err := os.Create(configJSONPathString)
		defer f.Close()
		if err != nil {
			return err
		}
		_, err = f.WriteString(configJSON)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteContainer(containerName string) error {
	deleteCmd := exec.Command(runc, "delete", "-f", containerName)
	deleteCmd.Run()
	for containerExist(containerName) {
		deleteCmd := exec.Command(runc, "delete", "-f", containerName)
		deleteCmd.Run()
	}
	/*var err error
	err = deleteCmd.Run()
	for err != nil {
		err = deleteCmd.Run()
	}*/
	return nil
}

func removeExistingSocket(sockName string) error {
	var err error
	if _, err = os.Stat(sockName); err == nil {
		err = os.Remove(sockName)
		return err
	}
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func pollForExistence(file string) {
	for {
		if _, err := os.Stat(file); err == nil {
			break
		}
	}
}

func containerExist(containerName string) bool {
	existCmd := exec.Command(runc, "state", containerName)
	err := existCmd.Run()
	if err == nil {
		return true
	}
	return false
}

func avg(results []int64) int64 {
	l := len(results)
	var sum int64 = 0
	for i := 0; i < l; i++ {
		sum += results[i]
	}
	return sum / int64(l)
}
