// +build ignore

package main

import (
	"encoding/binary"
	"log"
	"os"
)

func main() {
	for _, t := range []struct {
		name string
		size int
	}{
		{
			name: "data_float64_2x3_corder.dat",
			size: 2 * 3,
		},
		{
			name: "data_float64_2x3x4_corder.dat",
			size: 2 * 3 * 4,
		},
	} {
		f, err := os.Create(t.name)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		data := make([]float64, t.size)
		for i := range data {
			data[i] = float64(i + 1)
		}

		err = binary.Write(f, binary.LittleEndian, data)
		if err != nil {
			log.Fatal(err)
		}

		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}
