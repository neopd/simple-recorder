# Introduction
This module receives data from NATS and write to a file to measure the performance of streaming via NATS.

# How to test
* Run simpleRecorder
```
./simpleRecorder -nats=172.20.10.120:4222 -subject=area-*.cam-*.*
```
