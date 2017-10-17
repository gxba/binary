package binary

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"testing"
)

type baseStruct struct {
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
	Array      [4]uint8
	Bool       bool
	BoolArray  [9]bool
}

type littleStruct struct {
	String string
	Int16  int16
}

type fullStruct struct {
	BaseStruct    baseStruct
	LittleStruct  littleStruct
	PLittleStruct *littleStruct
	String        string
	PString       *string
	PInt32        *int32
	Slice         []*littleStruct
	PSlice        *[]*string
	Float64Slice  []float64
	BoolSlice     []bool
	Uint32Slice   []uint32
	Map           map[string]*littleStruct
	Map2          map[string]uint16
	IntSlice      []int
	UintSlice     []uint
}

var full = fullStruct{
	BaseStruct: baseStruct{
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
		Array:      [4]uint8{0x1, 0x2, 0x3, 0x4},
		Bool:       false,
		BoolArray:  [9]bool{true, false, false, false, false, true, true, false, true},
	},
	LittleStruct: littleStruct{
		String: "abc",
		Int16:  0x1234,
	},
	PLittleStruct: &littleStruct{
		String: "bcd",
		Int16:  0x2345,
	},
	String:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	PString: newString("hello"),
	PInt32:  newInt32(0x11223344),
	Slice: []*littleStruct{
		&littleStruct{
			String: "abc",
			Int16:  0x1122,
		},
		&littleStruct{
			String: "bcd",
			Int16:  0x2233,
		},
		&littleStruct{
			String: "cdef",
			Int16:  0x3344,
		},
	},
	PSlice:       &[]*string{newString("abc"), newString("def"), newString("ghijkl")},
	Float64Slice: []float64{3.141592654, 1.137856998, 6.789012345},
	BoolSlice:    []bool{false, true, false, false, true, true, false},
	Uint32Slice:  []uint32{0x12345678, 0x23456789, 0x34567890, 0x4567890a, 0x567890ab},
	Map:          map[string]*littleStruct{"a": &littleStruct{String: "a", Int16: 0x1122}, "b": &littleStruct{String: "b", Int16: 0x1122}},
	Map2:         map[string]uint16{"aaa": 0x5566, "bbb": 0x7788},
	IntSlice:     []int{0, -1, 1, -2, 2, -63, 63, -64, 64, -65, 65, -125, 125, -126, 126, -127, 127, -128, 128, -32765, 32765, -32766, 32766, -32767, 32767, -32768, 32768, -2147483645, 2147483645, -2147483646, 2147483646, -2147483647, 2147483647, -2147483648, 2147483648, -9223372036854775807, 9223372036854775806, -9223372036854775808, 9223372036854775807},
	UintSlice:    []uint{0, 1, 2, 127, 128, 32765, 32766, 32767, 32768, 65533, 65534, 65535, 65536, 0xFFFFFD, 0xFFFFFE, 0xFFFFFF, 0xFFFFFFFFFFFFFFFD, 0xFFFFFFFFFFFFFFFE, 0xFFFFFFFFFFFFFFFF},
}

func newString(s string) *string {
	p := new(string)
	*p = s
	return p
}
func newInt32(i int32) *int32 {
	p := new(int32)
	*p = i
	return p
}

