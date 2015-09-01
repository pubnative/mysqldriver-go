package mysqlclient

import (
	"errors"
	"net"

	"github.com/pubnative/mysqlproto-go"
)

type Conn struct {
	conn net.Conn
}

func NewConn(username, password, protocol, address, database string) (Conn, error) {
	conn, err := net.Dial(protocol, address)
	if err != nil {
		return Conn{}, err
	}

	if err = handshake(conn, username, password, database); err != nil {
		return Conn{}, err
	}

	if err = setUTF8Charset(conn); err != nil {
		return Conn{}, err
	}

	return Conn{conn}, nil
}

func (c Conn) Close() error {
	return c.conn.Close()
}

func handshake(conn net.Conn, username, password, database string) error {
	packet, err := mysqlproto.ReadHandshakeV10(conn)
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

	streamPkt := mysqlproto.NewStreamPacket(conn)
	packetOK, err := streamPkt.NextPacket()
	if err != nil {
		return err
	}

	if packetOK.Payload[0] != mysqlproto.PACKET_OK {
		return errors.New("Error occured during handshake with a server")
	}

	return nil
}

func setUTF8Charset(conn net.Conn) error {
	data := mysqlproto.ComQueryRequest([]byte("SET NAMES utf8"))
	if _, err := conn.Write(data); err != nil {
		return err
	}

	streamPkt := mysqlproto.NewStreamPacket(conn)
	packetOK, err := streamPkt.NextPacket()
	if err != nil {
		return err
	}

	if packetOK.Payload[0] != mysqlproto.PACKET_OK {
		return errors.New("Error occured during setting charset")
	}

	return nil
}
