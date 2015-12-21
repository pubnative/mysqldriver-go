package mysqldriver

import (
	"fmt"

	"github.com/pubnative/mysqlproto-go"
)

func handleOK(payload []byte, capabilityFlags uint32) error {
	if payload[0] == mysqlproto.OK_PACKET {
		return nil
	}

	if payload[0] == mysqlproto.ERR_PACKET {
		errPacket, err := mysqlproto.ParseERRPacket(payload, capabilityFlags)
		if err != nil {
			return err
		}
		return errPacket
	}

	return fmt.Errorf("mysqldriver: unknown error occured. Payload: %x", payload)
}