var bigFull = []byte{
	0x12,
	0x12, 0x34,
	0x12, 0x34, 0x56, 0x78,
	0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
	0x12,
	0x12, 0x34,
	0x71, 0x23, 0x45, 0x68,
	0xa1, 0x23, 0x45, 0x67, 0x89, 0xbc, 0xde, 0xf0,
	0x44, 0x9a, 0x52, 0x2b,
	0x40, 0xa2, 0x53, 0x5b, 0x98, 0xf0, 0x26, 0x6e, //Float64
	0x3f, 0x8f, 0xf1, 0xbb, 0x40, 0x16, 0xd, 0x5a,
	0x40, 0x74, 0xd7, 0x4f, 0xc9, 0x30, 0x96, 0x34, 0x40, 0x81, 0xba, 0xc4, 0x27, 0xba, 0x5d, 0x4c,
	0x4, 0x1, 0x2, 0x3, 0x4,
	0x0,
	0x9, 0x61, 0x1, //BoolArray, end of BaseStruct
	0x3, 0x61, 0x62, 0x63,
	0x12, 0x34, //end of LittleStruct
	0x3, 0x62, 0x63, 0x64,
	0x23, 0x45, //end of PLittleStruct
	0x40, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66,
	0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f,
	0x11, 0x22, 0x33, 0x44, //PInt32
	0x3,
	0x3, 0x61, 0x62, 0x63,
	0x11, 0x22,
	0x3, 0x62, 0x63, 0x64,
	0x22, 0x33,
	0x4, 0x63, 0x64, 0x65, 0x66,
	0x33, 0x44, //end of Slice
	0x3,
	0x3, 0x61, 0x62, 0x63,
	0x3, 0x64, 0x65, 0x66,
	0x6, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c,
	0x3,
	0x40, 0x9, 0x21, 0xfb, 0x54, 0x52, 0x45, 0x50,
	0x3f, 0xf2, 0x34, 0xa9, 0x8a, 0x1e, 0xf4, 0xaf,
	0x40, 0x1b, 0x27, 0xf2, 0xda, 0x27, 0xa9, 0x3c, //end of Float64Slice
	0x7, 0x32, //end of BoolSlice
	0x5,
	0x12, 0x34, 0x56, 0x78,
	0x23, 0x45, 0x67, 0x89,
	0x34, 0x56, 0x78, 0x90,
	0x45, 0x67, 0x89, 0xa,
	0x56, 0x78, 0x90, 0xab, //end of Uint32Slice
	0x2, 0x1, 0x61, 0x1, 0x61, 0x11, 0x22, 0x1, 0x62, 0x1, 0x62, 0x11, 0x22,
	0x2, 0x3, 0x61, 0x61, 0x61, 0x55, 0x66, 0x3, 0x62, 0x62, 0x62, 0x77, 0x88,
	//IntSlice
	0x27, 0x0, 0x1, 0x2, 0x3, 0x4, 0x7d, 0x7e, 0x7f, 0x80, 0x1, 0x81, 0x1, 0x82, 0x1, 0xf9, 0x1, 0xfa, 0x1, 0xfb, 0x1, 0xfc, 0x1, 0xfd, 0x1, 0xfe, 0x1, 0xff, 0x1, 0x80, 0x2, 0xf9, 0xff, 0x3, 0xfa, 0xff, 0x3, 0xfb, 0xff, 0x3, 0xfc, 0xff, 0x3, 0xfd, 0xff, 0x3, 0xfe, 0xff, 0x3, 0xff, 0xff, 0x3, 0x80, 0x80, 0x4, 0xf9, 0xff, 0xff, 0xff, 0xf, 0xfa, 0xff, 0xff, 0xff, 0xf, 0xfb, 0xff, 0xff, 0xff, 0xf, 0xfc, 0xff, 0xff, 0xff, 0xf, 0xfd, 0xff, 0xff, 0xff, 0xf, 0xfe, 0xff, 0xff, 0xff, 0xf, 0xff, 0xff, 0xff, 0xff, 0xf, 0x80, 0x80, 0x80, 0x80, 0x10, 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1, 0xfc, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1,
	//UintSlice
	0x13, 0x0, 0x1, 0x2, 0x7f, 0x80, 0x1, 0xfd, 0xff, 0x1, 0xfe, 0xff, 0x1, 0xff, 0xff, 0x1, 0x80, 0x80, 0x2, 0xfd, 0xff, 0x3, 0xfe, 0xff, 0x3, 0xff, 0xff, 0x3, 0x80, 0x80, 0x4, 0xfd, 0xff, 0xff, 0x7, 0xfe, 0xff, 0xff, 0x7, 0xff, 0xff, 0xff, 0x7, 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1,
}

