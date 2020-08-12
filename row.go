package mysqldriver

import (
	"strconv"
)

// Row reads the entire row.
// This function is identical to read each column successively.
//  rows, _ := conn.Query("SELECT id, name FROM people")
//  for rows.Next() {
//		row := rows.Row() // reads both columns
//		// the same as reading each column separately
//  	// rows.Bytes()
//  	// rows.Bytes()
//  }
// It is possible to call Row() function after reading some of the columns.
// In this case, Row() will read the rest of the columns of the row.
// 	rows, _ := conn.Query("SELECT id, name, age FROM people")
//  for rows.Next() {
// 		rows.Int() // read ID
//		row := rows.Row() // reads name, age and return the full row
//		fmt.Println(row.Int("id"), row.String("name"), row.Int("age"))
//  }
func (r *Rows) Row() Row {
	for range r.resultSet.Columns[r.readColumns:] {
		r.NullBytes()
	}

	row := Row{
		rows:    r,
		columns: r.columns,
	}
	return row
}

// Row represents a single DB row of the query results
type Row struct {
	rows    *Rows // used to set an error
	columns map[string]columnValue
}

// NullBytes returns value as a slice of bytes
// and NULL indicator. When value is NULL, second parameter is true.
//
// IMPORTANT. This function panics if it can't find the column by the name.
//
// All other type-specific functions are based on this one.
func (r Row) NullBytes(col string) ([]byte, bool) {
	column, ok := r.columns[col]
	if !ok {
		msg := `mysqldriver: column "` + col + `" doesn't exist.`
		if len(r.columns) > 0 {
			msg += ` Available columns are: `
			var i int
			for _, c := range r.rows.resultSet.Columns {
				if i > 0 {
					msg += ", "
				}
				msg += `"` + c.Name + `"`
				i += 1
			}
		}
		panic(msg)
	}

	return column.data, column.null
}

// Bytes returns value as slice of bytes.
// NULL value is represented as empty slice.
func (r Row) Bytes(col string) []byte {
	value, _ := r.NullBytes(col)
	return value
}

// NullString returns string as a value and
// NULL indicator. When value is NULL, second parameter is true.
func (r Row) NullString(col string) (string, bool) {
	value, null := r.NullBytes(col)
	return string(value), null
}

// String returns value as a string.
// NULL value is represented as an empty string.
func (r Row) String(col string) string {
	value, _ := r.NullString(col)
	return value
}

// NullInt returns value as an int and NULL indicator.
// When value is NULL, second parameter is true.
// NullInt method uses strconv.Atoi to convert string into int.
// (see https://golang.org/pkg/strconv/#Atoi)
func (r Row) NullInt(col string) (int, bool) {
	value, null := r.NullBytes(col)
	if null {
		return 0, true
	}

	num, err := atoi(value)
	if err != nil {
		r.rows.errParse = err
	}

	return num, false
}

// Int returns value as an int.
// NULL value is represented as 0.
// Int method uses strconv.Atoi to convert string into int.
// (see https://golang.org/pkg/strconv/#Atoi)
func (r Row) Int(col string) int {
	value, _ := r.NullInt(col)
	return value
}

// NullInt8 returns value as an int8 and NULL indicator.
// When value is NULL, second parameter is true.
// NullInt8 method uses strconv.ParseInt to convert string into int8.
// (see https://golang.org/pkg/strconv/#ParseInt)
func (r Row) NullInt8(col string) (int8, bool) {
	str, null := r.NullString(col)
	if null {
		return 0, true
	}

	num, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		r.rows.errParse = err
	}

	return int8(num), false
}

// Int8 returns value as an int8.
// NULL value is represented as 0.
// Int8 method uses strconv.ParseInt to convert string into int8.
// (see https://golang.org/pkg/strconv/#ParseInt)
func (r Row) Int8(col string) int8 {
	num, _ := r.NullInt8(col)
	return num
}

// NullInt16 returns value as an int8 and NULL indicator.
// When value is NULL, second parameter is true.
// NullInt16 method uses strconv.ParseInt to convert string into int16.
// (see https://golang.org/pkg/strconv/#ParseInt)
func (r Row) NullInt16(col string) (int16, bool) {
	str, null := r.NullString(col)
	if null {
		return 0, true
	}

	num, err := strconv.ParseInt(str, 10, 16)
	if err != nil {
		r.rows.errParse = err
	}

	return int16(num), false
}

