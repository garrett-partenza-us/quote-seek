package main

import (
	"encoding/binary"
	"bytes"
	"log"
	"os"
)

func ReadBinaryArrayFile1D(path string, cols int) []float64 {

	// Read the binary array file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Get the file size in bytes
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := fileInfo.Size()

	// Ensure file bytes is multiple of float64 byte size
	if fileSize%8 != 0 {
		log.Fatal("The file size is not a multiple of 4 bytes (size of float64)")
	}

	// Read the binary data into a byte slice
	data := make([]byte, fileSize)
	_, err = file.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the byte slice into a slice of float64
	numFloats := fileSize / 8
	floats := make([]float64, numFloats)
	err = binary.Read(bytes.NewReader(data), binary.LittleEndian, &floats)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the number of elements matches the expected 2D array size
	if cols != len(floats) {
		log.Print("Number of columns expected: ", cols)
		log.Print("Number of columns recieved: ", len(floats))
		log.Fatal("The dimensions of the array do not match the expected size.")
	}

	return floats
}

type StandardScaler struct {
	Mean		[]float64
	Scale		[]float64
}

func NewStandardScaler(size int, pathMean, pathScale string) StandardScaler {

	mean := ReadBinaryArrayFile1D(pathMean, size)
	scale := ReadBinaryArrayFile1D(pathScale, size)

	return StandardScaler{
		Mean:  mean,
		Scale: scale,
	}

}

func (s *StandardScaler) ScaleVector(data []float64) []float64 {
	scaledData := make([]float64, len(data))
	for col, val := range data {
		scaledData[col] = (val - s.Mean[col]) / s.Scale[col]
	}
	return scaledData
}
