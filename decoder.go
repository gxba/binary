package binary

import (
	"fmt"
	"io"
	"math"
	"reflect"
)

// NewDecoder make a new Decoder object with buffer.
func NewDecoder(buffer []byte) *Decoder {
	return NewDecoderEndian(buffer, DefaultEndian)
}

// NewDecoderEndian make a new Decoder object with buffer and endian.
func NewDecoderEndian(buffer []byte, endian Endian) *Decoder {
	p := &Decoder{}
	p.Init(buffer, endian)
	return p
}

// Decoder is used to decode byte array to go data.
type Decoder struct {
	coder
	reader io.Reader //for decode from reader only
}

// Skip ignore the next size of bytes for encoding/decoding.
// It will panic If space not enough.
// It will return -1 if size <= 0.
func (this *Decoder) Skip(size int) int {
	if nil == this.reserve(size) {
		return -1
	}
	return size
}

// reserve returns next size bytes for encoding/decoding.
func (this *Decoder) reserve(size int) []byte {
	if this.reader != nil { //decode from reader
		if size > len(this.buff) {
			this.buff = make([]byte, size)
		}
		buff := this.buff[:size]
		if n, _ := this.reader.Read(buff); n < size {
			panic(io.ErrUnexpectedEOF)
		}
		return buff
	} else { //decode from bytes buffer
		return this.coder.reserve(size)
	}
}

// Init initialize Encoder with buffer and endian.
func (this *Decoder) Init(buffer []byte, endian Endian) {
	this.buff = buffer
	this.pos = 0
	this.endian = endian
}

// Bool decode a bool value from Decoder buffer.
// It will panic if buffer is not enough.
func (this *Decoder) Bool() bool {
	b := this.reserve(1)
	x := b[0]
	return x != 0
}

// Int8 decode an int8 value from Decoder buffer.
// It will panic if buffer is not enough.
func (this *Decoder) Int8() int8 {
	return int8(this.Uint8())
}

// Uint8 decode a uint8 value from Decoder buffer.
// It will panic if buffer is not enough.
func (this *Decoder) Uint8() uint8 {
	b := this.reserve(1)
	x := b[0]
	return x
}

// Int16 decode an int16 value from Decoder buffer.
// It will panic if buffer is not enough.
func (this *Decoder) Int16() int16 {
	return int16(this.Uint16())
}

// Uint16 decode a uint16 value from Decoder buffer.
// It will panic if buffer is not enough.
func (this *Decoder) Uint16() uint16 {
	b := this.reserve(2)
	x := this.endian.Uint16(b)
	return x
}

// Int32 decode an int32 value from Decoder buffer.
// It will panic if buffer is not enough.
func (this *Decoder) Int32() int32 {
	return int32(this.Uint32())
}

// Uint32 decode a uint32 value from Decoder buffer.
// It will panic if buffer is not enough.
func (this *Decoder) Uint32() uint32 {
	b := this.reserve(4)
	x := this.endian.Uint32(b)
	return x
}

// Int64 decode an int64 value from Decoder buffer.
// It will panic if buffer is not enough.
func (this *Decoder) Int64() int64 {
	return int64(this.Uint64())
}

// Uint64 decode a uint64 value from Decoder buffer.
// It will panic if buffer is not enough.
func (this *Decoder) Uint64() uint64 {
	b := this.reserve(8)
	x := this.endian.Uint64(b)
	return x
}

// Float32 decode a float32 value from Decoder buffer.
// It will panic if buffer is not enough.
func (this *Decoder) Float32() float32 {
	x := math.Float32frombits(this.Uint32())
	return x
}

// Float64 decode a float64 value from Decoder buffer.
// It will panic if buffer is not enough.
func (this *Decoder) Float64() float64 {
	x := math.Float64frombits(this.Uint64())
	return x
}

// Complex64 decode a complex64 value from Decoder buffer.
// It will panic if buffer is not enough.
func (this *Decoder) Complex64() complex64 {
	x := complex(this.Float32(), this.Float32())
	return x
}

