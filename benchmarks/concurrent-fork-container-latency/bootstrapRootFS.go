package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var zygoteConfigJSON = `{
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

var spinConfigJSON = `{
	"ociVersion": "1.0.1-dev",
	"process": {
		"terminal": false,
		"user": {
			"uid": 0,
			"gid": 0
		},
		"args": [
			"/spin"
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
				"type": "pid"
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

var zygoteRootfs = ".base/container%d/rootfs"
var zygoteConfigJSONPath = ".base/container%d/config.json"
var zygoteBaseImage = "python-base-image"

var spinRootfs = ".base/spin%d/rootfs"
var spinConfigJSONPath = ".base/spin%d/config.json"
var spinBaseImage = "spin-base-image"

var zygoteContainerID string
var spinContainerID string

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <zygote-container-count> <spin-container-count>\n", os.Args[0])
		return
	}

	createZygoteContainer := exec.Command("docker", "create", zygoteBaseImage)
	output, err := createZygoteContainer.Output()
	zygoteContainerID = string(output)
	zygoteContainerID = strings.Trim(zygoteContainerID, " \n")
	removeZygoteContainer := exec.Command("docker", "rm", "-f", zygoteContainerID)
	defer removeZygoteContainer.Run()
	zygoteCount, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = initZygoteEnviron(zygoteCount)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = prepareAllZygoteRootfs(zygoteCount)
	if err != nil {
		fmt.Println(err)
		return
	}

	createSpinContainer := exec.Command("docker", "create", spinBaseImage)
	output, err = createSpinContainer.Output()
	spinContainerID = string(output)
	spinContainerID = strings.Trim(spinContainerID, " \n")
	removeSpinContainer := exec.Command("docker", "rm", "-f", spinContainerID)
	defer removeSpinContainer.Run()
	spinCount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = initSpinEnviron(spinCount)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = prepareAllSpinRootfs(spinCount)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func initZygoteEnviron(parallelCount int) error {
	for i := 0; i < parallelCount; i++ {
		rootfsString := fmt.Sprintf(zygoteRootfs, i)
		configJSONPathString := fmt.Sprintf(zygoteConfigJSONPath, i)
		os.MkdirAll(rootfsString, os.ModePerm)
		f, err := os.Create(configJSONPathString)
		defer f.Close()
		if err != nil {
			return err
		}
		_, err = f.WriteString(zygoteConfigJSON)
		if err != nil {
			return err
		}
	}
	return nil
}

func initSpinEnviron(parallelCount int) error {
	for i := 0; i < parallelCount; i++ {
		rootfsString := fmt.Sprintf(spinRootfs, i)
		configJSONPathString := fmt.Sprintf(spinConfigJSONPath, i)
		os.MkdirAll(rootfsString, os.ModePerm)
		f, err := os.Create(configJSONPathString)
		defer f.Close()
		if err != nil {
			return err
		}
		_, err = f.WriteString(spinConfigJSON)
		if err != nil {
			return err
		}
	}
	return nil
}

func prepareAllZygoteRootfs(parallelCount int) error {
	for i := 0; i < parallelCount; i++ {
		targetRootfs := fmt.Sprintf(zygoteRootfs, i)
		cmd := fmt.Sprintf("docker export %s | tar -C %s -xf -", zygoteContainerID, targetRootfs)
		fmt.Println(cmd)
		exportContainer := exec.Command("bash", "-c", cmd)
		err := exportContainer.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func prepareAllSpinRootfs(parallelCount int) error {
	for i := 0; i < parallelCount; i++ {
		targetRootfs := fmt.Sprintf(spinRootfs, i)
		cmd := fmt.Sprintf("docker export %s | tar -C %s -xf -", spinContainerID, targetRootfs)
		fmt.Println(cmd)
		exportContainer := exec.Command("bash", "-c", cmd)
		err := exportContainer.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
