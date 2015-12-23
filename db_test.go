package mysqldriver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBGetConnSuccessfullyEstablishConnection(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 1)
	conn, err := db.GetConn()
	assert.Nil(t, err)
	assert.True(t, conn.conn.CapabilityFlags>uint32(0))
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
