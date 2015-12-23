package mysqldriver

import (
	"testing"

	"github.com/pubnative/mysqlproto-go"
	"github.com/stretchr/testify/assert"
)

func TestDBGetConnSuccessfullyEstablishConnection(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 1)
	conn, err := db.GetConn()
	assert.Nil(t, err)
	assert.True(t, conn.conn.CapabilityFlags > uint32(0))
}

func TestDBGetConnReturnsConnectionFromThePool(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 2)
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
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 2)
	db.Close()
	_, err := db.GetConn()
	assert.Equal(t, err, ErrClosedDB)
}

func TestDBPutConnAddsConnectionToThePool(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 2)
	assert.Len(t, db.conns, 0)
	conn, _ := db.GetConn()
	assert.Nil(t, db.PutConn(conn))
	assert.Len(t, db.conns, 1)
}

func TestDBPutConnAddsUpTpPoolSize(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 2)
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

func TestDBPutConnClosesConnectionWhenDBIsClosed(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 2)
	db.Close()

	s := &stream{}
	conn := Conn{mysqlproto.Conn{mysqlproto.NewStream(s), 0}}
	db.PutConn(conn)
	assert.True(t, s.closed)
}

func TestDBCloseClosesAllConnections(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 2)
	s1 := &stream{}
	conn1 := Conn{mysqlproto.Conn{mysqlproto.NewStream(s1), 0}}
	db.PutConn(conn1)
	s2 := &stream{}
	conn2 := Conn{mysqlproto.Conn{mysqlproto.NewStream(s2), 0}}
	db.PutConn(conn2)

	assert.Len(t, db.conns, 2)
	assert.False(t, s1.closed)
	assert.False(t, s2.closed)
	db.Close()
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
func (s *stream) Read([]byte) (int, error)  { return 0, nil }
func (s *stream) Close() error              { s.closed = true; return nil }