// Int16 returns value as an int16.
// NULL value is represented as 0.
// Int16 method uses strconv.ParseInt to convert string into int16.
// (see https://golang.org/pkg/strconv/#ParseInt)
func (r Row) Int16(col string) int16 {
	num, _ := r.NullInt16(col)
	return num
}

// NullInt32 returns value as an int32 and NULL indicator.
// When value is NULL, second parameter is true.
// NullInt32 method uses strconv.ParseInt to convert string into int32.
// (see https://golang.org/pkg/strconv/#ParseInt)
func (r Row) NullInt32(col string) (int32, bool) {
	str, null := r.NullString(col)
	if null {
		return 0, true
	}

	num, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		r.rows.errParse = err
	}

	return int32(num), false
}

// Int32 returns value as an int32.
// NULL value is represented as 0.
// Int32 method uses strconv.ParseInt to convert string into int32.
// (see https://golang.org/pkg/strconv/#ParseInt)
func (r Row) Int32(col string) int32 {
	num, _ := r.NullInt32(col)
	return num
}

// NullInt64 returns value as an int64 and NULL indicator.
// When value is NULL, second parameter is true.
// NullInt64 method uses strconv.ParseInt to convert string into int64.
// (see https://golang.org/pkg/strconv/#ParseInt)
func (r Row) NullInt64(col string) (int64, bool) {
	str, null := r.NullString(col)
	if null {
		return 0, true
	}

	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		r.rows.errParse = err
	}

	return int64(num), false
}

// Int64 returns value as an int64.
// NULL value is represented as 0.
// Int64 method uses strconv.ParseInt to convert string into int64.
// (see https://golang.org/pkg/strconv/#ParseInt)
func (r Row) Int64(col string) int64 {
	num, _ := r.NullInt64(col)
	return num
}

// NullFloat32 returns value as a float32 and NULL indicator.
// When value is NULL, second parameter is true.
// NullFloat32 method uses strconv.ParseFloat to convert string into float32.
// (see https://golang.org/pkg/strconv/#ParseFloat)
func (r Row) NullFloat32(col string) (float32, bool) {
	str, null := r.NullString(col)
	if null {
		return 0, true
	}

	num, err := strconv.ParseFloat(str, 32)
	if err != nil {
		r.rows.errParse = err
	}

	return float32(num), false
}

// Float32 returns value as a float32.
// NULL value is represented as 0.0.
// Float32 method uses strconv.ParseFloat to convert string into float32.
// (see https://golang.org/pkg/strconv/#ParseFloat)
func (r Row) Float32(col string) float32 {
	num, _ := r.NullFloat32(col)
	return num
}

// NullFloat64 returns value as a float64 and NULL indicator.
// When value is NULL, second parameter is true.
// NullFloat64 method uses strconv.ParseFloat to convert string into float64.
// (see https://golang.org/pkg/strconv/#ParseFloat)
func (r Row) NullFloat64(col string) (float64, bool) {
	str, null := r.NullString(col)
	if null {
		return 0, true
	}

	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		r.rows.errParse = err
	}

	return num, false
}

// Float64 returns value as a float64.
// NULL value is represented as 0.0.
// Float64 method uses strconv.ParseFloat to convert string into float64.
// (see https://golang.org/pkg/strconv/#ParseFloat)
func (r Row) Float64(col string) float64 {
	num, _ := r.NullFloat64(col)
	return num
}

// NullBool returns value as a bool and NULL indicator.
// When value is NULL, second parameter is true.
// NullBool method uses strconv.ParseBool to convert string into bool.
// (see https://golang.org/pkg/strconv/#ParseBool)
func (r Row) NullBool(col string) (bool, bool) {
	str, null := r.NullBytes(col)
	if null {
		return false, true
	}

	b, err := parseBool(str)
	if err != nil {
		r.rows.errParse = err
	}
	return b, false
}

// Bool returns value as a bool.
// NULL value is represented as false.
// Bool method uses strconv.ParseBool to convert string into bool.
// (see https://golang.org/pkg/strconv/#ParseBool)
func (r Row) Bool(col string) bool {
	b, _ := r.NullBool(col)
	return b
}
