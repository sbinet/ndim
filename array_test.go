package ndim

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestReadArray(t *testing.T) {
	f, err := os.Open("testdata/data_float64_2x3_corder.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var data []float64
	dec, err := NewDecoder(f)
	if err != nil {
		t.Fatal(err)
	}

	err = dec.Decode(&data)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("data: %v\n", data)

	const size = 2 * 3
	if len(data) != size {
		log.Fatalf("invalid size. got=%d. want=%d\n", len(data), size)
	}

}

func TestReadArray2D(t *testing.T) {
	f, err := os.Open("testdata/data_float64_2x3_corder.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var data ndArray
	dec, err := NewDecoder(f)
	if err != nil {
		t.Fatal(err)
	}

	err = dec.Decode(&data)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("data: %v\n", data)

	const size = 2 * 3
	if data.Size() != size {
		log.Fatalf("invalid size. got=%d. want=%d\n", data.Size(), size)
	}

	if !reflect.DeepEqual(data.Dims(), []int{2, 3}) {
		log.Fatalf("invalid dims. got=%v. want=%v\n", data.Dims(), []int{2, 3})
	}
}

func TestReadArray3D(t *testing.T) {
	f, err := os.Open("testdata/data_float64_2x3x4_corder.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var data ndArray
	dec, err := NewDecoder(f)
	if err != nil {
		t.Fatal(err)
	}

	dec.hdr.shape = []int{2, 3, 4}
	err = dec.Decode(&data)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("data: %v\n", data)

	const size = 2 * 3 * 4
	if data.Size() != size {
		log.Fatalf("invalid size. got=%d. want=%d\n", data.Size(), size)
	}

	if !reflect.DeepEqual(data.Dims(), []int{2, 3, 4}) {
		log.Fatalf("invalid dims. got=%v. want=%v\n", data.Dims(), []int{2, 3, 4})
	}

	f.Seek(0, 0)
	err = Read(dec.hdr, f, nil)
	if err != nil {
		t.Fatal(err)
	}
}

type ndArray struct {
	Data  []float64
	Shape []int
}

func (arr *ndArray) Rio(dt Dtype, r io.Reader) error {
	arr.Shape = make([]int, len(dt.Dims()))
	copy(arr.Shape, dt.Dims())

	size := 1
	for _, dim := range arr.Shape {
		size *= dim
	}

	arr.Data = make([]float64, size)
	return binary.Read(r, dt.Order(), arr.Data)
}

func (arr *ndArray) Size() int {
	sz := 1
	for _, dim := range arr.Shape {
		sz *= dim
	}
	return sz
}

func (arr *ndArray) Dims() []int {
	return arr.Shape
}
