package mysqldriver

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
	value, offset, _ := mysqlproto.ReadRowValue(r.packet, r.offset)
	r.offset = offset
	return value
}

func (r *Rows) NullBytes() ([]byte, bool) {
	value, offset, null := mysqlproto.ReadRowValue(r.packet, r.offset)
	r.offset = offset
	return value, null
}

func (r *Rows) String() string {
	return string(r.Bytes())
}

func (r *Rows) NullString() (string, bool) {
	data, null := r.NullBytes()
	return string(data), null
}

func (r *Rows) Int() int {
	num, _ := r.NullInt()
	return num
}

func (r *Rows) NullInt() (int, bool) {
	str, null := r.NullString()
	if null {
		return 0, true
	}

	num, err := strconv.Atoi(str)
	if err != nil {
		r.err = err
	}

	return num, false
}

func (r *Rows) Int8() int8 {
	num, _ := r.NullInt8()
	return num
}

func (r *Rows) NullInt8() (int8, bool) {
	str, null := r.NullString()
	if null {
		return 0, true
	}

	num, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		r.err = err
	}

	return int8(num), false
}

func (r *Rows) Int16() int16 {
	num, _ := r.NullInt16()
	return num
}

func (r *Rows) NullInt16() (int16, bool) {
	str, null := r.NullString()
	if null {
		return 0, true
	}

	num, err := strconv.ParseInt(str, 10, 16)
	if err != nil {
		r.err = err
	}

	return int16(num), false
}

func (r *Rows) Int32() int32 {
	num, _ := r.NullInt32()
	return num
}

func (r *Rows) NullInt32() (int32, bool) {
	str, null := r.NullString()
	if null {
		return 0, true
	}

	num, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		r.err = err
	}

	return int32(num), false
}

func (r *Rows) Int64() int64 {
	num, _ := r.NullInt64()
	return num
}

func (r *Rows) NullInt64() (int64, bool) {
	str, null := r.NullString()
	if null {
		return 0, true
	}

	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		r.err = err
	}

	return int64(num), false
}

func (r *Rows) Float32() float32 {
	num, _ := r.NullFloat32()
	return num
}

func (r *Rows) NullFloat32() (float32, bool) {
	str, null := r.NullString()
	if null {
		return 0, true
	}

	num, err := strconv.ParseFloat(str, 32)
	if err != nil {
		r.err = err
	}

	return float32(num), false
}

func (r *Rows) Float64() float64 {
	num, _ := r.NullFloat64()
	return num
}

func (r *Rows) NullFloat64() (float64, bool) {
	str, null := r.NullString()
	if null {
		return 0, true
	}

	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		r.err = err
	}

	return num, false
}

func (r *Rows) Bool() bool {
	b, _ := r.NullBool()
	return b
}

func (r *Rows) NullBool() (bool, bool) {
	str, null := r.NullString()
	if null {
		return false, true
	}

	b, err := strconv.ParseBool(str)
	if err != nil {
		r.err = err
	}
	return b, false
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

func (c Conn) Exec(sql string) (mysqlproto.OKPacket, error) {
	req := mysqlproto.ComQueryRequest([]byte(sql))
	if _, err := c.conn.Write(req); err != nil {
		return mysqlproto.OKPacket{}, err
	}

	packet, err := c.conn.NextPacket()
	if err != nil {
		return mysqlproto.OKPacket{}, err
	}

	if packet.Payload[0] == mysqlproto.OK_PACKET {
		pkt, err := mysqlproto.ParseOKPacket(packet.Payload, c.conn.CapabilityFlags)
		return pkt, err
	} else {
		pkt, err := mysqlproto.ParseERRPacket(packet.Payload, c.conn.CapabilityFlags)
		if err == nil {
			return mysqlproto.OKPacket{}, pkt
		} else {
			return mysqlproto.OKPacket{}, err
		}
	}
}