// Complex128 decode a complex128 value from Decoder buffer.
// It will panic if buffer is not enough.
func (this *Decoder) Complex128() complex128 {
	x := complex(this.Float64(), this.Float64())
	return x
}

// String decode a string value from Decoder buffer.
// It will panic if buffer is not enough.
func (this *Decoder) String() string {
	s, _ := this.Uvarint()
	size := int(s)
	b := this.reserve(size)
	return string(b)
}

// Int decode an int value from Decoder buffer.
// It will panic if buffer is not enough.
// It use Varint() to decode as varint(1~10 bytes)
func (this *Decoder) Int() int {
	n, _ := this.Varint()
	return int(n)
}

// Uint decode a uint value from Decoder buffer.
// It will panic if buffer is not enough.
// It use Uvarint() to decode as uvarint(1~10 bytes)
func (this *Decoder) Uint() uint {
	n, _ := this.Uvarint()
	return uint(n)
}

// Varint decode an int64 value from Decoder buffer with varint(1~10 bytes).
// It will panic if buffer is not enough.
func (this *Decoder) Varint() (int64, int) {
	ux, n := this.Uvarint() // ok to continue in presence of error
	return ToVarint(ux), n
}

// Uvarint decode a uint64 value from Decoder buffer with varint(1~10 bytes).
// It will panic if buffer is not enough.
// It will return n <= 0 if varint error
func (this *Decoder) Uvarint() (uint64, int) {
	var x uint64 = 0
	var bit uint = 0
	i := 0
	for i = 0; i < MaxVarintLen64; i++ {
		b := this.Uint8()
		x |= uint64(b&0x7f) << bit
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				break // overflow
			}
			return x, i + 1
		}
		bit += 7
	}
	//return 0, 0
	panic(fmt.Errorf("binary.Decoder.Uvarint: overflow 64-bits value(pos:%d/%d).", this.Len(), this.Cap()))
}

// Value decode an interface value from Encoder buffer.
// x must be interface of pointer for modify.
// It will return none-nil error if x contains unsupported types
// or buffer is not enough.
func (this *Decoder) Value(x interface{}) (err error) {
	defer func() {
		if info := recover(); info != nil {
			err = info.(error)
			assert(err != nil, info)
		}
	}()
	if this.fastValue(x) { //fast value path
		return nil
	}

	v := reflect.ValueOf(x)

	if p, ok := x.(Unpacker); ok {
		size := 0
		if sizer, _ok := x.(Sizer); _ok { //interface verification
			size = sizer.Size()
		} else {
			panic(fmt.Errorf("expect but not Sizer: %s", v.Type().String()))
		}
		if _, _ok := x.(Packer); !_ok { //interface verification
			panic(fmt.Errorf("unexpect but not Packer: %s", v.Type().String()))
		}
		err := p.Unpack(this.buff[this.pos:])
		if err != nil {
			return err
		}
		this.reserve(size)
		return nil
	} else {
		if _, _ok := x.(Sizer); _ok { //interface verification
			panic(fmt.Errorf("unexpected Sizer: %s", v.Type().String()))
		}
		if _, _ok := x.(Packer); _ok { //interface verification
			panic(fmt.Errorf("unexpected Packer: %s", v.Type().String()))
		}
	}

	if v.Kind() == reflect.Ptr { //only support decode for pointer interface
		return this.value(v)
	} else {
		return fmt.Errorf("binary.Decoder.Value: non-pointer type %s", v.Type().String())
	}
}

