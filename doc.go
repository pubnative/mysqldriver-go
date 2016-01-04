/*
Package mysqldriver is a GC optimized MySQL driver

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

Reading rows

mysqldriver reads data from the DB in a sequential order
which means the whole result set of first query must be read
before executing another one.

Number of read column's values and their types must match
with the number of columns in a query.

 rows, err := conn.Query("SELECT id, name, married FROM people")
 if err != nil {
 	// handle error
 }
 for rows.Next() { // always read all rows
 	id := rows.Int()       // order of columns must be preserved
 	name := rows.String()  // type of the column must match with DB type
 	married := rows.Bool() // all column's values must be read
 }
 if err = rows.LastError(); err != nil {
 	// handle error if any occurred during reading packets from DB
 }
*/
package mysqldriver
