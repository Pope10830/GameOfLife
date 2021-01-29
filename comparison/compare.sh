#!/usr/bin/env bash

rm -f your-time.txt
touch your-time.txt
rm -f your-out.txt
touch your-out.txt
rm -f base-time.txt
touch base-time.txt
rm -f base-out.txt
touch base-out.txt

go test -c -o gameoflife.test

echo "Benchmarking..."

benchtime=10x

#for b in 128x128x2 128x128x4 128x128x8 512x512x2 512x512x4 512x512x8
for b in 128x128x2 128x128x4 128x128x8
do
    echo ${b} on your solution
    \time -f '%P' -o your-time.txt -a ./gameoflife.test -test.run XXX -test.bench /${b} -test.benchtime ${benchtime} >> your-out.txt
    echo ${b} on baseline solution
    \time -f '%P' -o base-time.txt -a ./baseline.test -test.run XXX -test.bench /${b} -test.benchtime ${benchtime} >> base-out.txt
done

go build comparison/compare.go
./compare base-time.txt your-time.txt base-out.txt your-out.txt