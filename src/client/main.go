package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

/*

Load Test Helper Functions

*/

func sendHashLoadTest(cfg LoadConfig, stateSeed int64) error {
	conn, err := net.Dial("tcp", cfg.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := jsonrpc.NewClient(conn)

	randGen := rand.New(rand.NewSource(stateSeed))
	choice := randGen.Intn(100 - (0 + 1)) // rand int between 0 and 100

	var hashReply string // there may be a race here lol
	var hargs HashArgs
	if choice < cfg.HeavyMix {
		// do a heavy operation
		hargs = HashArgs{[]byte("As the blue one says, Gotta go fast"), 14}
	} else {
		hargs = HashArgs{[]byte(LARGE_TEXT), 14}
	}

	call := client.Go("GetHash.HashCompute", hargs, &hashReply, nil)

	<-call.Done // wait for response
	return call.Error
}

func sendMatMuxLoadTest(cfg LoadConfig, stateSeed int64) error {
	conn, err := net.Dial("tcp", cfg.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := jsonrpc.NewClient(conn)

	randGen := rand.New(rand.NewSource(stateSeed))
	choice := randGen.Intn(100 - (0 + 1)) // rand int between 0 and 100
	var matmutArgs MatMutArgs
	var arr1, arr2, matReply []float64

	if choice < cfg.HeavyMix {
		arr1 = LARGE_ARR1
		arr2 = LARGE_ARR2
		matmutArgs = MatMutArgs{arr1, arr2, 6}
	} else {
		arr1 = []float64{1, 2, 3, 4}
		arr2 = []float64{5, 6, 7, 8}
		matmutArgs = MatMutArgs{arr1, arr2, 2}
	}

	call := client.Go("MatrixMultiply.MultiplyMatrix", matmutArgs, &matReply, nil)

	<-call.Done // wait for response
	return call.Error
}

func sendZlibCompressLoadTest(cfg LoadConfig, stateSeed int64) error {
	conn, err := net.Dial("tcp", cfg.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := jsonrpc.NewClient(conn)

	randGen := rand.New(rand.NewSource(stateSeed))
	choice := randGen.Intn(100 - (0 + 1)) // rand int between 0 and 100

	var compResp []byte
	var compArgs ZlibArgs

	if choice < cfg.HeavyMix {
		compArgs = ZlibArgs{[]byte(LARGE_TEXT), len(LARGE_TEXT)}
	} else {
		compArgs = ZlibArgs{[]byte("I am crushed and reborn!"), 24}
	}

	call := client.Go("Zlib.ZlibCompress", compArgs, &compResp, nil)

	<-call.Done // wait for response
	return call.Error
}

func sendArraySortLoadTest(cfg LoadConfig, stateSeed int64) error {
	conn, err := net.Dial("tcp", cfg.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := jsonrpc.NewClient(conn)

	randGen := rand.New(rand.NewSource(stateSeed))
	choice := randGen.Intn(100 - (0 + 1)) // rand int between 0 and 100

	var sortData, sortReply []int32

	if choice < cfg.HeavyMix {
		sortData = LARGE_ARR300
	} else {
		sortData = []int32{1, 5, 9, 27, 3, 5, 8, 1, 9, 7, 11}
	}

	sargs := SortArgs{sortData, len(sortData)}
	call := client.Go("ArraySort.SortArray", sargs, &sortReply, nil)

	<-call.Done // wait for response
	return call.Error
}

/*

Main Functions

*/

func sendSync(serverAddr string) {
	// Connect to the server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal("Dialing:", err)
	}
	defer conn.Close()

	// Create JSON-RPC client
	client := jsonrpc.NewClient(conn)
	defer client.Close()

	log.Println("Starting Small Synchronus set of requests")

	// Get Hash
	phrase := "As the blue one says, Gotta go fast"
	log.Printf("Hashing the phrase: %s\n", phrase)
	hargs := HashArgs{[]byte(phrase), 14}
	var reply string
	err = client.Call("GetHash.HashCompute", hargs, &reply)
	if err != nil {
		log.Fatal("Hash Compute error: ", err)
	}
	log.Printf("The returned hash is: %s\n", reply)

	// Sort array
	sortData := []int32{1, 5, 9, 27, 3, 5, 8, 1, 9, 7, 11}
	log.Printf("Sorting the array %v\n", sortData)
	sargs := SortArgs{sortData, 10}
	var sortReply []int32
	err = client.Call("ArraySort.SortArray", sargs, &sortReply)
	if err != nil {
		log.Fatal("Array Sort error: ", err)
	}
	log.Printf("The returned array is: %v\n", sortReply)

	// Matrix multiply
	arr1 := []float64{1, 2, 3, 4}
	arr2 := []float64{5, 6, 7, 8}
	log.Printf("Multiplying the matricies %v and %v\n", arr1, arr2)
	matmutArgs := MatMutArgs{arr1, arr2, 2}
	var matReply []float64
	err = client.Call("MatrixMultiply.MultiplyMatrix", matmutArgs, &matReply)
	if err != nil {
		log.Fatal("Matrix Multiplication error: ", err)
	}
	log.Printf("The multiplied matrix array is: %v\n", matReply)

	// data compression and decompression
	phrase = "I am crushed and reborn!"
	log.Printf("Zlib compressing the phrase: %s\n", phrase)
	compArgs := ZlibArgs{[]byte(phrase), 24}
	var compResp []byte
	err = client.Call("Zlib.ZlibCompress", compArgs, &compResp)
	if err != nil {
		log.Fatal("Zlib Compression error: ", err)
	}
	log.Printf("The compressed data is is: %v\n", compResp)

	log.Printf("Decomressing the previously compressed Phrase")
	decompArgs := ZlibArgs{compResp, 12}
	var decompResp []byte
	err = client.Call("Zlib.ZlibDecompress", decompArgs, &decompResp)
	if err != nil {
		log.Fatal("Zlib Decompression error: ", err)
	}
	log.Printf("The decompressed data is is: %v\n", string(decompResp))

	log.Println("Finished Small Synchronus set of requests")
}

func sendAsync(serverAddr string, seed int64) {
	// Connect to the server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal("Dialing:", err)
	}
	defer conn.Close()

	// Create JSON-RPC client
	client := jsonrpc.NewClient(conn)
	defer client.Close()

	if seed < 0 {
		seed = 1
	}
	log.Println("Starting Small Asynchronus set of requests with seed: #", seed)
	randGen := rand.New(rand.NewSource(seed))

	// 2) Issue many concurrent asynchronous calls and wait for all results
	var hashReply string
	var matReply []float64
	var compResp []byte
	var sortReply []int32

	var wg sync.WaitGroup
	nCalls := 5 // make this configurable
	wg.Add(nCalls)
	for i := 0; i < nCalls; i++ {
		choice := randGen.Intn(100 - (0 + 1)) // rand int between 0 and 100

		go func(i int) {
			defer wg.Done()
			var callPtr *rpc.Call
			if choice < 25 {
				log.Println("Async Hashing was chosed as call #", i)
				hargs := HashArgs{[]byte("Gotta go fast"), 14}
				callPtr = client.Go("GetHash.HashCompute", hargs, &hashReply, nil)
			} else if choice < 50 {
				log.Println("Matrix Multiplication was chosed as call #", i)
				arr1 := []float64{1, 2, 3, 4}
				arr2 := []float64{5, 6, 7, 8}
				matmutArgs := MatMutArgs{arr1, arr2, 2}
				callPtr = client.Go("MatrixMultiply.MultiplyMatrix", matmutArgs, &matReply, nil)
			} else if choice < 75 {
				log.Println("Zlib Compression was chosed as call #", i)
				compArgs := ZlibArgs{[]byte("I am crushed and reborn!"), 24}
				callPtr = client.Go("Zlib.ZlibCompress", compArgs, &compResp, nil)
			} else {
				log.Println("Array Sorting was chosed as call #", i)
				sortData := []int32{1, 5, 9, 27, 3, 5, 8, 1, 9, 7, 11}
				sargs := SortArgs{sortData, 10}
				callPtr = client.Go("ArraySort.SortArray", sargs, &sortReply, nil)
			}

			select {
			case res := <-callPtr.Done:
				if res.Error != nil {
					log.Printf("call %d error: %v", i, res.Error)
					return
				}
				log.Printf("call %d Returned\n", i)
			case <-time.After(3 * time.Second):
				log.Printf("call %d timed out\n", i)
			}
		}(i)
	}
	wg.Wait()
}

func loadTest(cfg LoadConfig) []Result {
	results := make([]Result, 0, cfg.Rate*int(cfg.Duration.Seconds())) // hold all the latencies, array of size 0 with space to hold rate * second entries
	resultsMu := sync.Mutex{}                                          // mutex to make sure adding to the results array is safe

	var wg sync.WaitGroup
	interval := time.Second / time.Duration(cfg.Rate)
	endTime := time.Now().Add(cfg.Duration)
	ticker := time.NewTicker(interval)

	randGen := rand.New(rand.NewSource(cfg.Seed))
	var thread int64 = 0

	var upper, lower int
	switch cfg.Mode {
	case 0:
		upper = 100
		lower = 0
	case 1: //Hasing only
		upper = 24
		lower = 0
	case 2: // Matrix multiplication only
		upper = 49
		lower = 25
	case 3: // Zlib Compression only
		upper = 74
		lower = 50
	default: // Array sort only
		upper = 100
		lower = 75
	}

	log.Printf("Starting Load test with Parameters: %v\n", cfg)

	for time.Now().Before(endTime) {
		<-ticker.C
		wg.Add(1)
		go func() {
			defer wg.Done()
			choice := randGen.Intn(upper-lower) + lower // rand int between 0 and 100
			if choice < 25 {
				start := time.Now() // start timeing
				err := sendHashLoadTest(cfg, thread)
				lat := time.Since(start) // finish timing to calculate the latency
				resultsMu.Lock()
				thread++
				results = append(results, Result{Latency: lat, Error: err}) //err})
				resultsMu.Unlock()
			} else if choice < 50 {
				start := time.Now() // start timeing
				err := sendMatMuxLoadTest(cfg, thread)
				lat := time.Since(start) // finish timing to calculate the latency
				resultsMu.Lock()
				thread++
				results = append(results, Result{Latency: lat, Error: err}) //err})
				resultsMu.Unlock()
			} else if choice < 75 {
				start := time.Now() // start timeing
				err := sendZlibCompressLoadTest(cfg, thread)
				lat := time.Since(start) // finish timing to calculate the latency
				resultsMu.Lock()
				thread++
				results = append(results, Result{Latency: lat, Error: err}) //err})
				resultsMu.Unlock()
			} else {
				start := time.Now() // start timeing
				err := sendArraySortLoadTest(cfg, thread)
				lat := time.Since(start) // finish timing to calculate the latency
				resultsMu.Lock()
				thread++
				results = append(results, Result{Latency: lat, Error: err}) //err})
				resultsMu.Unlock()
			}
		}()
	}

	ticker.Stop()
	wg.Wait()

	log.Println("Finished Load Test")
	return results
}

func percentiles(results []Result) (p50, p95, p99 float64) {
	var latencies []float64
	for _, r := range results {
		if r.Error == nil {
			latencies = append(latencies, float64(r.Latency.Microseconds()))
		}
	}
	if len(latencies) == 0 {
		return math.NaN(), math.NaN(), math.NaN()
	}
	sort.Float64s(latencies)
	p50 = selectPercentile(latencies, 0.50) / 1000
	p95 = selectPercentile(latencies, 0.95) / 1000
	p99 = selectPercentile(latencies, 0.99) / 1000
	return p50, p95, p99
}

func selectPercentile(data []float64, pct float64) float64 {
	if len(data) == 0 {
		return math.NaN()
	}
	idx := int(float64(len(data)-1) * pct)
	return data[idx]
}

func report(results []Result, cfg LoadConfig) Summary {
	var latencies []float64
	var errors int
	for _, r := range results {
		if r.Error != nil {
			errors++
			continue
		}
		latencies = append(latencies, float64(r.Latency.Microseconds())/1000.0)
	}
	sort.Float64s(latencies)

	sum := 0.0
	for _, v := range latencies {
		sum += v
	}
	avg := sum / float64(len(latencies))
	throughput := float64(len(latencies)) / cfg.Duration.Seconds() //float64(cfg.Duration) //

	var op string
	switch cfg.Mode {
	case 0:
		op = "Mixed Operations"
	case 1:
		op = "String Hashing"
	case 2:
		op = "Matrix Multiplication"
	case 3:
		op = "Zlib Compression"
	default:
		op = "Array Sort"
	}

	p50, p95, p99 := percentiles(results)

	summary := Summary{
		Operation:  op,
		Seed:       cfg.Seed,
		Rate:       cfg.Rate,
		AvgLatency: avg,
		P50Latency: p50,
		P95Latency: p95,
		P99Latency: p99,
		Throughput: throughput,
		Errors:     errors,
	}
	log.Printf("Load Test Summary Results: %v\n", summary)

	f, err := os.OpenFile(cfg.ResultFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // 0644 gives read and write permisisons
	if err != nil {
		log.Println("Unable to open file to write summary record")
		log.Fatal(err)
	}
	err = json.NewEncoder(f).Encode(summary)
	if err != nil {
		log.Println("Summary contained NaN due to low performance, cannot write this record")
	}
	f.Close()

	return summary
}

func help() {
	fmt.Println(helpMessage)
}

// Usage: ./main -a (for async)
// -lt for load test
func main() {
	argsLen := len(os.Args)
	if argsLen > 1 {
		switch os.Args[1] {
		case "-h":
			help()
		case "-a":
			if argsLen == 4 {
				seed, err := strconv.ParseInt(os.Args[3], 10, 64)
				if err != nil {
					help()
					return
				} else {
					sendAsync(os.Args[2], seed)
				}
			} else {
				help()
			}
		case "-s":
			if argsLen == 3 {
				sendSync(os.Args[2])
			} else {
				help()
			}
		case "-lt":
			if argsLen == 9 {

				rate, err := strconv.Atoi(os.Args[3])
				if err != nil {
					help()
					return
				}
				durr, err := strconv.Atoi(os.Args[4])
				if err != nil {
					help()
					return
				}
				seed, err := strconv.ParseInt(os.Args[5], 10, 64)
				if err != nil {
					help()
					return
				}

				mode, err := strconv.Atoi(os.Args[6])
				if err != nil {
					help()
					return
				} else if mode < 0 || mode > 4 {
					help()
					return
				}

				heavyMix, err := strconv.Atoi(os.Args[7])
				if err != nil {
					help()
					return
				} else if heavyMix > 100 {
					help()
					return
				}

				config := LoadConfig{os.Args[2], rate, time.Duration(durr) * time.Second, seed, mode, heavyMix, os.Args[8] + ".jsonl"}
				report(loadTest(config), config)
			} else {
				help()
			}
		case "-lt1":
			var config LoadConfig
			if len(os.Args) == 3 {
				log.Println("Processing Load Test, please wait 1 minute!")
				for i := 1; i < 21; i++ {
					config = LoadConfig{os.Args[2], 100 * i, time.Duration(1) * time.Second, 1, 0, 0, "load_test_eg1.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 100 * i, time.Duration(1) * time.Second, 1, 1, 0, "load_test_eg1.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 100 * i, time.Duration(1) * time.Second, 1, 2, 0, "load_test_eg1.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 100 * i, time.Duration(1) * time.Second, 1, 3, 0, "load_test_eg1.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 100 * i, time.Duration(1) * time.Second, 1, 4, 0, "load_test_eg1.jsonl"}
					report(loadTest(config), config)
				}
				log.Println("Finished Processing test, Results in load_test_eg1.jsonl")
			}
		case "-lt2":
			// "localhost:1234"
			var config LoadConfig
			if len(os.Args) == 3 {
				log.Println("Processing Load Test, please wait 1 minute!")
				for i := 1; i < 10; i++ {
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 0, 0, "load_test_eg2.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 1, 0, "load_test_eg2.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 2, 0, "load_test_eg2.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 3, 0, "load_test_eg2.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 4, 0, "load_test_eg2.jsonl"}
					report(loadTest(config), config)
				}
				log.Println("Finished Processing test, Results in load_test_eg2.jsonl")
			}
		case "-lt3":
			var config LoadConfig
			if len(os.Args) == 3 {
				log.Println("Processing Load Test, please wait 1 minute!")
				for i := 1; i < 10; i++ {
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 0, 50, "load_test_eg3.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 1, 50, "load_test_eg3.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 2, 50, "load_test_eg3.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 3, 50, "load_test_eg3.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 4, 50, "load_test_eg3.jsonl"}
					report(loadTest(config), config)
				}
				log.Println("Finished Processing test, Results in load_test_eg3.jsonl")
			}
		case "-lt4":
			var config LoadConfig
			if len(os.Args) == 3 {
				log.Println("Processing Load Test, please wait 1 minute!")
				for i := 1; i < 10; i++ {
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 0, 100, "load_test_eg4.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 1, 100, "load_test_eg4.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 2, 100, "load_test_eg4.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 3, 100, "load_test_eg4.jsonl"}
					report(loadTest(config), config)
					config = LoadConfig{os.Args[2], 300 + (100 * i), time.Duration(1) * time.Second, 1, 4, 100, "load_test_eg4.jsonl"}
					report(loadTest(config), config)
				}
				log.Println("Finished Processing test, Results in load_test_eg4.jsonl")
			}
		case "-lt5":
			var config LoadConfig
			if len(os.Args) == 3 {
				log.Println("Processing Load Test, please wait 1 minute!")
				// TODO: Fine a way that is somewhat near runtime.nanotime so we an filter out interference
				// 			in loggin gbefore deploying the workoad
				start := time.Now() // for linux this may be monotonic?
				print("Start Time: ", time.Since(start).Nanoseconds(), "\n")
				config = LoadConfig{os.Args[2], 20, time.Duration(10) * time.Second, 1, 0, 0, "load_test_eg5.jsonl"}
				report(loadTest(config), config)
				print("End Time: ", time.Since(start).Nanoseconds(), "\n")
				log.Println("Finished Processing test, Results in load_test_eg1.jsonl")
			}
		case "-g":
			if len(os.Args) == 3 {
				data, err := getSummaryData(os.Args[2])
				if err != nil {
					log.Fatalf("failed reading the input file")
				}
				makeGraphs(data)
				if err != nil {
					log.Fatalf("failed making graphs")
				}
			}
		case "-inst":
			if len(os.Args) == 3 {
				data, err := getInstrumentationData(os.Args[2])
				if err != nil {
					log.Fatalf("failed reading the input file")
				}
				makeCreationLatencyHistogram(data)
				makeCreationLatencyCDF(data)
				makeGoroutinesCreated(data)
				if err != nil {
					log.Fatalf("failed making graphs")
				}
			}
		case "-gstat":
			if len(os.Args) == 3 {
				datags, err := getGStatusData(os.Args[2])
				if err != nil {
					log.Fatalf("failed reading the input file")
				}

				makeSchedulingLatencyCDF(datags)
				makeSchedulingLatencyHistogram(datags)
				if err != nil {
					log.Fatalf("failed making graphs")
				}
			}
		case "-pg":
			if len(os.Args) == 3 {
				data, err := getSummaryData(os.Args[2])
				if err != nil {
					log.Fatalf("failed reading the input file")
				}
				printSummary(data)
				if err != nil {
					log.Fatalf("failed printing graphs")
				}
			}

		}
		// runtime.DumpCreationLogs()
		// runtime.DumpTimingLogs()
	}
}
