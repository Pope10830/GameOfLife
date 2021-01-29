focus ='worker|distributor'
ignore = 'strings|fmt'

gol:
	go build
	./gameoflife


# Add -run /[NAME]
# eg: -run /16x16x2-0
# to run a specific test
test:
	go test


# Use -benchtime [TIME][UNIT]
# eg: -benchtime 60s
# to force the benchmark to run for the specified amount of time

# Use -bench /[NAME]
# eg: -bench /16x16x2
# to run a specific benchmark

# bench will run all tests before benchmarking - they must all pass
bench:
	go test -bench .

compare:
	./comparison/compare.sh

trace:
	go test -run=Test/trace -trace trace.out
	go tool trace trace.out


# Requires graphviz to work correctly
cpuprofile:
	go test -bench /512x512x8  -cpuprofile cpu.prof
	go tool pprof -pdf -nodefraction=0 -unit=ms -focus=$(focus) -ignore=$(ignore) cpu.prof


# Interactive mode - Useful commands:
# top10
# list worker
# list distributor
cpuprofile-i:
	go test -bench /512x512x8  -cpuprofile cpu.prof
	go tool pprof -nodefraction=0 -unit=ms -focus=$(focus) -ignore=$(ignore) cpu.prof


# Requires graphviz to work correctly
memprofile:
	go test -bench /512x512x8  -memprofile mem.prof --memprofilerate=1
	go tool pprof -pdf -alloc_space -nodefraction=0 -unit=B -focus=$(focus) -ignore=$(ignore) mem.prof


# Interactive mode - Useful commands:
# top10
# list worker
# list distributor
memprofile-i:
	go test -bench /512x512x8  -memprofile mem.prof --memprofilerate=1
	go tool pprof -alloc_space -nodefraction=0 -unit=B -focus=$(focus) -ignore=$(ignore) mem.prof


perf:
	sudo perf stat -d go test -bench /512x512x2
	sudo perf stat -d go test -bench /512x512x4
	sudo perf stat -d go test -bench /512x512x8


time:
	time go test -bench /512x512x2
	time go test -bench /512x512x4
	time go test -bench /512x512x8


.PHONY: gameoflife compare baseline baseline.test
