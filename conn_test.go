package mysqldriver

import (
	"context"
	"testing"

	"github.com/pubnative/mysqlproto-go"
	"github.com/stretchr/testify/assert"
)

func TestNewConnSuccess(t *testing.T) {
	conn, err := NewConn("root", "", "tcp", "127.0.0.1:3306", "test")
	assert.Nil(t, err)
	assert.True(t, conn.valid)
}

func TestNewConnError(t *testing.T) {
	conn, err := NewConn("root", "", "tcp", "127.0.0.1:3306", "unknown")
	assert.NotNil(t, err)
	errPkt, ok := err.(mysqlproto.ERRPacket)
	assert.True(t, ok)
	assert.Equal(t, errPkt.ErrorCode, uint16(1049))
	assert.Equal(t, errPkt.SQLState, "42000")
	assert.Equal(t, errPkt.ErrorMessage, "Unknown database 'unknown'")
	assert.False(t, conn.valid)
}

func TestNewConnContextSuccess(t *testing.T) {
	conn, err := NewConnContext(context.Background(), "root", "", "tcp", "127.0.0.1:3306", "test")
	assert.NoError(t, err)
	assert.True(t, conn.valid)
}

func TestNewConnContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := NewConnContext(ctx, "root", "", "tcp", "127.0.0.1:3306", "test")
	assert.EqualError(t, err, "dial tcp 127.0.0.1:3306: operation was canceled")
}

func TestConnClose(t *testing.T) {
	conn, err := NewConn("root", "", "tcp", "127.0.0.1:3306", "test")
	assert.Nil(t, err)
	assert.Nil(t, conn.Close())
	assert.True(t, conn.closed)
}
