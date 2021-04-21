import os
import re
import sys
import time
TEST_TIMES = 16
TEST_INVOKETIME_PATTERN = {"baseline": "start run container", "fork": "start fork"}
USAGE="python3 test_baseline.py [test], test can be \"baseline\" or \"fork\"\nIf no test is specified, it runs all tests by default"

def test_fork_start():
    ENDPOINT_BUNDLE="%s/.base/spin0/rootfs" %os.environ['HOME']
    COMMAND_FORK = "./run_fork_loop.sh"
    COMMAND_DUMP_RUN = "./mem-usage-fork-show.sh"
    #COMMAND_LOOP_RUN = "./run_baseline_loop.sh %d"

    for i in range(TEST_TIMES):
        os.system(COMMAND_FORK)
        time.sleep(0.1)
        time.sleep(1)
        print("[Test Result, (fork) concurrent cases: %d]" % i)
        os.system(COMMAND_DUMP_RUN)
    # print(latencies)

# pre-requisite: finish the building of the baseline container bundle, i.e., ~/.base/baseline/rootfs and ~/.base/baseline/config.json
def test_baseline_start():
    # PWD must be in scripts/tests
    COMMAND_RUN = "./run_baseline.sh %d"
    COMMAND_LOOP_RUN = "./run_baseline_loop.sh %d"
    COMMAND_DUMP_RUN = "./mem-usage-show.sh"

    for i in range(TEST_TIMES):
        os.system(COMMAND_LOOP_RUN % i)
        time.sleep(0.1)
        print("[Test Result, (baseline) concurrent cases: %d]" % i)
        os.system(COMMAND_DUMP_RUN)
        #print(start_latency, e2e_latency)
    # print(latencies)


def parse_output_lines(output_lines, test):
    invokeTime_pattern_line = TEST_INVOKETIME_PATTERN[test]
    startTime_pattern_line = "\'startTime\': 1[0-9]{12}"
    retTime_pattern_line = "\'retTime\': 1[0-9]{12}"
    time_pattern = "1[0-9]{12}"

    # Find lines that contains the invokeTime and the startTime
    invokeTime_line = None
    startTime_line = None
    retTime_line = None
    for line in output_lines:
        invokeTime_match = re.search(invokeTime_pattern_line, line)
        if invokeTime_match != None:
            # print(line)
            invokeTime_line = line
            continue

        startTime_match = re.search(startTime_pattern_line, line)
        if startTime_match != None:
            # print(line)
            startTime_line = line

        retTime_match = re.search(retTime_pattern_line, line)
        if retTime_match != None:
            # print(line)
            retTime_line = line


    if startTime_line == None or invokeTime_line == None or retTime_line == None:
        print("error output: can't find the startTime or invokeTime")
        exit()
    # Find the time value in the lines
    invokeTime_search = re.search(time_pattern, invokeTime_line)
    startTime_search = re.search(startTime_pattern_line, startTime_line)
    retTime_search = re.search(retTime_pattern_line, retTime_line)

    invokeTime = invokeTime_line[invokeTime_search.span()[0]: invokeTime_search.span()[0] + 13]
    startTime = startTime_line[startTime_search.span()[1] - 13 : startTime_search.span()[1]]
    retTime = startTime_line[retTime_search.span()[1] - 13 : retTime_search.span()[1]]

    #print(startTime, retTime)

    return int(invokeTime), int(startTime), int(retTime)

def format_result(latencies, test):
    request_num = len(latencies)
    print("=============== %s result ===============" %test)
    latencies.sort()
    latency_sum = 0
    for latency in latencies:
        latency_sum += latency
    averageLatency = latency_sum / request_num
    _50pcLatency = latencies[int(request_num * 0.5) - 1]
    _75pcLatency = latencies[int(request_num * 0.75) - 1]
    _90pcLatency = latencies[int(request_num * 0.9) - 1]
    _95pcLatency = latencies[int(request_num * 0.95) - 1]
    _99pcLatency = latencies[int(request_num * 0.99) - 1]
    print("latency (ms):\navg\t50%\t75%\t90%\t95%\t99%")
    print("%.2f\t%d\t%d\t%d\t%d\t%d" %(averageLatency,_50pcLatency,_75pcLatency,_90pcLatency,_95pcLatency,_99pcLatency))

def format_scale_result(latencies, test):
    request_num = len(latencies)
    print("=============== %s Scale result ===============" %test)
    latencies.sort()
    latency_sum = 0
    i = 0
    for latency in latencies:
        print("%d\t%d" % (i, latency))
        i += 1

if __name__ == '__main__':
    if len(sys.argv) == 2:
        if sys.argv[1] == "fork":
            test_fork_start()
        elif sys.argv[1] == "baseline":
            test_baseline_start()
        else:
            print(USAGE)
    else:
        test_fork_start()
        test_baseline_start()
