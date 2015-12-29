package mysqldriver

import (
	"net"

	"github.com/pubnative/mysqlproto-go"
)

var capabilityFlags = mysqlproto.CLIENT_LONG_PASSWORD |
	mysqlproto.CLIENT_FOUND_ROWS |
	mysqlproto.CLIENT_LONG_FLAG |
	mysqlproto.CLIENT_CONNECT_WITH_DB |
	mysqlproto.CLIENT_PLUGIN_AUTH |
	mysqlproto.CLIENT_TRANSACTIONS |
	mysqlproto.CLIENT_PROTOCOL_41 |
	mysqlproto.CLIENT_SESSION_TRACK

// Conn represents connection to MySQL server
type Conn struct {
	conn mysqlproto.Conn
}

// Contains connection statistics
type Stats struct {
	Syscalls int // number of system calls performed to read all packets
}

// NewConn establishes connection to the DB. After obtaining the connection,
// it sends "SET NAMES utf8" command to the DB
func NewConn(username, password, protocol, address, database string) (Conn, error) {
	conn, err := net.Dial(protocol, address)
	if err != nil {
		return Conn{}, err
	}

	stream, err := mysqlproto.ConnectPlainHandshake(
		conn, capabilityFlags,
		username, password, database, nil,
	)

	if err != nil {
		return Conn{}, err
	}

	if err = setUTF8Charset(stream); err != nil {
		return Conn{}, err
	}

	return Conn{stream}, nil
}

// Close closes the connection
func (c Conn) Close() error {
	return c.conn.Close()
}

// Returns statistics about the connection
func (c Conn) Stats() Stats {
	return Stats{
		Syscalls: c.conn.Syscalls(),
	}
}

// Add sum ups all stats
func (s Stats) Add(stats Stats) Stats {
	return Stats{
		Syscalls: s.Syscalls + stats.Syscalls,
	}
}

func setUTF8Charset(conn mysqlproto.Conn) error {
	data := mysqlproto.ComQueryRequest([]byte("SET NAMES utf8"))
	if _, err := conn.Write(data); err != nil {
		return err
	}

	packet, err := conn.NextPacket()
	if err != nil {
		return err
	}

	return handleOK(packet.Payload, conn.CapabilityFlags)
}