var littleFull = []byte{
	0x12,
	0x34, 0x12,
	0x78, 0x56, 0x34, 0x12,
	0xf0, 0xde, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0x12,
	0x12,
	0x34, 0x12,
	0x68, 0x45, 0x23, 0x71,
	0xf0, 0xde, 0xbc, 0x89, 0x67, 0x45, 0x23, 0xa1,
	0x2b, 0x52, 0x9a, 0x44,
	0x6e, 0x26, 0xf0, 0x98, 0x5b, 0x53, 0xa2, 0x40, //Float64
	0xbb, 0xf1, 0x8f, 0x3f, 0x5a, 0xd, 0x16, 0x40,
	0x34, 0x96, 0x30, 0xc9, 0x4f, 0xd7, 0x74, 0x40, 0x4c, 0x5d, 0xba, 0x27, 0xc4, 0xba, 0x81, 0x40,
	0x4, 0x1, 0x2, 0x3, 0x4,
	0x0,
	0x9, 0x61, 0x1, //BoolArray, end of BaseStruct
	0x3, 0x61, 0x62, 0x63,
	0x34, 0x12, //end of LittleStruct
	0x3, 0x62, 0x63, 0x64,
	0x45, 0x23, //end of PLittleStruct
	0x40, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66,
	0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f,
	0x44, 0x33, 0x22, 0x11, //PInt32
	0x3,
	0x3, 0x61, 0x62, 0x63,
	0x22, 0x11, //end of Slice[0]
	0x3, 0x62, 0x63, 0x64,
	0x33, 0x22,
	0x4, 0x63, 0x64, 0x65, 0x66,
	0x44, 0x33, //end of Slice[2]
	0x3,
	0x3, 0x61, 0x62, 0x63,
	0x3, 0x64, 0x65, 0x66,
	0x6, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, //end of PSlice
	0x3,
	0x50, 0x45, 0x52, 0x54, 0xfb, 0x21, 0x9, 0x40,
	0xaf, 0xf4, 0x1e, 0x8a, 0xa9, 0x34, 0xf2, 0x3f,
	0x3c, 0xa9, 0x27, 0xda, 0xf2, 0x27, 0x1b, 0x40, //end of Float64Slice
	0x7, 0x32, //end of BoolSlice
	0x5,
	0x78, 0x56, 0x34, 0x12,
	0x89, 0x67, 0x45, 0x23,
	0x90, 0x78, 0x56, 0x34,
	0xa, 0x89, 0x67, 0x45,
	0xab, 0x90, 0x78, 0x56, //end of Uint32Slice
	0x2, 0x1, 0x61, 0x1, 0x61, 0x22, 0x11, 0x1, 0x62, 0x1, 0x62, 0x22, 0x11,
	0x2, 0x3, 0x61, 0x61, 0x61, 0x66, 0x55, 0x3, 0x62, 0x62, 0x62, 0x88, 0x77,
	//IntSlice
	0x27, 0x0, 0x1, 0x2, 0x3, 0x4, 0x7d, 0x7e, 0x7f, 0x80, 0x1, 0x81, 0x1, 0x82, 0x1, 0xf9, 0x1, 0xfa, 0x1, 0xfb, 0x1, 0xfc, 0x1, 0xfd, 0x1, 0xfe, 0x1, 0xff, 0x1, 0x80, 0x2, 0xf9, 0xff, 0x3, 0xfa, 0xff, 0x3, 0xfb, 0xff, 0x3, 0xfc, 0xff, 0x3, 0xfd, 0xff, 0x3, 0xfe, 0xff, 0x3, 0xff, 0xff, 0x3, 0x80, 0x80, 0x4, 0xf9, 0xff, 0xff, 0xff, 0xf, 0xfa, 0xff, 0xff, 0xff, 0xf, 0xfb, 0xff, 0xff, 0xff, 0xf, 0xfc, 0xff, 0xff, 0xff, 0xf, 0xfd, 0xff, 0xff, 0xff, 0xf, 0xfe, 0xff, 0xff, 0xff, 0xf, 0xff, 0xff, 0xff, 0xff, 0xf, 0x80, 0x80, 0x80, 0x80, 0x10, 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1, 0xfc, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1,
	//UintSlice
	0x13, 0x0, 0x1, 0x2, 0x7f, 0x80, 0x1, 0xfd, 0xff, 0x1, 0xfe, 0xff, 0x1, 0xff, 0xff, 0x1, 0x80, 0x80, 0x2, 0xfd, 0xff, 0x3, 0xfe, 0xff, 0x3, 0xff, 0xff, 0x3, 0x80, 0x80, 0x4, 0xfd, 0xff, 0xff, 0x7, 0xfe, 0xff, 0xff, 0x7, 0xff, 0xff, 0xff, 0x7, 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1,
}

