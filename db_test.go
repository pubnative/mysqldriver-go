package mysqldriver

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/pubnative/mysqlproto-go"
	"github.com/stretchr/testify/assert"
)

func TestDBGetConnSuccessfullyEstablishConnection(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 1, time.Duration(0))
	conn, err := db.GetConn()
	assert.Nil(t, err)
	assert.True(t, conn.conn.CapabilityFlags > uint32(0))
}

func TestDBGetConnReturnsConnectionFromThePool(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 2, time.Duration(0))
	conn1, _ := db.GetConn()
	conn2, _ := db.GetConn()
	db.PutConn(conn1)
	db.PutConn(conn2)

	assert.Len(t, db.conns, 2)
	db.GetConn()
	assert.Len(t, db.conns, 1)
	db.GetConn()
	assert.Len(t, db.conns, 0)
}

func TestDBGetConnReturnsErrorWhenDBIsClosed(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 2, time.Duration(0))
	errors := db.Close()
	assert.Nil(t, errors)
	_, err := db.GetConn()
	assert.Equal(t, err, ErrClosedDB)
}

func TestDBPutConnAddsConnectionToThePool(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 2, time.Duration(0))
	assert.Len(t, db.conns, 0)
	conn, _ := db.GetConn()
	assert.Nil(t, db.PutConn(conn))
	assert.Len(t, db.conns, 1)
}

func TestDBPutConnAddsUpToPoolSize(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 2, time.Duration(0))
	conn1, _ := db.GetConn()
	conn2, _ := db.GetConn()
	conn3, _ := db.GetConn()
	assert.Len(t, db.conns, 0)
	assert.Nil(t, db.PutConn(conn1))
	assert.Len(t, db.conns, 1)
	assert.Nil(t, db.PutConn(conn2))
	assert.Len(t, db.conns, 2)
	assert.Nil(t, db.PutConn(conn3))
	assert.Len(t, db.conns, 2)
}

func TestDBPutConnClosesConnectionWhenItIsInvalid(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 2, time.Duration(0))
	errors := db.Close()
	assert.Nil(t, errors)

	s := &stream{}
	conn := &Conn{mysqlproto.Conn{mysqlproto.NewStream(s, time.Duration(0)), 0}, false, false}
	db.PutConn(conn)
	assert.True(t, s.closed)
	assert.Len(t, db.conns, 0)
}

func TestDBPutConnDiscardsConnectionWhenItIsClosedAlready(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 2, time.Duration(0))
	errors := db.Close()
	assert.Nil(t, errors)

	conn := &Conn{conn: mysqlproto.Conn{nil, 0}, valid: true, closed: true}
	assert.Nil(t, db.PutConn(conn))
	assert.Len(t, db.conns, 0)
}

func TestDBPutConnClosesConnectionWhenDBIsClosed(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 2, time.Duration(0))
	errors := db.Close()
	assert.Nil(t, errors)

	s := &stream{}
	conn := &Conn{mysqlproto.Conn{mysqlproto.NewStream(s, time.Duration(0)), 0}, true, false}
	db.PutConn(conn)
	assert.True(t, s.closed)
	assert.Len(t, db.conns, 0)
}

func TestDBCloseClosesAllConnections(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 2, time.Duration(0))
	s1 := &stream{}
	conn1 := &Conn{mysqlproto.Conn{mysqlproto.NewStream(s1, time.Duration(0)), 0}, true, false}
	db.PutConn(conn1)
	s2 := &stream{}
	conn2 := &Conn{mysqlproto.Conn{mysqlproto.NewStream(s2, time.Duration(0)), 0}, true, false}
	db.PutConn(conn2)

	assert.Len(t, db.conns, 2)
	assert.False(t, s1.closed)
	assert.False(t, s2.closed)
	errors := db.Close()
	assert.Nil(t, errors)
	assert.True(t, s1.closed)
	assert.True(t, s2.closed)
	_, more := <-db.conns
	assert.False(t, more)
}

func TestParseDataSourceFull(t *testing.T) {
	source := "root:123@tcp(127.0.0.1:3306)/test"
	usr, pass, proto, addr, dbname := parseDataSource(source)
	assert.Equal(t, usr, "root")
	assert.Equal(t, pass, "123")
	assert.Equal(t, proto, "tcp")
	assert.Equal(t, addr, "127.0.0.1:3306")
	assert.Equal(t, dbname, "test")
}

func TestParseDataSourceWithoutPassword(t *testing.T) {
	source := "root@tcp(127.0.0.1:3306)/test"
	usr, pass, proto, addr, dbname := parseDataSource(source)
	assert.Equal(t, usr, "root")
	assert.Equal(t, pass, "")
	assert.Equal(t, proto, "tcp")
	assert.Equal(t, addr, "127.0.0.1:3306")
	assert.Equal(t, dbname, "test")
}

func TestParseDataSourceWithoutDatabase(t *testing.T) {
	source := "root@tcp(127.0.0.1:3306)"
	usr, pass, proto, addr, dbname := parseDataSource(source)
	assert.Equal(t, usr, "root")
	assert.Equal(t, pass, "")
	assert.Equal(t, proto, "tcp")
	assert.Equal(t, addr, "127.0.0.1:3306")
	assert.Equal(t, dbname, "")
}

type stream struct{ closed bool }

func (s *stream) Write([]byte) (int, error) { return 0, nil }
func (s *stream) Read([]byte) (int, error)  { return 0, io.EOF }
func (s *stream) Close() error              { s.closed = true; return nil }
func (s *stream) RemoteAddr() net.Addr { return MockAddr{} }
func (s *stream) LocalAddr() net.Addr { return MockAddr{} }
func (s *stream) SetDeadline(t time.Time) error { return nil}
func (s *stream) SetReadDeadline(t time.Time) error { return nil}
func (s *stream) SetWriteDeadline(t time.Time) error { return nil}

type MockAddr struct {}
func (m MockAddr) Network() string { return "" }
func (m MockAddr) String() string { return "" }
// Initializes the pool for 10 connections
func ExampleNewDB() {
	NewDB("root@tcp(127.0.0.1:3306)/test", 10, time.Duration(0))
}
