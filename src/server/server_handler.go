package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/instrumentation_export"
)

const helpMessage = `
Usage:
  ./main <server:port>
  `

// Compute Hash

type HashArgs struct {
	Data []byte `json:"data"`
	Size int    `json:"size"`
}

type GetHash struct{}

func (gh GetHash) HashCompute(args HashArgs, reply *string) error {
	hash := sha256.Sum256(args.Data)
	*reply = hex.EncodeToString(hash[:])
	return nil
}

// Array sort

type SortArgs struct {
	Data []int32 `json:"data"`
	Size int     `json:"size"`
}

type ArraySort struct{}

func (as ArraySort) SortArray(args SortArgs, reply *[]int32) error {
	*reply = quicksort(args.Data)
	return nil
}

func quicksort(arr []int32) []int32 {
	if len(arr) < 2 {
		return arr
	}

	left, right := 0, len(arr)-1

	// Choose a pivot (here we pick the middle element)
	pivotIndex := len(arr) / 2
	arr[pivotIndex], arr[right] = arr[right], arr[pivotIndex]

	// Partition
	for i := range arr {
		if arr[i] < arr[right] {
			arr[i], arr[left] = arr[left], arr[i]
			left++
		}
	}

	// Put pivot into correct place
	arr[left], arr[right] = arr[right], arr[left]

	// Recursively sort left and right partitions
	quicksort(arr[:left])
	quicksort(arr[left+1:])

	return arr
}

// Multiply matricies

type MatMutArgs struct {
	Arr1 []float64 `json:"arr1"`
	Arr2 []float64 `json:"arr2"`
	Size int       `json:"size"`
}

type MatrixMultiply struct{}

func (mm MatrixMultiply) MultiplyMatrix(args MatMutArgs, reply *[]float64) error {
	A := args.Arr1
	B := args.Arr2
	n := args.Size
	if len(A) != n*n || len(B) != n*n {
		log.Printf("Matrix is not a square of size %dx%d\n", n, n)
		msg := fmt.Sprintf("Matrix is not a square of size %dx%d\n", n, n) // TODO: Make this a log
		return errors.New(msg)
	}

	C := make([]float64, n*n)
	var sum float64
	for i := 0; i < n; i++ { // row in A
		for j := 0; j < n; j++ { // col in B
			sum = 0
			for k := 0; k < n; k++ { // dot product
				sum += A[i*n+k] * B[k*n+j]
			}
			C[i*n+j] = sum
		}
	}
	*reply = C

	return nil
}

// compress data

type ZlibArgs struct {
	Data []byte `json:"data"`
	Size int    `json:"size"`
}

type Zlib struct{}

func (zc Zlib) ZlibCompress(args ZlibArgs, reply *[]byte) error {
	var b bytes.Buffer

	// Create a new zlib writer
	w := zlib.NewWriter(&b)

	// Write data to it
	_, err := w.Write(args.Data)
	if err != nil {
		return err
	}

	// Close to flush all data
	w.Close()

	*reply = b.Bytes()
	return nil
}

func (zc Zlib) ZlibDecompress(args ZlibArgs, reply *[]byte) error {
	r, err := zlib.NewReader(bytes.NewReader(args.Data))
	if err != nil {
		return err
	}
	var out bytes.Buffer
	io.Copy(&out, r)
	r.Close()

	*reply = out.Bytes()
	return nil
}

type Shutdown struct{}

func (s *Shutdown) Exit(args struct{}, reply *string) error {
	log.Println("Shutdown requested. Dumping instrumentation logs...")
	// instrumentation_export.DumpInstrumentationLogs() // <-- YOUR FUNCTION
	// instrumentation_export.DumpQSizeLogs()
	// instrumentation_export.DumpGStatusLogs()
	// instrumentation_export.DumpCyclesLogs()
	instrumentation_export.DumpCyclesLogsToFile("../json_results/cycles_events.jsonl")
	instrumentation_export.DumpInstrumentationLogsToFile("../json_results/instrumentation.jsonl")
	instrumentation_export.DumpGStatusLogsToFile("../json_results/goroutine_status.jsonl")
	instrumentation_export.DumpQSizeLogsToFile("../json_results/queue_size.jsonl")
	*reply = "Server shutting down and logs dumped."
	go func() {
		// Give the RPC response a moment to be sent before exiting
		// Otherwise client might not receive it
		os.Exit(0)
	}()
	return nil
}
