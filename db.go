package mysqldriver

import (
	"errors"
	"strings"
)

var ErrClosedDB = errors.New("mysqldriver: can't get connection from the closed DB")

type DB struct {
	conns    chan Conn
	username string
	password string
	protocol string
	address  string
	database string
}

func NewDB(dataSource string, pool int) *DB {
	usr, pass, proto, addr, dbname := parseDataSource(dataSource)
	conns := make(chan Conn, pool)
	return &DB{
		conns:    conns,
		username: usr,
		password: pass,
		protocol: proto,
		address:  addr,
		database: dbname,
	}
}

func (db *DB) GetConn() (Conn, error) {
	select {
	case conn, more := <-db.conns:
		if !more {
			return Conn{}, ErrClosedDB
		}
		return conn, nil
	default:
		return db.dial()
	}
}

func (db *DB) PutConn(conn Conn) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = conn.Close()
			return
		}
	}()

	conn.conn.ResetStats()

	select {
	case db.conns <- conn:
	default:
		err = conn.Close()
	}

	return
}

func (db *DB) Close() {
	close(db.conns)
	for {
		conn, more := <-db.conns
		if more {
			conn.Close()
		} else {
			break
		}
	}
}

func (db *DB) dial() (Conn, error) {
	return NewConn(db.username, db.password, db.protocol, db.address, db.database)
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
