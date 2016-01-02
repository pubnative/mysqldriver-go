# mysqldriver-go
[![Build Status](https://travis-ci.org/pubnative/mysqldriver-go.svg?branch=master)](https://travis-ci.org/pubnative/mysqldriver-go)
[![GoDoc](https://godoc.org/github.com/pubnative/mysqldriver-go?status.svg)](https://godoc.org/github.com/pubnative/mysqldriver-go)

## Motivation
There are already many MySQL drivers which implement [database/sql](https://golang.org/pkg/database/sql/) interface however using this generic interface, especialy [Scan](https://golang.org/pkg/database/sql/#Row.Scan) method, requires to store many objects in a HEAP. Reading massive number of records from DB can significantly increase GC pause time which is very sensitive for low-latency applications. Due to this issue, was made a decision to write another MySQL driver which is GC friendly as much as possible, and not to follow [database/sql](https://golang.org/pkg/database/sql/) interface.

[Benchmark](https://github.com/pubnative/mysqldriver-go/blob/master/benchmarks/main.go) was performed on MacBook Pro (Retina, 13-inch, Late 2013), 2.8 GHz Intel Core i7, 16 GB 1600 MHz DDR3
[![comparison](https://cloud.githubusercontent.com/assets/296795/12074709/9dbf19a2-b162-11e5-8dd0-a973b57895b0.png)](https://jsfiddle.net/zs83oze6/1/)
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
Main goals of this library are: *performance* over flexibility, *simplicity* over complexity. Any new feature shouldn't decrease performance of exising code, any improvements to productivity are always welcome. There are no plan to convert this library into ORM, it should stay simple however support all MySQL features.

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

	db.Close() // close the pool and all connections in it
}
```
