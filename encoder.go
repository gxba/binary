package binary

import (
	"fmt"
	"math"
	"reflect"
)

// NewEncoder make a new Encoder object with buffer size.
func NewEncoder(size int) *Encoder {
	return NewEncoderEndian(size, DefaultEndian)
}

// NewEncoderBuffer make a new Encoder object with buffer.
func NewEncoderBuffer(buffer []byte) *Encoder {
	p := &Encoder{}
	//assert(buffer != nil, "nil buffer")
	p.buff = buffer
	p.endian = DefaultEndian
	p.pos = 0
	return p
}

// NewEncoderEndian make a new Encoder object with buffer size and endian.
func NewEncoderEndian(size int, endian Endian) *Encoder {
	p := &Encoder{}
	p.Init(size, endian)
	return p
}

// Encoder is used to encode go data to byte array.
type Encoder struct {
	coder
}

// Init initialize Encoder with buffer size and endian.
func (encoder *Encoder) Init(size int, endian Endian) {
	encoder.buff = make([]byte, size)
	encoder.pos = 0
	encoder.endian = endian
}

// ResizeBuffer confirm that len(buffer) >= size and alloc larger buffer if necessary
// It will call Reset to initial encoder state of buffer
func (encoder *Encoder) ResizeBuffer(size int) bool {
	ok := len(encoder.buff) < size
	if ok {
		encoder.buff = make([]byte, size)
	}
	encoder.Reset()
	return ok
}

// Bool encode a bool value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Bool(x bool) error {
	if encoder.boolBit == 0 {
		b, err := encoder.reserve(1)
		if err != nil {
			return err
		}
		b[0] = 0
		encoder.boolPos = encoder.pos - 1
	}

	if mask := byte(1 << encoder.boolBit); x {
		encoder.buff[encoder.boolPos] |= mask
	}
	encoder.boolBit = (encoder.boolBit + 1) % 8
	return nil
}

// Int8 encode an int8 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Int8(x int8) error {
	return encoder.Uint8(uint8(x))
}

// Uint8 encode a uint8 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Uint8(x uint8) error {
	b, err := encoder.reserve(1)
	if err != nil {
		return err
	}
	b[0] = x
	return nil
}

// Int16 encode an int16 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Int16(x int16, packed bool) (err error) {
	if packed {
		_, err = encoder.Varint(int64(x))
	} else {
		err = encoder.Uint16(uint16(x), false)
	}
	return
}

// Uint16 encode a uint16 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Uint16(x uint16, packed bool) (err error) {
	if packed {
		_, err = encoder.Uvarint(uint64(x))
	} else {
		b, e := encoder.reserve(2)
		if e != nil {
			return e
		}
		encoder.endian.PutUint16(b, x)
	}
	return
}

// Int32 encode an int32 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Int32(x int32, packed bool) (err error) {
	if packed {
		_, err = encoder.Varint(int64(x))
	} else {
		err = encoder.Uint32(uint32(x), false)
	}
	return
}

// Uint32 encode a uint32 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Uint32(x uint32, packed bool) (err error) {
	if packed {
		_, err = encoder.Uvarint(uint64(x))
	} else {
		b, e := encoder.reserve(4)
		if e != nil {
			return e
		}
		encoder.endian.PutUint32(b, x)
	}
	return
}

// Int64 encode an int64 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Int64(x int64, packed bool) (err error) {
	if packed {
		_, err = encoder.Varint(x)
	} else {
		err = encoder.Uint64(uint64(x), false)
	}
	return
}

// Uint64 encode a uint64 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Uint64(x uint64, packed bool) (err error) {
	if packed {
		_, err = encoder.Uvarint(x)
	} else {
		b, e := encoder.reserve(8)
		if e != nil {
			return e
		}
		encoder.endian.PutUint64(b, x)
	}
	return
}

// Float32 encode a float32 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Float32(x float32) error {
	return encoder.Uint32(math.Float32bits(x), false)
}

// Float64 encode a float64 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Float64(x float64) error {
	return encoder.Uint64(math.Float64bits(x), false)
}

// Complex64 encode a complex64 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Complex64(x complex64) (err error) {
	err = encoder.Uint32(math.Float32bits(real(x)), false)
	if err == nil {
		err = encoder.Uint32(math.Float32bits(imag(x)), false)
	}
	return
}

// Complex128 encode a complex128 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Complex128(x complex128) (err error) {
	err = encoder.Uint64(math.Float64bits(real(x)), false)
	if err == nil {
		err = encoder.Uint64(math.Float64bits(imag(x)), false)
	}
	return
}

// String encode a string value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) String(x string) (err error) {
	_b := []byte(x)
	size := len(_b)
	_, err = encoder.Uvarint(uint64(size))
	if err != nil {
		return
	}
	buff, e := encoder.reserve(size)
	if e != nil {
		return e
	}
	copy(buff, _b)
	return
}

// Int encode an int value to Encoder buffer.
// It will panic if buffer is not enough.
// It use Varint() to encode as varint(1~10 bytes)
func (encoder *Encoder) Int(x int) (err error) {
	_, err = encoder.Varint(int64(x))
	return
}

