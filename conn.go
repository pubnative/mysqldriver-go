package mysqldriver

import (
	"errors"
	"net"

	"github.com/pubnative/mysqlproto-go"
)

type Conn struct {
	stream *mysqlproto.Stream
}

type Stats struct {
	Syscalls int
}

func NewConn(username, password, protocol, address, database string) (Conn, error) {
	conn, err := net.Dial(protocol, address)
	if err != nil {
		return Conn{}, err
	}

	stream := mysqlproto.NewStream(conn)

	if err = handshake(stream, username, password, database); err != nil {
		return Conn{}, err
	}

	if err = setUTF8Charset(stream); err != nil {
		return Conn{}, err
	}

	return Conn{stream}, nil
}

func (c Conn) Close() error {
	return c.stream.Close()
}

func (c Conn) Stats() Stats {
	return Stats{
		Syscalls: c.stream.Syscalls(),
	}
}

func handshake(stream *mysqlproto.Stream, username, password, database string) error {
	packet, err := mysqlproto.ReadHandshakeV10(stream)
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

	if _, err := stream.Write(res); err != nil {
		return err
	}

	packetOK, err := stream.NextPacket()
	if err != nil {
		return err
	}

	if packetOK.Payload[0] != mysqlproto.PACKET_OK {
		return errors.New("Error occured during handshake with a server")
	}

	return nil
}

func setUTF8Charset(stream *mysqlproto.Stream) error {
	data := mysqlproto.ComQueryRequest([]byte("SET NAMES utf8"))
	if _, err := stream.Write(data); err != nil {
		return err
	}

	packetOK, err := stream.NextPacket()
	if err != nil {
		return err
	}

	if packetOK.Payload[0] != mysqlproto.PACKET_OK {
		return errors.New("Error occured during setting charset")
	}

	return nil
}
