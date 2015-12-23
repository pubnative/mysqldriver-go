package mysqldriver

import (
	"testing"

	"github.com/pubnative/mysqlproto-go"
	"github.com/stretchr/testify/assert"
)

func TestHandleOKNoError(t *testing.T) {
	data := []byte{0x00}
	err := handleOK(data, 0)
	assert.Nil(t, err)
}

func TestHandleOKPacketContainsError(t *testing.T) {
	data := []byte{
		0xff, 0x48, 0x04, 0x23, 0x48, 0x59,
		0x30, 0x30, 0x30, 0x4e, 0x6f, 0x20,
		0x74, 0x61, 0x62, 0x6c, 0x65, 0x73,
		0x20, 0x75, 0x73, 0x65, 0x64,
	}
	err := handleOK(data, 0)
	assert.NotNil(t, err)
	_, ok := err.(mysqlproto.ERRPacket)
	assert.True(t, ok)
}

func TestHandleOKPacketContainsBrokenErrorPayout(t *testing.T) {
	data := []byte{0x01}
	err := handleOK(data, 0)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "mysqldriver: unknown error occured. Payload: 01")
}
