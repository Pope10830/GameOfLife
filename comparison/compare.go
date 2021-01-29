package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
)

type bench struct {
	name   string
	result int64
}

type time struct {
	result int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readCpuTimes(file []byte) []time {
	numberRegex, _ := regexp.Compile(`\d+`)
	rows := numberRegex.FindAllString(string(file), -1)
	times := make([]time, len(rows))
	for i, row := range rows {
		result, err := strconv.Atoi(row)
		check(err)
		times[i] = time{result:result}
	}
	return times
}

//noinspection GoUnhandledErrorResult
func analyseCpuTimes() {
	baseBenchmarksFile, err := ioutil.ReadFile(os.Args[3])
	check(err)
	baseBenchmarks := readBenchmarks(baseBenchmarksFile)

	baseCpuTimesFile, err := ioutil.ReadFile(os.Args[1])
	check(err)
	newCpuTimesFile, err := ioutil.ReadFile(os.Args[2])
	check(err)
	baseTimes := readCpuTimes(baseCpuTimesFile)
	newTimes := readCpuTimes(newCpuTimesFile)

	if len(baseTimes) != len(newTimes) {
		panic("CPU time lengths don't match")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "Benchmark\tBaseline CPU usage\tYour CPU usage\t% Difference")
	fmt.Fprintln(w, "\t\t\tThe smaller the better")
	for i, b := range baseBenchmarks {
		fmt.Fprintln(w, b.name, "\t", baseTimes[i].result, "%\t", newTimes[i].result, "%\t", baseTimes[i].result*100/newTimes[i].result, "%")
	}
	fmt.Fprintln(w, "This is the percentage of the CPU that this job got. It's computed as (U + S) / E")
	fmt.Fprintln(w, "Where")
	fmt.Fprintln(w, "U\tTotal number of CPU-seconds that the process spent in user mode.")
	fmt.Fprintln(w, "S\tTotal number of CPU-seconds that the process spent in kernel mode.")
	fmt.Fprintln(w, "E\tElapsed real time")
}

func readBenchmarks(file []byte) []bench {
	numberRegex, _ := regexp.Compile(`\d+`)
	rowRegex, _ := regexp.Compile(`\d+x\d+x\d+-?\d*\s+\d+\s+\d+ ns/op`)
	rows := rowRegex.FindAllString(string(file), -1)
	benchmarks := make([]bench, len(rows))
	for i, row := range rows {
		fields := strings.Fields(row)
		result, err := strconv.ParseInt(numberRegex.FindString(fields[2]), 10, 64)
		check(err)
		benchmarks[i] = bench{name: fields[0], result: result}
	}
	return benchmarks
}

//noinspection GoUnhandledErrorResult
func analyseBenchmarks() {
	baseBenchmarksFile, err := ioutil.ReadFile(os.Args[3])
	check(err)
	newBenchmarksFile, err := ioutil.ReadFile(os.Args[4])
	check(err)
	baseBenchmarks := readBenchmarks(baseBenchmarksFile)
	newBenchmarks := readBenchmarks(newBenchmarksFile)
	if len(baseBenchmarks) != len(newBenchmarks) {
		panic("Benchmark lengths don't match")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "Benchmark\tBaseline result\tYour result\t% Difference")
	fmt.Fprintln(w, "\t(ns/1000 turns)\t(ns/1000 turns)\tThe bigger the better")
	for i, b := range baseBenchmarks {
		baselineResult := b.result
		newResult := newBenchmarks[i].result
		fmt.Fprintln(w, b.name, "\t", baselineResult, "\t", newResult, "\t", baselineResult*100/newResult, "%")

	}
}

func main() {
	fmt.Println()
	fmt.Println("TIME RESULTS")
	analyseBenchmarks()
	fmt.Println()
	fmt.Println("CPU USAGE RESULTS")
	analyseCpuTimes()
}
