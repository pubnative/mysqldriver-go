package mysqlclient

import (
	"strconv"

	"github.com/pubnative/mysqlproto-go"
)

type Rows struct {
	resultSet mysqlproto.ResultSet
	packet    []byte
	offset    uint64
	eof       bool
	err       error
}

func (r *Rows) Next() bool {
	if r.eof {
		return false
	}

	if r.err != nil {
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

func (r *Rows) Bytes() []byte {
	value, offset := mysqlproto.ReadRowValue(r.packet, r.offset)
	r.offset = offset
	return value
}

func (r *Rows) String() string {
	return string(r.Bytes())
}

func (r *Rows) Int() int {
	num, err := strconv.Atoi(r.String())
	if err != nil {
		r.err = err
	}

	return num
}

func (r *Rows) Int8() int8 {
	num, err := strconv.ParseInt(r.String(), 10, 8)
	if err != nil {
		r.err = err
	}
	return int8(num)
}

func (r *Rows) Int16() int16 {
	num, err := strconv.ParseInt(r.String(), 10, 16)
	if err != nil {
		r.err = err
	}
	return int16(num)
}

func (r *Rows) Int32() int32 {
	num, err := strconv.ParseInt(r.String(), 10, 32)
	if err != nil {
		r.err = err
	}
	return int32(num)
}

func (r *Rows) Int64() int64 {
	num, err := strconv.ParseInt(r.String(), 10, 64)
	if err != nil {
		r.err = err
	}
	return int64(num)
}

func (r *Rows) Float32() float32 {
	num, err := strconv.ParseFloat(r.String(), 32)
	if err != nil {
		r.err = err
	}
	return float32(num)
}

func (r *Rows) Float64() float64 {
	num, err := strconv.ParseFloat(r.String(), 64)
	if err != nil {
		r.err = err
	}
	return num
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
