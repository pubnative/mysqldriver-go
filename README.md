# mysqldriver-go
[![Build Status](https://travis-ci.org/pubnative/mysqldriver-go.svg?branch=master)](https://travis-ci.org/pubnative/mysqldriver-go)
[![GoDoc](https://godoc.org/github.com/pubnative/mysqldriver-go?status.svg)](https://godoc.org/github.com/pubnative/mysqldriver-go)

## Table of contents

- [Motivation](#motivation)
- [Goal](#goal)
- [Documentation](#documentation)
- [Dependencies](#dependencies)
- [Installation](#installation)
- [Quick Start](#quick-start)

## Motivation
There are many MySQL drivers that implement the [database/sql](https://golang.org/pkg/database/sql/) interface.
However, using this generic interface, especialy in the [`Scan`](https://golang.org/pkg/database/sql/#Row.Scan) method, requires the storage of many objects in the heap. 

Reading a massive number of records from a DB can significantly increase the Garbage Collection (GC) pause-time that can be very sensitive for high-throughput, low-latency applications. 

Because of the above and the need for a GC-friendly MySQL driver, we've decided not to follow the [database/sql](https://golang.org/pkg/database/sql/) interface and write this driver.

The following [Benchmark](https://github.com/pubnative/mysqldriver-go/blob/master/benchmarks/main.go) was run on a `MacBook Pro (Retina, 13-inch, Late 2013), 2.8 GHz Intel Core i7, 16 GB 1600 MHz DDR3` using `Go 1.5.2`:

[![comparison](https://cloud.githubusercontent.com/assets/296795/12080839/72fcf55c-b268-11e5-9632-743ec07c2b80.png)](https://jsfiddle.net/zs83oze6/3/)
```zsh
➜  benchmarks git:(master) ✗ go run main.go 
mysqldriver: records read 100  HEAP 129  time 722.293µs
go-sql-driver: records read 100  HEAP 335  time 716.416µs
mysqldriver: records read 1000  HEAP 1015  time 633.537µs
go-sql-driver: records read 1000  HEAP 3010  time 798.109µs
mysqldriver: records read 10000  HEAP 10092  time 3.137886ms
go-sql-driver: records read 10000  HEAP 30010  time 3.377241ms
```

## Goal
The main goals of this library are: *performance* over flexibility, *simplicity* over complexity. Any new feature shouldn't decrease the performance of the exising code base. 

Any improvements to productivity are always welcome. There is no plan to convert this library into an ORM. The plan is to keep it simple, and still keep supporting all of the MySQL features.

## Documentation
1. [API Reference](https://godoc.org/github.com/pubnative/mysqldriver-go)
2. [Official MySQL Protocol Documentation](https://dev.mysql.com/doc/internals/en/client-server-protocol.html)

## Dependencies
1. [pubnative/mysqlproto-go](https://github.com/pubnative/mysqlproto-go) MySQL protocol implementation

## Installation
`go get github.com/pubnative/mysqldriver-go`

## Quick Start
```go
package main

import (
	"fmt"
	"strconv"

	"github.com/pubnative/mysqldriver-go"
)

type Person struct {
	Name    string
	Age     int
	Married bool
}

func main() {
	// initialize DB pool of 10 connections
	db := mysqldriver.NewDB("root@tcp(127.0.0.1:3306)/test", 10)

	// obtain connection from the pool
	conn, err := db.GetConn()
	if err != nil {
		panic(err)
	}

	if _, err := conn.Exec(`CREATE TABLE IF NOT EXISTS people (
        id int NOT NULL AUTO_INCREMENT,
    	name varchar(255),
    	age int,
        married tinyint,
        PRIMARY KEY (id)
    )`); err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		num := strconv.Itoa(i)
		_, err := conn.Exec(`
            INSERT INTO people(name,age,married) 
            VALUES("name` + num + `",` + num + `,` + strconv.Itoa(i%2) + `)
        `)
		if err != nil {
			panic(err)
		}
	}

	rows, err := conn.Query("SELECT name,age,married FROM people")
	if err != nil {
		panic(err)
	}

	for rows.Next() { // switch cursor to the next unread row
		person := Person{
			Name:    rows.String(),
			Age:     rows.Int(),
			Married: rows.Bool(),
		}
		fmt.Printf("%#v\n", person)
	}

	// always should be checked if there is an error during reading rows
	if err := rows.LastError(); err != nil {
		panic(err)
	}

	// return connection to the pool for further reuse
	if err := db.PutConn(conn); err != nil {
		panic(err)
	}

	if errors := db.Close(); errors != nil { // close the pool and all connections in it
	    for _, err := range errors {
	        _ = err // handle error
        }
	}
}
```
