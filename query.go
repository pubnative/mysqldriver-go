package mysqldriver

import (
	"strconv"

	"github.com/pubnative/mysqlproto-go"
)

// Rows represents result set of SELECT query
type Rows struct {
	resultSet mysqlproto.ResultSet
	packet    []byte
	offset    uint64
	eof       bool
	err       error
}

// Next moves cursor to the next unread row.
// It returns false when there are no more rows left
// or an error occurred during reading rows (see LastError() function)
// This function must be called before reading first row
// and continue being called until it returns false.
//  rows, _ := conn.Query("SELECT * FROM people LIMIT 2")
//  rows.Next() // move cursor to the first row
//  // read values from the first row
//  rows.Next() // move cursor to the second row
//  // read values from the second row
//  rows.Next() // drain the stream
// Best practice is to call Next() function in a loop:
//  rows, _ := conn.Query("SELECT * FROM people")
//  for rows.Next() {
//  	// read values from the row
//  }
// It's required to read all rows before performing another query
// because connection contains sequential stream of rows.
//  rows, _ := conn.Query("SELECT name FROM dogs LIMIT 1")
//  rows.Next()   // move cursor to the first row
//  rows.String() // dog's name
//  rows, _ = conn.Query("SELECT name FROM cats LIMIT 2")
//  rows.Next()   // move cursor to the second row of first query
//  rows.String() // still dog's name
//  rows.Next()   // returns false. closes the first stream of rows
//  rows.Next()   // move cursor to the first row of second query
//  rows.String() // cat's name
//  rows.Next()   // returns false. closes the second stream of rows
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

// Bytes returns value as slice of bytes.
// NULL value is represented as empty slice.
func (r *Rows) Bytes() []byte {
	value, _ := r.NullBytes()
	return value
}

// NullBytes returns value as a slice of bytes
// and NULL indicator. When value is NULL, second parameter is true.
func (r *Rows) NullBytes() ([]byte, bool) {
	value, offset, null := mysqlproto.ReadRowValue(r.packet, r.offset)
	r.offset = offset
	return value, null
}

// String returns value as a string.
// NULL value is represented as an empty string.
func (r *Rows) String() string {
	value, _ := r.NullString()
	return value
}

// NullString returns string as a value and
// NULL indicator. When value is NULL, second parameter is true.
func (r *Rows) NullString() (string, bool) {
	data, null := r.NullBytes()
	return string(data), null
}

// Int returns value as an int.
// NULL value is represented as 0.
// Int method uses strconv.Atoi to convert string into int.
// (see https://golang.org/pkg/strconv/#Atoi)
func (r *Rows) Int() int {
	num, _ := r.NullInt()
	return num
}

// NullInt returns value as an int and NULL indicator.
// When value is NULL, second parameter is true.
// NullInt method uses strconv.Atoi to convert string into int.
// (see https://golang.org/pkg/strconv/#Atoi)
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

// Int8 returns value as an int8.
// NULL value is represented as 0.
// Int8 method uses strconv.ParseInt to convert string into int8.
// (see https://golang.org/pkg/strconv/#ParseInt)
func (r *Rows) Int8() int8 {
	num, _ := r.NullInt8()
	return num
}

// NullInt8 returns value as an int8 and NULL indicator.
// When value is NULL, second parameter is true.
// NullInt8 method uses strconv.ParseInt to convert string into int8.
// (see https://golang.org/pkg/strconv/#ParseInt)
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

// Int16 returns value as an int16.
// NULL value is represented as 0.
// Int16 method uses strconv.ParseInt to convert string into int16.
// (see https://golang.org/pkg/strconv/#ParseInt)
func (r *Rows) Int16() int16 {
	num, _ := r.NullInt16()
	return num
}

// NullInt16 returns value as an int8 and NULL indicator.
// When value is NULL, second parameter is true.
// NullInt16 method uses strconv.ParseInt to convert string into int16.
// (see https://golang.org/pkg/strconv/#ParseInt)
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

// Int32 returns value as an int32.
// NULL value is represented as 0.
// Int32 method uses strconv.ParseInt to convert string into int32.
// (see https://golang.org/pkg/strconv/#ParseInt)
func (r *Rows) Int32() int32 {
	num, _ := r.NullInt32()
	return num
}

// NullInt32 returns value as an int32 and NULL indicator.
// When value is NULL, second parameter is true.
// NullInt32 method uses strconv.ParseInt to convert string into int32.
// (see https://golang.org/pkg/strconv/#ParseInt)
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

// Int64 returns value as an int64.
// NULL value is represented as 0.
// Int64 method uses strconv.ParseInt to convert string into int64.
// (see https://golang.org/pkg/strconv/#ParseInt)
func (r *Rows) Int64() int64 {
	num, _ := r.NullInt64()
	return num
}

// NullInt64 returns value as an int64 and NULL indicator.
// When value is NULL, second parameter is true.
// NullInt64 method uses strconv.ParseInt to convert string into int64.
// (see https://golang.org/pkg/strconv/#ParseInt)
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

// Float32 returns value as an float32.
// NULL value is represented as 0.0.
// Float32 method uses strconv.ParseFloat to convert string into float32.
// (see https://golang.org/pkg/strconv/#ParseFloat)
func (r *Rows) Float32() float32 {
	num, _ := r.NullFloat32()
	return num
}

// NullFloat32 returns value as an float32 and NULL indicator.
// When value is NULL, second parameter is true.
// NullFloat32 method uses strconv.ParseFloat to convert string into float32.
// (see https://golang.org/pkg/strconv/#ParseFloat)
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

// Float64 returns value as an float64.
// NULL value is represented as 0.0.
// Float64 method uses strconv.ParseFloat to convert string into float64.
// (see https://golang.org/pkg/strconv/#ParseFloat)
func (r *Rows) Float64() float64 {
	num, _ := r.NullFloat64()
	return num
}

// NullFloat64 returns value as an float64 and NULL indicator.
// When value is NULL, second parameter is true.
// NullFloat64 method uses strconv.ParseFloat to convert string into float64.
// (see https://golang.org/pkg/strconv/#ParseFloat)
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

// Bool returns value as an bool.
// NULL value is represented as false.
// Bool method uses strconv.ParseBool to convert string into bool.
// (see https://golang.org/pkg/strconv/#ParseBool)
func (r *Rows) Bool() bool {
	b, _ := r.NullBool()
	return b
}

// NullBool returns value as an bool and NULL indicator.
// When value is NULL, second parameter is true.
// NullBool method uses strconv.ParseBool to convert string into bool.
// (see https://golang.org/pkg/strconv/#ParseBool)
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

// LastError returns the error if any occurred during
// reading result set of SELECT query. This method should
// be always called after reading all rows.
//  rows, err := conn.Query("SELECT * FROM dogs")
//  if err != nil {
//  	// handle error
//  }
//  for rows.Next() {
//  	// read values
//  }
//  if err = rows.LastError(); err != nil {
//  	// handle error
//  }
func (r *Rows) LastError() error {
	return r.err
}

// Query function is used only for SELECT query.
// For all other queries and commands see func (c Conn) Exec
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

// Exec executes queries or other commands which expect to return OK_PACKET
// including INSERT/UPDATE/DELETE queries. For SELECT query see func (Conn) Query
//  okPacket, err := conn.Exec("DELETE FROM dogs WHERE id = 1")
//	if err == nil {
//  	return nil // query was performed successfully
//  }
//  if errPacket, ok := err.(mysqlproto.ERRPacket); ok {
//  	return errPacket // retrieve more information about the error
//  } else {
//  	return err // generic error
//  }
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
