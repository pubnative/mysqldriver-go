// +build ignore

package main

import (
	"database/sql"
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/pubnative/mysqldriver-go"
)

func main() {
	debug.SetGCPercent(-1) // disable GC

	db := mysqldriver.NewDB("root@tcp(127.0.0.1:3306)/test", 10)
	sqlDB, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}

	preFillRecords(100)
	objectsInHEAP(func() string { return readAllMysqldriver(db) })
	objectsInHEAP(func() string { return readAllGoSqlDriver(sqlDB) })

	preFillRecords(1000)
	objectsInHEAP(func() string { return readAllMysqldriver(db) })
	objectsInHEAP(func() string { return readAllGoSqlDriver(sqlDB) })

	preFillRecords(10000)
	objectsInHEAP(func() string { return readAllMysqldriver(db) })
	objectsInHEAP(func() string { return readAllGoSqlDriver(sqlDB) })
}

func readAllMysqldriver(db *mysqldriver.DB) string {
	conn, err := db.GetConn()
	if err != nil {
		panic(err)
	}
	defer db.PutConn(conn)

	rows, err := conn.Query("SELECT name FROM mysqldriver_benchmarks")
	if err != nil {
		panic(err)
	}

	count := 0
	for rows.Next() {
		name := rows.String()
		count++
		_ = name
	}

	return "mysqldriver: records read " + strconv.Itoa(count)
}

func readAllGoSqlDriver(db *sql.DB) string {
	rows, err := db.Query("SELECT name FROM mysqldriver_benchmarks")
	if err != nil {
		panic(err)
	}

	count := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			panic(err)
		}
		count++
		_ = name
	}

	return "go-sql-driver: records read " + strconv.Itoa(count)
}

func objectsInHEAP(fn func() string) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	objects := memStats.HeapObjects
	now := time.Now()
	prefix := fn()
	took := time.Since(now)
	runtime.ReadMemStats(memStats)
	diff := memStats.HeapObjects - objects
	fmt.Println(prefix, " HEAP", diff, " time", took)
}

func preFillRecords(num int) {
	db := mysqldriver.NewDB("root@tcp(127.0.0.1:3306)/test", 10)
	conn, err := db.GetConn()
	if err != nil {
		panic(err)
	}
	if _, err := conn.Exec(`DROP TABLE IF EXISTS mysqldriver_benchmarks`); err != nil {
		panic(err)
	}
	if _, err := conn.Exec(`CREATE TABLE mysqldriver_benchmarks (
		id int NOT NULL AUTO_INCREMENT,
		name varchar(255),
		age int,
		PRIMARY KEY (id)
	)`); err != nil {
		panic(err)
	}

	for i := 0; i < num; i++ {
		_, err := conn.Exec(`INSERT INTO mysqldriver_benchmarks(name) VALUES("name` + strconv.Itoa(i) + `")`)
		if err != nil {
			panic(err)
		}
	}

	db.Close()
}