// Uint encode a uint value to Encoder buffer.
// It will panic if buffer is not enough.
// It use Uvarint() to encode as uvarint(1~10 bytes)
func (encoder *Encoder) Uint(x uint) (err error) {
	_, err = encoder.Uvarint(uint64(x))
	return
}

// Varint encode an int64 value to Encoder buffer with varint(1~10 bytes).
// It will panic if buffer is not enough.
func (encoder *Encoder) Varint(x int64) (int, error) {
	return encoder.Uvarint(ToUvarint(x))
}

// Uvarint encode a uint64 value to Encoder buffer with varint(1~10 bytes).
// It will panic if buffer is not enough.
func (encoder *Encoder) Uvarint(x uint64) (int, error) {
	size := SizeofUvarint(x)
	b, err := encoder.reserve(size)
	if err != nil {
		return 0, err
	}
	x_ := x
	for i, s := 0, size-1; i < s; i++ {
		b[i] = byte((x_ & 0x7f) | 0x80)
		x_ >>= 7
	}
	b[size-1] = byte(x_)
	return size, nil
}

// Value encode an interface value to Encoder buffer.
// It will return none-nil error if x contains unsupported types
// or buffer is not enough.
// It will check if x implements interface BinaryEncoder and use x.Encode first.
func (encoder *Encoder) Value(x interface{}) error {
	return encoder.ValueX(x, defaultSerializer)
}

// ValueX encode an interface value to Encoder buffer.
// enableSerializer switch if need check BinarySerilizer.
// It will return none-nil error if x contains unsupported types
// or buffer is not enough.
// It will check if x implements interface BinaryEncoder and use x.Encode first.
func (encoder *Encoder) ValueX(x interface{}, enableSerializer bool) (err error) {
	//	defer func() {
	//		if e := recover(); e != nil {
	//			err = e.(error)
	//		}
	//	}()

	encoder.resetBoolCoder()               //reset bool writer
	if ok, e := encoder.fastValue(x); ok { //fast value path
		return e
	}

	v := reflect.ValueOf(x)
	return encoder.value(reflect.Indirect(v), false, toplvSerializer(enableSerializer))
}

func (encoder *Encoder) fastValue(x interface{}) (ok bool, err error) {
	switch d := x.(type) {
	case int:
		err = encoder.Int(d)
	case uint:
		err = encoder.Uint(d)

	case bool:
		err = encoder.Bool(d)
	case int8:
		err = encoder.Int8(d)
	case uint8:
		err = encoder.Uint8(d)
	case int16:
		err = encoder.Int16(d, false)
	case uint16:
		err = encoder.Uint16(d, false)
	case int32:
		err = encoder.Int32(d, false)
	case uint32:
		encoder.Uint32(d, false)
	case float32:
		err = encoder.Float32(d)
	case int64:
		err = encoder.Int64(d, false)
	case uint64:
		err = encoder.Uint64(d, false)
	case float64:
		err = encoder.Float64(d)
	case complex64:
		err = encoder.Complex64(d)
	case complex128:
		err = encoder.Complex128(d)
	case string:
		err = encoder.String(d)
	case []bool:
		l := len(d)
		_, err = encoder.Uvarint(uint64(l))
		var b []byte
		for i := 0; i < l; i++ {
			bit := i % 8
			mask := byte(1 << uint(bit))
			if bit == 0 {
				b, err = encoder.reserve(1)
				if err != nil {
					return true, err
				}
				b[0] = 0
			}
			if x := d[i]; x {
				b[0] |= mask
			}
		}
	case []int8:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Int8(d[i])
		}
	case []uint8:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Uint8(d[i])
		}
	case []int16:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Int16(d[i], false)
		}
	case []uint16:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Uint16(d[i], false)
		}
	case []int32:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Int32(d[i], false)
		}
	case []uint32:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Uint32(d[i], false)
		}
	case []int64:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Int64(d[i], false)
		}
	case []uint64:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Uint64(d[i], false)
		}
	case []float32:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Float32(d[i])
		}
	case []float64:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Float64(d[i])
		}
	case []complex64:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Complex64(d[i])
		}
	case []complex128:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Complex128(d[i])
		}
	case []string:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.String(d[i])
		}
	case []int:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Uvarint(ToUvarint(int64(d[i])))
			//encoder.Int(d[i])
		}
	case []uint:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Uvarint(uint64(d[i]))
			//encoder.Uint(d[i])
		}
	default:
		return false, nil
	}
	return true, err

}

// use BinarySerializer interface to encode this value
func (encoder *Encoder) useSerializer(v reflect.Value) error {
	return encoder.Serializer(v.Interface())
}

// Serializer encode BinarySerializer x.
func (encoder *Encoder) Serializer(x interface{}) error {
	//	t := reflect.TypeOf(x)
	//	if _, _, _, err := deepRegableType(t, true); err != nil {
	//		return err
	//	}
	if p, ok := x.(BinaryEncoder); ok {
		r, err := p.Encode(encoder.buff[encoder.pos:])
		if err != nil {
			return err

		}
		encoder.mustReserve(len(r))
		return nil
	}

	return typeError("binary: expect implements BinarySerializer %s", reflect.TypeOf(x), true)
}

