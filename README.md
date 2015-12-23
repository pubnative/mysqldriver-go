# mysqldriver-go
[![Build Status](https://travis-ci.org/pubnative/mysqldriver-go.svg?branch=master)](https://travis-ci.org/pubnative/mysqldriver-go)
[![GoDoc](https://godoc.org/github.com/pubnative/mysqldriver-go?status.svg)](https://godoc.org/github.com/pubnative/mysqldriver-go)

## Motivation
There are already many MySQL drivers which implement [database/sql](https://golang.org/pkg/database/sql/) interface however using this generic interface, especialy [Scan](https://golang.org/pkg/database/sql/#Row.Scan) method, requires to store many objects in a HEAP. Reading massive number of records from DB can significantly increase GC pause time which is very sensitive for low-latency applications. Due to this issue, was made a decision to write another MySQL driver which is GC friendly as much as possible, and not to follow [database/sql](https://golang.org/pkg/database/sql/) interface.