func TestGenerate(t *testing.T) {
	//	var a, b int = -32767, 0
	//	var c, d int = 32767, 0
	//	var e, f uint = 0xfffffe, 0
	//	b1, _ := Pack(a, nil)
	//	Unpack(b1, &b)
	//	b2, _ := Pack(c, nil)
	//	Unpack(b2, &d)
	//	b3, _ := Pack(e, nil)
	//	Unpack(b3, &f)
	//	fmt.Printf("1 %#v %X \n%#v\n%#v\n", a, a, b, b1)
	//	fmt.Printf("2 %#v %X \n%#v\n%#v\n", c, c, d, b2)
	//	fmt.Printf("2 %#v %X \n%#v\n%#v\n", e, e, f, b3)
	//	b1, _ := Pack(full.IntSlice, nil)
	//	b2, _ := Pack(full.UintSlice, nil)
	//	fmt.Printf("%#v\n%#v\n", full.IntSlice, b1)
	//	fmt.Printf("%#v\n%#v\n", full.UintSlice, b2)
	//	for i, v := range full.IntSlice {
	//		b, _ := Pack(v, nil)
	//		fmt.Printf("int %d %x %x %d %#v\n", i, v, ToUvarint(int64(v)), len(b), b)
	//	}
	//	for i, v := range full.UintSlice {
	//		b, _ := Pack(v, nil)
	//		fmt.Printf("int %d %x %d %#v\n", i, v, len(b), b)
	//	}
}

