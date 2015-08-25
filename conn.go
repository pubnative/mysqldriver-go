package mysqlclient

import (
	"errors"
	"net"
	"strings"

	"github.com/pubnative/mysqlproto-go"
)

type Conn struct {
	conn net.Conn
}

func NewConn(dataSource string) (Conn, error) {
	usr, pass, proto, addr, dbname := parseDataSource(dataSource)

	conn, err := net.Dial(proto, addr)
	if err != nil {
		return Conn{}, err
	}

	if err = handshake(conn, usr, pass, dbname); err != nil {
		return Conn{}, err
	}

	return Conn{conn}, nil
}

func parseDataSource(dataSource string) (username, password, protocol, address, database string) {
	params := strings.Split(dataSource, "@")

	userData := strings.Split(params[0], ":")
	serverData := strings.Split(params[1], "/")

	username = userData[0]
	if len(userData) > 1 {
		password = userData[1]
	}

	if len(serverData) > 1 {
		database = serverData[1]
	}

	protoHost := strings.Split(serverData[0], "(")
	protocol = protoHost[0]
	address = protoHost[1][:len(protoHost[1])-1]

	return
}

func handshake(conn net.Conn, username, password, database string) error {
	packet, err := mysqlproto.NewHandshakeV10(conn)
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

	packetOK, err := mysqlproto.ReadPacket(conn)
	if err != nil {
		return err
	}

	if packetOK.Payload[0] != mysqlproto.PACKET_OK {
		return errors.New("Error occured during handshake with a server")
	}

	return nil
}
