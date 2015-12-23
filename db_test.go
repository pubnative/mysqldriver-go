package mysqldriver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
