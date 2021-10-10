package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

var RUNNING = true
var PARALLEL_COUNT = 12
var TEST_DURATION_SECONDS = 10

func main() {
	// s := launchProcess()
	// fmt.Printf("%s", s)
	i := 0
	c := make(chan int64)
	throughput := int64(0)
	for i < PARALLEL_COUNT {
		go benchmark(c, i)
		i = i + 1
	}
	// TODO: sleep for a while
	time.Sleep(time.Duration(TEST_DURATION_SECONDS) * time.Second)
	RUNNING = false
	i = 0
	for i < PARALLEL_COUNT {
		throughput += <-c
		i = i + 1
	}
	fmt.Printf("Parallel Count: %d, Test duration in Second: %d, Total Throughput (op/second): %d\n", PARALLEL_COUNT, TEST_DURATION_SECONDS, throughput/int64(TEST_DURATION_SECONDS))
}

func benchmark(c chan int64, goroutine_id int) {
	count := int64(0)
	for RUNNING == true {
		launchProcess(goroutine_id)
		count = count + 1
	}
	c <- count
}

func launchProcess(goroutine_id int) string {
	app_container := fmt.Sprintf("app%d", goroutine_id)
	cmd := exec.Command("/home/wtx/runc/runc", "fork2container", "--zygote", "c-zygote", "--target", app_container) // the main evaluation part
	// now := getTimeStamp()
	// fmt.Printf("%d\n", now)
	// err := cmd.Run()
	output, err := cmd.CombinedOutput()
	// now = getTimeStamp()
	// fmt.Printf("%d\n", now)
	if err != nil {
		log.Println(string(output))
		log.Fatal(err)
	}
	return string("")
}

func getTimeStamp() int64 {
	return time.Now().UnixNano()
}
