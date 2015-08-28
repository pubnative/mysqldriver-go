package mysqlclient

import (
	"errors"
	"net"

	"github.com/pubnative/mysqlproto-go"
)

type Conn struct {
	conn  net.Conn
	proto mysqlproto.Proto
}

func NewConn(username, password, protocol, address, database string) (Conn, error) {
	conn, err := net.Dial(protocol, address)
	if err != nil {
		return Conn{}, err
	}

	proto := mysqlproto.NewProto()
	if err = handshake(proto, conn, username, password, database); err != nil {
		return Conn{}, err
	}

	return Conn{conn, proto}, nil
}

func (c Conn) Close() error {
	return c.conn.Close()
}

func handshake(proto mysqlproto.Proto, conn net.Conn, username, password, database string) error {
	packet, err := proto.NewHandshakeV10(conn)
	if err != nil {
		return err
	}

	flags := packet.CapabilityFlags
	flags &= ^mysqlproto.CLIENT_SSL
	flags &= ^mysqlproto.CLIENT_COMPRESS

	res := mysqlproto.HandshakeResponse41(
		packet.CapabilityFlags&(flags),
		packet.CharacterSet,
		username,
		password,
		packet.AuthPluginData,
		database,
		packet.AuthPluginName,
		nil,
	)

	if _, err := conn.Write(res); err != nil {
		return err
	}

	packetOK, err := proto.ReadPacket(conn)
	if err != nil {
		return err
	}

	if packetOK.Payload[0] != mysqlproto.PACKET_OK {
		return errors.New("Error occured during handshake with a server")
	}

	return nil
}