// valueSerializer encode v with serializer check
func (encoder *Encoder) value(v reflect.Value, packed bool, serializer serializerSwitch) error {
	k := v.Kind()
	if serializer.checkOk() ||
		serializer.needCheck() && k != reflect.Ptr && querySerializer(v.Type()) {
		return encoder.useSerializer(v)
	}

	switch k {
	case reflect.Int:
		encoder.Uvarint(ToUvarint(v.Int()))
		//encoder.Int(int(v.Int()))
	case reflect.Uint:
		encoder.Uvarint(v.Uint())
		//encoder.Uint(uint(v.Uint()))
	case reflect.Bool:
		encoder.Bool(v.Bool())
	case reflect.Int8:
		b := encoder.mustReserve(1)
		b[0] = uint8(v.Int())
		//encoder.Int8(int8(v.Int()))
	case reflect.Int16:
		encoder.Int16(int16(v.Int()), packed)
	case reflect.Int32:
		encoder.Int32(int32(v.Int()), packed)
	case reflect.Int64:
		encoder.Int64(v.Int(), packed)
	case reflect.Uint8:
		encoder.Uint8(uint8(v.Uint()))
	case reflect.Uint16:
		encoder.Uint16(uint16(v.Uint()), packed)
	case reflect.Uint32:
		if packed {
			encoder.Uvarint(v.Uint())
		} else {
			b := encoder.mustReserve(4)
			encoder.endian.PutUint32(b, uint32(v.Uint()))
		}
		//encoder.Uint32(uint32(v.Uint()), packed)
	case reflect.Uint64:
		x := v.Uint()
		if packed {
			encoder.Uvarint(x)
		} else {
			b := encoder.mustReserve(8)
			encoder.endian.PutUint64(b, x)
		}
		//encoder.Uint64(v.Uint(), packed)
	case reflect.Float32:
		encoder.Float32(float32(v.Float()))
	case reflect.Float64:
		encoder.Float64(v.Float())
	case reflect.Complex64:
		x := v.Complex()
		encoder.Complex64(complex64(x))
	case reflect.Complex128:
		x := v.Complex()
		encoder.Complex128(x)
	case reflect.String:
		encoder.String(v.String())

	case reflect.Slice, reflect.Array:
		elemT := v.Type().Elem()
		if !validUserType(elemT) { //verify array element is valid
			return fmt.Errorf("binary.Encoder.Value: unsupported type %s", v.Type().String())
		}
		elemSerializer := serializer.subSwitchCheck(elemT)
		if encoder.boolArray(v) < 0 { //deal with bool array first
			l := v.Len()
			encoder.Uvarint(uint64(l))
			for i := 0; i < l; i++ {
				encoder.value(v.Index(i), packed, elemSerializer) //ignore error
			}
		}

	case reflect.Map:
		t := v.Type()
		kt, vt := t.Key(), t.Elem()
		if !validUserType(kt) || !validUserType(vt) { //verify map key and value type are both valid
			return fmt.Errorf("binary.Decoder.Value: unsupported type %s", v.Type().String())
		}

		keySerilaizer := serializer.subSwitchCheck(kt)
		valueSerilaizer := serializer.subSwitchCheck(vt)
		keys := v.MapKeys()
		l := len(keys)
		encoder.Uvarint(uint64(l))
		for i := 0; i < l; i++ {
			key := keys[i]
			encoder.value(key, packed, keySerilaizer)               //ignore error
			encoder.value(v.MapIndex(key), packed, valueSerilaizer) //ignore error
		}

	case reflect.Struct:
		return queryStruct(v.Type()).encode(encoder, v, serializer)

	case reflect.Ptr:
		if !validUserType(v.Type()) {
			return fmt.Errorf("binary.Encoder.Value: unsupported type %s", v.Type().String())
		}
		if !v.IsNil() {
			encoder.Bool(true)
			if e := v.Elem(); e.Kind() != reflect.Ptr {
				return encoder.value(e, packed, serializer)
			}
		} else {
			encoder.Bool(false) //put a bool to mark nil pointer
		}

	default:
		return typeError("binary.Encoder.Value: unsupported type [%s]", v.Type(), true)
	}
	return nil
}

// encode bool array
func (encoder *Encoder) boolArray(v reflect.Value) int {
	if k := v.Kind(); k == reflect.Slice || k == reflect.Array {
		if v.Type().Elem().Kind() == reflect.Bool {
			l := v.Len()
			encoder.Uvarint(uint64(l))
			var b []byte
			for i := 0; i < l; i++ {
				bit := i % 8
				mask := byte(1 << uint(bit))
				if bit == 0 {
					b = encoder.mustReserve(1)
					b[0] = 0
				}
				if x := v.Index(i).Bool(); x {
					b[0] |= mask
				}
			}
			return sizeofBoolArray(l)
		}
	}
	return -1
}
