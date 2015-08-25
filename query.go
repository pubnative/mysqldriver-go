package mysqlclient

import (
	"github.com/pubnative/mysqlproto-go"
)

type Rows struct {
	resultSet mysqlproto.ResultSet
	packet []byte
	offset uint64
	eof bool
	err error
}

func (r *Rows) Next() bool {
	if r.eof {
		return false
	}

	packet, err := r.resultSet.Row()
	if err != nil {
		r.err = err
		r.eof = true
		return false
	}

	if packet == nil {
		r.eof = true
		return false
	} else {
		r.packet = packet
		r.offset = 0
		return true
	}
}

func (r *Rows) String() string {
	value, offset := mysqlproto.ParseString(r.packet, r.offset)
	r.offset = offset
	return value
}

func (r *Rows) LastError() error {
	return r.err
}

func (c Conn) Query(sql string) (*Rows, error) {
	req := mysqlproto.ComQueryRequest([]byte(sql))
	if _, err := c.conn.Write(req); err != nil {
		return nil, err
	}

	resultSet, err := mysqlproto.ComQueryResponse(c.conn)
	if err != nil {
		return nil, err
	}

	return &Rows{resultSet: resultSet}, nil
}
