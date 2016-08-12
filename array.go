package ndim

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

type ndarray struct {
	buf     unsafe.Pointer
	shape   []int
	strides []int
	etype   reflect.Type
}

func (arr *ndarray) Dims() []int {
	return arr.shape
}

func (arr *ndarray) Strides() []int {
	return arr.strides
}

func (arr *ndarray) Type() reflect.Type {
	return arr.etype
}

type Dtype interface {
	Dims() []int
	Order() binary.ByteOrder
	Type() reflect.Type
}

type Header struct {
	Dims    []int
	Strides []int
	Order   binary.ByteOrder
	Elem    reflect.Type
}

func (hdr Header) dtype() header {
	return header{
		shape:   hdr.Dims,
		strides: hdr.Strides,
		order:   hdr.Order,
		etype:   hdr.Elem,
	}
}

func (hdr Header) Dtype() Dtype {
	return hdr.dtype()
}

type header struct {
	shape   []int
	strides []int
	order   binary.ByteOrder
	etype   reflect.Type
}

func (hdr header) Dims() []int {
	return hdr.shape
}

func (hdr header) Strides() []int {
	return hdr.strides
}

func (hdr header) Order() binary.ByteOrder {
	return hdr.order
}

func (hdr header) Type() reflect.Type {
	return hdr.etype
}

type Reader struct {
	r   io.Reader
	hdr header
}

func NewReader(r io.Reader) (*Reader, error) {
	return nil, nil
}

func (r *Reader) Read() {}

type Decoder struct {
	r   io.Reader
	hdr header
}

func NewDecoder(r io.Reader) (*Decoder, error) {
	dec := &Decoder{
		r: r,
	}
	err := dec.decodeHeader()
	if err != nil {
		return nil, err
	}

	return dec, nil
}

func NewDecoderFrom(r io.Reader, hdr Header) *Decoder {
	return &Decoder{
		r:   r,
		hdr: hdr.dtype(),
	}
}

func (dec *Decoder) Dtype() Dtype {
	return dec.hdr
}

func (dec *Decoder) decodeHeader() error {
	dec.hdr = header{
		shape:   []int{2, 3},
		strides: nil,
		order:   binary.LittleEndian,
		etype:   reflect.TypeOf(float64(0)),
	}
	return nil
}

func (dec *Decoder) Decode(ptr interface{}) error {
	var err error
	if ptr, ok := ptr.(BioReader); ok {
		return ptr.Rio(dec.Dtype(), dec.r)
	}

	switch v := ptr.(type) {
	case *[]float64:
		n := 1
		for _, nn := range dec.Dtype().Dims() {
			n *= nn
		}
		for i := 0; i < n; i++ {
			var vv float64
			err = binary.Read(dec.r, dec.hdr.order, &vv)
			if err != nil {
				return err
			}
			*v = append(*v, vv)
		}
	}
	return err
}

type BioReader interface {
	Rio(dt Dtype, r io.Reader) error
}

type BioWriter interface {
	Wio(dt Dtype, w io.Writer) error
}

type BioMarshaler interface {
	MarshalBio(dt Dtype) (data []byte, err error)
}

type BioUnmarshaler interface {
	UnmarshalBio(dt Dtype, data []byte) error
}

func Read(dt Dtype, r io.Reader, ptr interface{}) error {
	size := 1
	for _, dim := range dt.Dims() {
		size *= dim
	}
	for i := 0; i < size; i++ {
		vv := reflect.New(dt.Type()).Elem()
		err := binary.Read(r, dt.Order(), vv.Addr().Interface())
		if err != nil {
			return err
		}
		fmt.Printf(">>> data[%d]= %#v\n", i, vv.Interface())
	}
	return nil
}