func (this *Decoder) value(v reflect.Value) error {
	// check Packer interface for every value is perfect
	// but this is too costly
	//
	//	if t := v.Type(); t.Implements(tUnpacker) {
	//		if !t.Implements(tPacker) { //interface verification
	//			panic(fmt.Errorf("unexpect but not Packer: %s", v.Type().String()))
	//		}
	//		if !t.Implements(tSizer) { //interface verification
	//			panic(fmt.Errorf("expect but not Sizer: %s", t.String()))
	//		}

	//		unpacker := v.Interface().(PackUnpacker)
	//		size := unpacker.Size()
	//		err := unpacker.Unpack(this.buff[this.pos:])
	//		if err != nil {
	//			return err
	//		}
	//		this.reserve(size)
	//		return nil
	//	} else {
	//		//interface verification
	//		if t.Implements(tSizer) {
	//			panic(fmt.Errorf("unexpected Sizer: %s", t.String()))
	//		}
	//		if t.Implements(tPacker) {
	//			panic(fmt.Errorf("unexpected Packer: %s", t.String()))
	//		}
	//	}

	switch k := v.Kind(); k {
	case reflect.Int:
		v.SetInt(int64(this.Int()))
	case reflect.Uint:
		v.SetUint(uint64(this.Uint()))

	case reflect.Bool:
		v.SetBool(this.Bool())

	case reflect.Int8:
		v.SetInt(int64(this.Int8()))
	case reflect.Int16:
		v.SetInt(int64(this.Int16()))
	case reflect.Int32:
		v.SetInt(int64(this.Int32()))
	case reflect.Int64:
		v.SetInt(this.Int64())

	case reflect.Uint8:
		v.SetUint(uint64(this.Uint8()))
	case reflect.Uint16:
		v.SetUint(uint64(this.Uint16()))
	case reflect.Uint32:
		v.SetUint(uint64(this.Uint32()))
	case reflect.Uint64:
		v.SetUint(this.Uint64())

	case reflect.Float32:
		v.SetFloat(float64(this.Float32()))
	case reflect.Float64:
		v.SetFloat(this.Float64())

	case reflect.Complex64:
		v.SetComplex(complex128(this.Complex64()))

	case reflect.Complex128:
		v.SetComplex(this.Complex128())

	case reflect.String:
		v.SetString(this.String())

	case reflect.Slice, reflect.Array:
		if sizeofNilPointer(v.Type().Elem()) < 0 { //verify array element is valid
			return fmt.Errorf("binary.Decoder.Value: unsupported type %s", v.Type().String())
		}
		if this.boolArray(v) < 0 { //deal with bool array first
			s, _ := this.Uvarint()
			size := int(s)
			if size > 0 && k == reflect.Slice { //make a new slice
				ns := reflect.MakeSlice(v.Type(), size, size)
				v.Set(ns)
			}

			l := v.Len()
			for i := 0; i < size; i++ {
				if i < l {
					this.value(v.Index(i))
				} else {
					this.skipByType(v.Type().Elem())
				}
			}
		}
	case reflect.Map:
		t := v.Type()
		kt := t.Key()
		vt := t.Elem()
		if sizeofNilPointer(kt) < 0 ||
			sizeofNilPointer(vt) < 0 { //verify map key and value type are both valid
			return fmt.Errorf("binary.Decoder.Value: unsupported type %s", v.Type().String())
		}

		if v.IsNil() {
			newmap := reflect.MakeMap(v.Type())
			v.Set(newmap)
		}

		s, _ := this.Uvarint()
		size := int(s)
		for i := 0; i < size; i++ {
			key := reflect.New(kt).Elem()
			value := reflect.New(vt).Elem()
			this.value(key)
			this.value(value)
			v.SetMapIndex(key, value)
		}
	case reflect.Struct:
		return queryStruct(v.Type()).decode(this, v)

	default:
		if newPtr(v) {
			return this.value(v.Elem())
		} else {
			return fmt.Errorf("binary.Decoder.Value: unsupported type %s", v.Type().String())
		}
	}
	return nil
}

