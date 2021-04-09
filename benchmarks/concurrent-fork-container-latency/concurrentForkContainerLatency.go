package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var zygoteContainerName = "zygote%d"
var zygoteContainerBase = ".base/container%d"
var zygoteRootfs = ".base/container%d/rootfs"
var configJSONPath = ".base/container%d/config.json"
var runc = "runc"

var appContainerName = "app%d"
var appContainerBase = ".base/spin%d"
var appSocketPath = ".base/spin%d/rootfs/fork.sock"

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
	err = initEnviron(parallelCount)
	if err != nil {
		panic(err)
	}
	startContainers(parallelCount)
	defer removeContainers(parallelCount)
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
	time.Sleep(time.Second)
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
	_ = jsonResults
	// fmt.Println(string(jsonResults))
	fmt.Println(len(results))
	fmt.Println(avg(results))
}

func startContainers(parallelCount int) error {
	for i := 0; i < 1; i++ {
		// create 1 zygote container
		cmd := exec.Command(runc, "run", "-d", "--bundle", fmt.Sprintf(zygoteContainerBase, i), fmt.Sprintf(zygoteContainerName, i))
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	for i := 0; i < parallelCount; i++ {
		// create parallelCount app container
		cmd := exec.Command(runc, "run", "-d", "--bundle", fmt.Sprintf(appContainerBase, i), fmt.Sprintf(appContainerName, i))
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func removeContainers(parallelCount int) error {
	for i := 0; i < 1; i++ {
		// delete 1 zygote container
	deleteZygote:
		cmd := exec.Command(runc, "delete", "-f", fmt.Sprintf(zygoteContainerName, i))
		err := cmd.Run()
		if err != nil {
			goto deleteZygote
		}
	}

	for i := 0; i < parallelCount; i++ {
		// delete parallelCount app container
	deleteApp:
		cmd := exec.Command(runc, "delete", "-f", fmt.Sprintf(appContainerName, i))
		err := cmd.Run()
		if err != nil {
			goto deleteApp
		}
	}
	return nil
}

func benchmark(stop chan struct{}, resultsChan chan []int64, count int) {
	results := []int64{}
	socketName := fmt.Sprintf(appSocketPath, count)            // .base/spin0|1|2/rootfs/fork.sock
	thisContainerName := fmt.Sprintf(appContainerName, count)  // app0, app1, ...
	thisZygoteContainer := fmt.Sprintf(zygoteContainerName, 0) // only support zygote0
	for {
		select {
		case <-stop:
			resultsChan <- results
			return
		default:
		}
		result, err := routine(count, socketName, thisContainerName, thisZygoteContainer)
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

func routine(count int, socketName string, newContainerName string, thisZygoteContainer string) (int64, error) {

	// 0. delete the container, delete the socket file
	// socketName := fmt.Sprintf(socketPath, count)
	// newContainerName := fmt.Sprintf(containerName, count)
	// newContainerBase := fmt.Sprintf(containerBase, count)
	var start, end int64
	err := removeExistingSocket(socketName)
	if err != nil {
		panic(err)
		return -1, err
	}

	// startChpt!
	start = time.Now().UnixNano()

	// choice1: start by forking from a zygote
	startCmd := exec.Command(runc, "fork2container", "--zygote", thisZygoteContainer, "--target", newContainerName)

	// 1. wait for output
	output, err := startCmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		panic(err)
		return 0, err
	}

	// 2. check for the socket existence
	pollForExistence(socketName)

	// endChpt!
	end = time.Now().UnixNano()

	pid, err := strconv.Atoi(strings.Trim(string(output), " \n"))
	if err != nil {
		panic(err)
		return 0, err
	}

	err = killProcess(pid)

	if err != nil {
		panic(err)
		return 0, err
	}

	// fmt.Println(end - start) // TODO
	return end - start, nil
}

func initEnviron(parallelCount int) error {
	if runcPath, ok := os.LookupEnv("RUNC"); ok {
		runc = runcPath
	}
	cmd := exec.Command("./bootstrapRootFS", "1", fmt.Sprint(parallelCount)) // create 1 zygote and parallelCount spin containers to land
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
	}
	return err
}

func killProcess(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = p.Kill()
	if err != nil {
		return err
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
