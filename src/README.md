# FDS-RPC-GO

An RPC in Golang that supports Synchronous and Asynchronous requests, measures performance through Load Tests and outputs graphs and raw data. Created for Fundementals of Distributed Systems.

## Server Usage

### Building the Program

1. cd into the server folder: cd server
2. build server using: go build main.go server_handler.go

### Running the Program

1. run ./main \<server:port>
* Example: ./main localhost:1234

You should see a log confirming that the JSON-RPC is listening.

Ensure that the \<server:port> is the same being used by the client.

## Client Usage

### Building the Program

1. cd into the client folder
2. build the client by using: go build main.go client_handler.go analyze.go

### Running the Program

Client comes with many of options to run:

* <b>-h</b>: 
    * Displays a help usage message
* <b>-s</b>: 
    * Run a small batch of synchronus requests to the server
    * Format: ./main -s \<server:port> \<Seed>
    * Example: ./main -s localhost:1234 1
* <b>-a</b>:
    * Run a smalll batch of asynchronus requests to the server
    * Format: ./main -a \<server:port>
* <b>-lt</b>: 
    * Conduct a Single Load test and add the data to a file
    * Format: ./main -lt \<server:port> \<Rate> \<Duration> \<Seed> \<Mode> \<HeavyMix%> \<ResultFileName>
    * Descriptions:
        * \<Rate> --> The number of requests per second
        * \<Duration> --> THe number of seconds to run the load test for
        * \<Seed> --> A randomness seed used to determine the mix of operations and/or the size of the request selected (Big or small)
        * \<Mode> --> Changes the mix of the requests issued (values 0 - 4)
            * 0 --> Mixed Operations
            * 1 --> String Hashing Only
            * 2 --> Matrix Multiplication Only
            * 3 --> ZlibCompression Only
            * 4 --> Array Sort Only
        * \<HeavyMix%> -->  val from 0 to 100, percentage chance of requests that are "heavy"
        * \<ResultFileName> --> where the results of the loadtest will be stored in json format (The file does not have to exist prior to running, it will be created if it does not exist)
    * Example: ./main -lt localhost:1234 10 5 1 0 25 result
        * Run the load test at localhost port 1234 doing 10 requests per second for 5 seconds. use the randomness seed 1 and mode 0 to mix the operations sent. Let there be a 25% percentage chance of heavy instructions per each instruction. Store the results in the file results.jsonl.
* <b>-g</b>:
    * Create graphs Average, 50th Percentile, 95th Percentile, 88th Percentile for a conducted Load Test
    * Format: ./main -g \<filename>
* <b>-pg</b>:
    * Print the summary data used to create Load Test graphs to the console
    * Format ./main -pg \<filename>

Ensure that the \<server:port> is the same being used as the server.

## Pre-Prepared Load Tests

Pre-prepared Load test sequences have are available if you dont want to craft your own using -lt. All of these have the format:  ./main -lt# \<server:port>. Ensure that the server is active.

Options:
* -lt1 --> Rates from 100 to 2000 requests per second increasing in intervals of 100 req/s. Lasts one second for every request mode, zero chance of large requests. Stores results in load_test_eg1.jsonl
* -lt2 --> Rates from 400 to 1200 requests per second increasing in intervals of 100 req/s. Lasts one second for every request mode, zero chance of large requests. Stores results in load_test_eg2.jsonl
* -lt3 --> Rates from 400 to 1200 requests per second increasing in intervals of 100 req/s. Lasts one second for every request mode, fifty percent chance of large requests. Stores results in load_test_eg3.jsonl
* -lt4 --> Rates from 400 to 1200 requests per second increasing in intervals of 100 req/s. Lasts one second for every request mode, 100% chance of large requests. Stores results in load_test_eg4.jsonl

<b>Note:</b> The results files do not clear ofter running, if you want a clean slate either delete or rename the pre-prepared load test files!

## Design

This RPC uses TCP to send requests over a socket. Requests are Golang structs converted into JSON using functions from the /net/rpc/jsonrpc package. Under the hood the json is structured as followed:

`
{
  "method": <method name (str)>,
  "params": <Array of Parameter Map>,
}
`

For example:

`
{
  "method": "ArraySort.SortArray",
  "params": [{"arr1":[1, 2, 3, 4],"arr2":[5, 6, 7, 8],"size":2}],
}
`

The methods available and parameters are as follows:
### GetHash.HashCompute
Params:
* data: byte array
* size: int

### ArraySort.SortArray
Params:
* data: int32 array
* size: int

### MatrixMultiply.MultiplyMatrix
Params:
* arr1: float64 array
* arr2: float64 array
* size: int

### Zlib.ZlibCompress
Params:
* data: byte array
* size: int

### Zlib.ZlibDecompress (Only Appears in Small Sequential Test)
Params:
* data: byte array
* size: int

## Testing

To run the test in each module simply navigate to the folder and run 
`go test`