func (this *Decoder) fastValue(x interface{}) bool {
	switch d := x.(type) {
	case *int:
		*d = this.Int()
	case *uint:
		*d = this.Uint()

	case *bool:
		*d = this.Bool()
	case *int8:
		*d = this.Int8()
	case *uint8:
		*d = this.Uint8()

	case *int16:
		*d = this.Int16()
	case *uint16:
		*d = this.Uint16()

	case *int32:
		*d = this.Int32()
	case *uint32:
		*d = this.Uint32()
	case *float32:
		*d = this.Float32()

	case *int64:
		*d = this.Int64()
	case *uint64:
		*d = this.Uint64()
	case *float64:
		*d = this.Float64()
	case *complex64:
		*d = this.Complex64()

	case *complex128:
		*d = this.Complex128()

	case *string:
		*d = this.String()

	case *[]bool:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]bool, l)
		var b []byte
		for i := 0; i < l; i++ {
			_, bit := i/8, i%8
			mask := byte(1 << uint(bit))
			if bit == 0 {
				b = this.reserve(1)
			}
			x := ((b[0] & mask) != 0)
			(*d)[i] = x
		}

	case *[]int:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]int, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Int()
		}
	case *[]uint:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]uint, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Uint()
		}

	case *[]int8:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]int8, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Int8()
		}
	case *[]uint8:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]uint8, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Uint8()
		}
	case *[]int16:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]int16, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Int16()
		}
	case *[]uint16:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]uint16, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Uint16()
		}
	case *[]int32:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]int32, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Int32()
		}
	case *[]uint32:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]uint32, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Uint32()
		}
	case *[]int64:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]int64, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Int64()
		}
	case *[]uint64:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]uint64, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Uint64()
		}
	case *[]float32:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]float32, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Float32()
		}
	case *[]float64:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]float64, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Float64()
		}
	case *[]complex64:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]complex64, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Complex64()
		}
	case *[]complex128:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]complex128, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Complex128()
		}
	case *[]string:
		s, _ := this.Uvarint()
		l := int(s)
		*d = make([]string, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.String()
		}
	default:
		return false
	}
	return true
}

func (this *Decoder) skipByType(t reflect.Type) int {
	if s := _fixTypeSize(t); s > 0 {
		this.Skip(s)
		return s
	}
	switch t.Kind() {
	case reflect.Int:
		_, n := this.Varint()
		return n
	case reflect.Uint:
		_, n := this.Uvarint()
		return n
	case reflect.String:
		s, n := this.Uvarint()
		size := int(s) //string length and data
		this.Skip(size)
		return size + n
	case reflect.Slice, reflect.Array:
		s, sLen := this.Uvarint()
		cnt := int(s)
		e := t.Elem()
		if s := _fixTypeSize(e); s > 0 {
			if t.Elem().Kind() == reflect.Bool { //compressed bool array
				totalSize := sizeofBoolArray(cnt)
				size := totalSize - SizeofUvarint(uint64(cnt)) //cnt has been read
				this.Skip(size)
				return totalSize
			} else {
				size := cnt * s
				this.Skip(size)
				return size
			}
		} else {
			sum := sLen //array size
			for i, n := 0, cnt; i < n; i++ {
				s := this.skipByType(e)
				assert(s >= 0, e.String()) //I'm sure here cannot find unsupported type
				sum += s
			}
			return sum
		}
	case reflect.Map:
		s, sLen := this.Uvarint()
		cnt := int(s)
		kt := t.Key()
		vt := t.Elem()
		sum := sLen //array size
		for i, n := 0, cnt; i < n; i++ {
			sum += this.skipByType(kt)
			sum += this.skipByType(vt)
		}
		return sum

	case reflect.Struct:
		return queryStruct(t).decodeSkipByType(this, t)
	}
	return -1
}

// decode bool array
func (this *Decoder) boolArray(v reflect.Value) int {
	if k := v.Kind(); k == reflect.Slice || k == reflect.Array {
		if v.Type().Elem().Kind() == reflect.Bool {
			_l, _ := this.Uvarint()
			l := int(_l)
			if k == reflect.Slice && l > 0 { //make a new slice
				v.Set(reflect.MakeSlice(v.Type(), l, l))
			}
			var b []byte
			for i := 0; i < l; i++ {
				_, bit := i/8, i%8
				mask := byte(1 << uint(bit))
				if bit == 0 {
					b = this.reserve(1)
				}
				x := ((b[0] & mask) != 0)
				v.Index(i).SetBool(x)
			}
			return sizeofBoolArray(l)
		}
	}
	return -1
}
