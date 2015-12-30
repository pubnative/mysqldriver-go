/*
Package mysqldriver is a driver for MySQL database

Concurrency

DB struct manages pool of connections to MySQL. Connection itself
isn't thread-safe, so it should be obtained per every go-routine.
It's important to return a connection back to the pool
when it's not needed for further reuse.

 db := mysqldriver.NewDB("root@tcp(127.0.0.1:3306)/test", 10)
 for i := 0; i < 10; i++ {
 	go func() {
 		conn, err := db.GetConn()
 		if err != nil {
 			// handle error
 		}
 		defer db.PutConn(conn) // return connection to the pool
 		// perform queries
 	}()
 }
*/
package mysqldriver
