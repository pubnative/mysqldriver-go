# mysqldriver-go
[![Build Status](https://travis-ci.org/pubnative/mysqldriver-go.svg?branch=master)](https://travis-ci.org/pubnative/mysqldriver-go)
[![GoDoc](https://godoc.org/github.com/pubnative/mysqldriver-go?status.svg)](https://godoc.org/github.com/pubnative/mysqldriver-go)

## Motivation
There are already many MySQL drivers which implement [database/sql](https://golang.org/pkg/database/sql/) interface however using this generic interface, especialy [Scan](https://golang.org/pkg/database/sql/#Row.Scan) method, requires to store many objects in a HEAP. Reading massive number of records from DB can significantly increase GC pause time which is very sensitive for low-latency applications. Due to this issue, was made a decision to write another MySQL driver which is GC friendly as much as possible, and not to follow [database/sql](https://golang.org/pkg/database/sql/) interface.

Benchmark was performed on MacBook Pro (Retina, 13-inch, Late 2013), 2.8 GHz Intel Core i7, 16 GB 1600 MHz DDR3
```zsh
➜  benchmarks git:(master) ✗ go run main.go 
mysqldriver: records read 100  HEAP 129  time 722.293µs
go-sql-driver: records read 100  HEAP 335  time 716.416µs
mysqldriver: records read 1000  HEAP 1015  time 633.537µs
go-sql-driver: records read 1000  HEAP 3010  time 798.109µs
mysqldriver: records read 10000  HEAP 10092  time 3.137886ms
go-sql-driver: records read 10000  HEAP 30010  time 3.377241ms
```
