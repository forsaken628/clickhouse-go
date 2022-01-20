package column

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/binary"
)

type Date32 struct {
	values Int32
}

func (dt *Date32) Type() Type {
	return "Date32"
}

func (col *Date32) ScanType() reflect.Type {
	return scanTypeTime
}

func (dt *Date32) Rows() int {
	return len(dt.values)
}

func (dt *Date32) Row(i int, ptr bool) interface{} {
	value := dt.row(i)
	if ptr {
		return &value
	}
	return value
}

func (dt *Date32) ScanRow(dest interface{}, row int) error {
	switch d := dest.(type) {
	case *time.Time:
		*d = dt.row(row)
	case **time.Time:
		*d = new(time.Time)
		**d = dt.row(row)
	default:
		return &ColumnConverterError{
			Op:   "ScanRow",
			To:   fmt.Sprintf("%T", dest),
			From: "Date32",
		}
	}
	return nil
}

func (dt *Date32) Append(v interface{}) (nulls []uint8, err error) {
	switch v := v.(type) {
	case []time.Time:
		in := make([]int32, 0, len(v))
		for _, t := range v {
			in = append(in, timeToInt32(t))
		}
		dt.values, nulls = append(dt.values, in...), make([]uint8, len(v))
	case []*time.Time:
		nulls = make([]uint8, len(v))
		for i, v := range v {
			switch {
			case v != nil:
				dt.values = append(dt.values, timeToInt32(*v))
			default:
				dt.values, nulls[i] = append(dt.values, 0), 1
			}
		}
	default:
		return nil, &ColumnConverterError{
			Op:   "Append",
			To:   "Date32",
			From: fmt.Sprintf("%T", v),
		}
	}
	return
}

func (dt *Date32) AppendRow(v interface{}) error {
	var date int32
	switch v := v.(type) {
	case time.Time:
		date = timeToInt32(v)
	case *time.Time:
		if v != nil {
			date = timeToInt32(*v)
		}
	case nil:
	default:
		return &ColumnConverterError{
			Op:   "AppendRow",
			To:   "Date32",
			From: fmt.Sprintf("%T", v),
		}
	}
	dt.values = append(dt.values, date)
	return nil
}

func (dt *Date32) Decode(decoder *binary.Decoder, rows int) error {
	return dt.values.Decode(decoder, rows)
}

func (dt *Date32) Encode(encoder *binary.Encoder) error {
	return dt.values.Encode(encoder)
}

func (dt *Date32) row(i int) time.Time {
	return time.Unix((int64(dt.values[i]) * secInDay), 0).UTC()
}

func timeToInt32(t time.Time) int32 {
	return int32(t.Unix() / secInDay)
}

var _ Interface = (*Date32)(nil)