func TestPack(t *testing.T) {
	v := reflect.ValueOf(full)
	vt := v.Type()
	n := v.NumField()
	check := littleFull
	for i := 0; i < n; i++ {
		if !validField(vt.Field(i)) {
			continue
		}
		b, err := Pack(v.Field(i).Interface(), nil)
		c := check[:len(b)]
		check = check[len(b):]
		if err != nil {
			t.Error(err)
		}
		if vt.Field(i).Type.Kind() != reflect.Map && //map keys will be got as unspecified order, byte order may change but it doesn't matter
			!reflect.DeepEqual(b, c) {
			t.Errorf("field %d %s got %+v\nneed %+v\n", i, vt.Field(i).Name, b, c)
		}
	}

	//	b, err := Pack(full, nil)
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	if !reflect.DeepEqual(b, littleFull) {
	//		t.Errorf("got %+v\nneed %+v\n", b, littleFull)
	//	}
}

func TestPackEmptyPointer(t *testing.T) {
	var s struct {
		PString *string
		PSlice  *[]int
		PArray  *[2]bool
		PInt    *int32
		PStruct *struct{ A int }
		//PStruct2 *struct{ B uintptr }
	}
	b, err := Pack(&s, nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v\n%#v\n", s, b)
}

func TestUnpack(t *testing.T) {
	var v fullStruct
	err := Unpack(littleFull, &v)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(v, full) {
		t.Errorf("got %#v\nneed %#v\n", v, full)
	}
}

func BenchmarkEncoder(b *testing.B) {

}

func _TestGob(t *testing.T) {
	var s struct {
		P *int32
		Q int16
		R string
		T struct {
			A int8
			B string
		}
	}
	s.P = new(int32)
	s.Q = 5
	s.R = "hello"
	s.T.A = 3
	s.T.B = "abc"

	buff := make([]byte, 0, 32)
	bu := bytes.NewBuffer(buff)
	coder := gob.NewEncoder(bu)
	err := coder.Encode(s)
	fmt.Println(err)
	fmt.Printf("before pack: %#v\n", s)
	bb := bu.Bytes()
	fmt.Printf("gobencode: %d %#v\n", len(bb), bb)
	b, _ := Pack(s, nil)
	fmt.Printf("packencode: %d %#v\n", len(b), b)
}

func _TestGob2(t *testing.T) {
	type SubStruct struct {
		Strings []string
		Int16s  []int16
		Uint32s []uint32
	}
	type S struct {
		Bools      []bool
		Bools2     [9]bool
		Bool       bool
		Int8       int8
		Uint8      uint8
		Int16      int16
		Uint16     uint16
		Int32      int32
		Uint32     uint32
		Int64      int64
		Uint64     uint64
		Float32    float32
		Float64    float64
		Complex64  complex64
		Complex128 complex128
		String     string
		//		PString    *string
		Struct SubStruct
	}
	var s S = S{
		Bools:  []bool{true, true, false, false, true, true, false, false, false, true},
		Bools2: [9]bool{false, false, false, true, true, true, false, false},
		Int8:   5,
		Uint32: 277,
		String: "sss",
		//		PString: new(string),
		Struct: SubStruct{
			Strings: []string{"aaa", "bbb", "ccc"},
		},
	}
	buffer := bytes.NewBuffer(make([]byte, 0, 512))
	coder := gob.NewEncoder(buffer)
	err := coder.Encode(s)
	fmt.Println(err)
	fmt.Printf("before pack: %#v\n", s)
	bb := buffer.Bytes()
	fmt.Printf("gobencode: %d %#v\n", len(bb), bb)
	//fmt.Println(Size(s))
	b, _ := Pack(s, nil)
	fmt.Printf("packencode: %d %#v\n", len(b), b)
	var tt S
	err = Unpack(b, &tt)
	fmt.Println(err)
	fmt.Printf("after unpack: %#v\n", tt)
	fmt.Println(reflect.DeepEqual(s, tt))

}

func _TestSlice(t *testing.T) {
	var s = []byte("hello")
	testSlice(&s)
	fmt.Printf("%#v\n", s)
}

func testSlice(i interface{}) {
	//	v := reflect.ValueOf(i)
	//	//t := []byte("asdf")
	//	v.Elem().Set(reflect.MakeSlice(v.Elem().Type(), 2, 2))
	//	//v.Elem().Set(reflect.ValueOf(&t).Elem())

	//	var s []byte
	//	var ss []byte = []byte("hh")
	//	p := &s
	//	var q *[]byte
	//	r := &ss

	//	fmt.Printf("%#v\n", p)
	//	fmt.Printf("%#v\n", q)
	//	fmt.Printf("%#v\n", r)

	var s struct {
		P *int
	}

	//	var pi *int
	//	fmt.Printf("%#v\n", pi)
	//	v := reflect.ValueOf(&pi)
	//	vv := reflect.Indirect(v)
	//	vv.Set(reflect.New(vv.Type().Elem()))
	//	fmt.Printf("%#v\n", pi)
	//	pi = new(int)
	//	fmt.Printf("%#v\n", pi)

	v := reflect.ValueOf(&s)
	vv := reflect.Indirect(v)
	f := vv.Field(0)
	f.Set(reflect.New(f.Type().Elem()))
	fmt.Printf("%#v\n", s)
}

func __TestPack(t *testing.T) {
	var s struct {
		P *int32
		Q int16
		R string
		T struct {
			A int8
			B string
		}
	}
	s.P = new(int32)
	s.Q = 5
	s.R = "hello"
	s.T.A = 3
	s.T.B = "abc"

	b, err := Pack(s, nil)
	//fmt.Println("err:", err)
	fmt.Printf("before pack: %#v\n", s)
	fmt.Printf("pack data: len:%d, %#v\n", len(b), b)

	tt := s
	err2 := Unpack(b, &tt)
	err2 = err2
	err = err
	//fmt.Println("err:", err2)

	fmt.Printf("after unpack: %#v\n", tt)
}

func TestStruct(t *testing.T) {
	type T struct {
		A uint32
		b uint32
		_ uint32
		C uint32 `binary:"ignore"`
	}
	var s T
	s.A = 0x11223344
	s.b = 0x22334455
	s.C = 0x33445566
	check := []byte{0x44, 0x33, 0x22, 0x11}
	b, err := Pack(s, nil)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(b, check) {
		t.Errorf("%T: got %x; want %x", s, b, check)
	}
	var ss, ssCheck T
	ssCheck.A = s.A
	err = Unpack(b, &ss)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(ss, ssCheck) {
		t.Errorf("%T: got %q; want %q", s, ss, ssCheck)
	}
}
