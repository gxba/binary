//bench test case for binary and std

package bench

import (
	"bytes"
	std "encoding/binary"
	"reflect"
	"time"

	"github.com/vipally/binary"
)

type BenchCase struct {
	Id               int
	Name             string
	DoCnt            int
	EnableSerializer bool
	Data             interface{}
}

var cases []*BenchCase

type benchType byte

const (
	BenchStdWrite benchType = iota
	BenchStdRead
	BenchEncode
	BenchDecode

	benchDoCnt = 100000
)

func (bench benchType) String() string {
	switch bench {
	case BenchStdWrite:
		return "BenchStdWrite"
	case BenchStdRead:
		return "BenchStdRead"
	case BenchEncode:
		return "BenchEncode"
	case BenchDecode:
		return "BenchDecode"
	}
	panic("undefined benchType")
}

var (
	buff   = make([]byte, 8192)
	buffer = bytes.NewBuffer(buff[:0])
)

func BenchCases() []*BenchCase {
	return cases
}

// DoBench runs a bench test case for binary
func DoBench(bench benchType, data interface{},
	doCnt int, enableSerializer bool, name string) (t time.Duration, speed float64) {
	start := time.Now()
	var err error
	switch bench {
	case BenchStdWrite:
		s := std.Size(data)
		if s <= 0 {
			println(name, "unsupported ")
			return 0, 0
		}
		for i := 0; i < doCnt; i++ {
			buffer.Reset()
			std.Write(buffer, std.LittleEndian, data)
		}
	case BenchStdRead:
		s := std.Size(data)
		if s <= 0 {
			println(name, "unsupported ")
			return 0, 0
		}
		if err = std.Write(buffer, std.LittleEndian, data); err != nil {
			panic(err)
		}
		w := newSame(data)
		b := buffer.Bytes()
		for i := 0; i < doCnt; i++ {
			r := binary.BytesReader(b)
			std.Read(&r, std.LittleEndian, w)
		}
	case BenchEncode:
		for i := 0; i < doCnt; i++ {
			_, err = binary.EncodeX(data, buff, enableSerializer)
		}
		if err != nil {
			panic(err)
		}
	case BenchDecode:
		std.Write(buffer, std.LittleEndian, data)
		w := newSame(data)
		b := buffer.Bytes()
		for i := 0; i < doCnt; i++ {
			err = binary.DecodeX(b, w, enableSerializer)
		}
		if err != nil {
			panic(err)
		}
	}

	dur := time.Now().Sub(start)
	return dur, 0
}

func newSame(x interface{}) interface{} {
	t := reflect.TypeOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	_new := reflect.New(t)
	switch t.Kind() {
	case reflect.Slice:
		_new.Set(reflect.MakeSlice(t, t.Len(), t.Len()))
	case reflect.Map:
		_new.Set(reflect.MakeMap(t))
	}
	return _new.Interface()
}

type BaseStruct struct {
	Bool       bool
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Float32    float32
	Float64    float64
	Complex64  complex64
	Complex128 complex128

	Array     [4]uint8
	BoolArray [9]bool
}

var BaseStruct_ = BaseStruct{
	Bool:       false,
	Int8:       0x12,
	Int16:      0x1234,
	Int32:      0x12345678,
	Int64:      0x123456789abcdef0,
	Uint8:      0x12,
	Uint16:     0x1234,
	Uint32:     0x71234568,
	Uint64:     0xa123456789bcdef0,
	Float32:    1234.5678,
	Float64:    2345.6789012,
	Complex64:  complex(1.12456453, 2.344565),
	Complex128: complex(333.4569789789123, 567.34577890012),

	Array:     [4]uint8{0x1, 0x2, 0x3, 0x4},
	BoolArray: [9]bool{true, false, false, false, false, true, true, false, true},
}

type FastValues struct {
	Int             int
	Uint            uint
	Bool            bool
	Int8            int8
	Int16           int16
	Int32           int32
	Int64           int64
	Uint8           uint8
	Uint16          uint16
	Uint32          uint32
	Uint64          uint64
	Float32         float32
	Float64         float64
	Complex64       complex64
	Complex128      complex128
	String          string
	IntSlice        []int
	UintSlice       []uint
	BoolSlice       []bool
	Int8Slice       []int8
	Int16Slice      []int16
	Int32Slice      []int32
	Int64Slice      []int64
	Uint8Slice      []uint8
	Uint16Slice     []uint16
	Uint32Slice     []uint32
	Uint64Slice     []uint64
	Float32Slice    []float32
	Float64Slice    []float64
	Complex64Slice  []complex64
	Complex128Slice []complex128
	StringSlice     []string
}

var FastValues_ = FastValues{
	Int:             -2,
	Uint:            2,
	Bool:            true,
	Int8:            -3,
	Int16:           -4,
	Int32:           -5,
	Int64:           -6,
	Uint8:           3,
	Uint16:          4,
	Uint32:          5,
	Uint64:          6,
	Float32:         -7,
	Float64:         7,
	Complex64:       8,
	Complex128:      9,
	String:          "hello",
	IntSlice:        []int{-1, 2},
	UintSlice:       []uint{1, 3},
	BoolSlice:       []bool{false, true},
	Int8Slice:       []int8{-1, 2},
	Int16Slice:      []int16{-1, 3},
	Int32Slice:      []int32{-1, 4},
	Int64Slice:      []int64{-1, 5},
	Uint8Slice:      []uint8{1, 6},
	Uint16Slice:     []uint16{1, 7},
	Uint32Slice:     []uint32{1, 8},
	Uint64Slice:     []uint64{1, 9},
	Float32Slice:    []float32{1, 10.1},
	Float64Slice:    []float64{1, 11.2},
	Complex64Slice:  []complex64{1, 2.2},
	Complex128Slice: []complex128{1, 12.9},
	StringSlice:     []string{"abc", "bcd"},
}

type FullStruct struct {
}

var FullStruct_ = FullStruct{}

type RegedStruct struct {
}

var RegedStruct_ = RegedStruct{}

type Serializer struct {
}

func (s Serializer) Size() int {
	return 0
}

func (s Serializer) Encode(buffer []byte) ([]byte, error) {
	return nil, nil
}

func (s *Serializer) Decode(buffer []byte) error {
	return nil
}

func init() {
	binary.RegisterType((*RegedStruct)(nil))
	binary.RegisterType((*Serializer)(nil))

	cases = []*BenchCase{
		&BenchCase{0, "FastValues", benchDoCnt, false, FastValues_},
		&BenchCase{0, "BaseValues", benchDoCnt, false, BaseStruct_},
		&BenchCase{0, "FullValues", benchDoCnt, false, FullStruct_},
		&BenchCase{0, "RegedStruct", benchDoCnt, false, RegedStruct_},
	}
}