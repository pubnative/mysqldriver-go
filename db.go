package mysqldriver

import (
	"errors"
	"strings"
)

var ErrClosedDB = errors.New("mysqldriver: can't get connection from the closed DB")

// DB manages pool of connection
type DB struct {
	OnDial func(conn *Conn) error // called when new connection is established

	conns    chan *Conn
	username string
	password string
	protocol string
	address  string
	database string
}

// NewDB initializes pool of connections but doesn't
// establishes connection to DB.
//
// Pool size is fixed and can't be resized later.
// DataSource parameter has the following format:
// [username[:password]@][protocol[(address)]]/dbname
func NewDB(dataSource string, pool int) *DB {
	usr, pass, proto, addr, dbname := parseDataSource(dataSource)
	conns := make(chan *Conn, pool)
	return &DB{
		conns:    conns,
		username: usr,
		password: pass,
		protocol: proto,
		address:  addr,
		database: dbname,
	}
}

// GetConn gets connection from the pool if there is one or
// establishes a new one.This method always returns the connection
// regardless the pool size. When DB is closed, this method
// returns ErrClosedDB error.
func (db *DB) GetConn() (*Conn, error) {
	select {
	case conn, more := <-db.conns:
		if !more {
			return nil, ErrClosedDB
		}
		return conn, nil
	default:
		return db.dial()
	}
}

// PutConn returns connection to the pool. When pool is reached,
// connection is closed and won't be further reused.
// If connection is already closed, PutConn will discard it
// so it's safe to return closed connection to the pool.
func (db *DB) PutConn(conn *Conn) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = conn.Close()
			return
		}
	}()

	if !conn.valid {
		// broken connection shouldn't be in a pool
		return conn.Close()
	}

	if conn.closed {
		return nil
	}

	conn.conn.ResetStats()

	select {
	case db.conns <- conn:
	default:
		err = conn.Close()
	}

	return
}

// Close closes all connections in a pool and
// doesn't allow to establish new ones to DB any more.
// Returns slice of errors if any occurred.
func (db *DB) Close() []error {
	close(db.conns)
	var errors []error
	for {
		conn, more := <-db.conns
		if more {
			if err := conn.Close(); err != nil {
				errors = append(errors, err)
			}
		} else {
			break
		}
	}
	return errors
}

func (db *DB) dial() (*Conn, error) {
	conn, err := NewConn(db.username, db.password, db.protocol, db.address, db.database)
	if err != nil {
		return conn, err
	}
	if db.OnDial != nil {
		err = db.OnDial(conn)
	}
	return conn, err
}

func parseDataSource(dataSource string) (username, password, protocol, address, database string) {
	params := strings.Split(dataSource, "@")

	userData := strings.Split(params[0], ":")
	serverData := strings.Split(params[1], "/")

	username = userData[0]
	if len(userData) > 1 {
		password = userData[1]
	}

	if len(serverData) > 1 {
		database = serverData[1]
	}

	protoHost := strings.Split(serverData[0], "(")
	protocol = protoHost[0]
	address = protoHost[1][:len(protoHost[1])-1]

	return
}
