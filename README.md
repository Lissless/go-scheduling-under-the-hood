# Scheduling Under the Hood: Go RPC Insights

In a disaggregated architecture, computation and data often reside on different nodes, requiring efficient remote communication. Remote Procedure Calls (RPCs) serve as the foundation of this interaction, making their performance central to system efficiency. 

The Go programming language is widely used in modern datacenter environments for building scalable cloud services. It provides a powerful yet abstracted concurrency model using lightweight goroutines. However, the runtime’s internal behavior can significantly affect RPC latency and throughput. Having a better understanding on how much time is spent on runtime mechanisms vs data processing could lead to developments that can improve the performance of Golang implementations.


## Go-Instrumented

In order to re-make the core golang executable for the language you have to:

1) cd into src in the relelvant submodule
2) run ./make.bash


## Compile RPC-FDS-GO

### Build the client
1) cd into src/client
2) run ../../utils/go-instrumented/bin/go build main.go client_handler.go analyze.go

### Build the server
1) cd into src/server
2) Build the version depending on what you want to test:
    * Go-Instrumented (Cooperative Sceduling): 
        * ../../utils/go-instrumented/bin/go build main.go server_handler.go
    * Aspen-Go Instrumented (Preemptive Scheduling): 
        * ../../utils-preempt/go-preempt-instrumented/bin/go build main.go server_handler.go

## Running the tests

1) Start the server with the environment variable GOINSTRUMENT
    * GOINSTRUMENT=1 ./main \<server:port>
    * Example: GOINSTRUMENT=1 ./main localhost:1234
2) Run the client with the desired test, many different tests are detailed in src/README.md
    * Example: ./main -lt5 localhost:1234
3) Close the server after the desired tests are done, this will produce serveral files of json with data
    * TODO: include a verbose mode where certain data should also be dumped to the console
4) Run the desired analysis module on the data using the client
    * Example: ./main -inst ../json_results/instrumentation.jsonl
    * Example ./main -gstat ../json_results/goroutine_status.jsonl
    * Resulting graphs will be saved in the TODO: make all graphs appear in a certain file