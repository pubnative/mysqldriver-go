package mysqldriver

import (
	"fmt"

	"github.com/pubnative/mysqlproto-go"
)

func handleOK(payload []byte, capabilityFlags uint32) error {
	if payload[0] == mysqlproto.PACKET_OK {
		return nil
	}

	if payload[0] == mysqlproto.PACKET_ERR {
		errPacket, err := mysqlproto.ParseERRPacket(payload, capabilityFlags)
		if err != nil {
			return err
		}
		return errPacket
	}

	return fmt.Errorf("mysqldriver: unknown error occured. Payload: %x", payload)
}
