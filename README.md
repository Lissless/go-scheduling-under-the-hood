# Scheduling Under the Hood: Go RPC Insights

In a disaggregated architecture, computation and data often reside on different nodes, requiring efficient remote communication. Remote Procedure Calls (RPCs) serve as the foundation of this interaction, making their performance central to system efficiency. 

The Go programming language is widely used in modern datacenter environments for building scalable cloud services. It provides a powerful yet abstracted concurrency model using lightweight goroutines. However, the runtime’s internal behavior can significantly affect RPC latency and throughput. Having a better understanding on how the scheduler acts underneath the surface when faced with different kinds of workloads can help datacenters make decisions on what style would be best for them.

## Before Starting

This project needs local builds of the instrumented verisons of Golang. These are stored in the utils/go-instrumented and utils-preempt/go-preempt-instrumented as submodules to go-scheduling-under-the-hood.

Before proceeding check that tho prior mentioned directories are not empty, if they are run this command from the root:
 * git submodule update --init --recursive

## Go-Instrumented & Aspen-Go-Instrumented

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

## Running Experiments With Scripts

There are serveral scripts within src/scripts that are meant to run pre-designed experiemnts and collect data, there are two ways you can use them:

1) ./run_expr_suite
    * This will run all the experiments and place the results in the experiment-results folder
2) ./run_experiment -exprX
    * Where X is a number from 1 to 4 inclusive, this will run the specific experiment that was specified and place the results in the experiment-results folder

## Running Experiments Manually

If you desire you do have the ability to run your own tests manually.

1) Start the server with the environment variable GOINSTRUMENT
    * GOINSTRUMENT=1 ./main \<server:port>
    * Example: GOINSTRUMENT=1 ./main localhost:1234
        * If successful you should see the message "Goroutine instrumentation enabled: Cooperative/Preemptive" depending on which kind of go comiled the server
2) Run the client with the desired test, many different tests are detailed in src/README.md
    * Example: ./main -expr1 localhost:1234
3) Close the server after the desired tests are done, this will produce serveral files of json with data
    * the -exprX tests come with a shutdown instruction, others may not. If the server is still up when you are done just use Constrol-C
4) Run the desired analysis module on the data using the client
    * Example: ./main -inst ../json_results/instrumentation.jsonl
    * Example: ./main -gstat ../json_results/goroutine_status.jsonl
    * Example: ./main -cycle ../json_results/cycles_events.jsonl
    * Resulting graphs will be saved in the file taht ran the client application.