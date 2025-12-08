package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"

	"runtime/instrumentation_export"
	"syscall"
)

func help() {
	fmt.Println(helpMessage)
}

func main() {
	argsLen := len(os.Args)
	if argsLen == 2 {
		rpc.Register(new(GetHash))
		rpc.Register(new(ArraySort))
		rpc.Register(new(MatrixMultiply))
		rpc.Register(new(Zlib))

		listener, err := net.Listen("tcp", os.Args[1])
		if err != nil {
			log.Fatal("Listen error:", err)
		}
		log.Println("JSON-RPC server listening on: ", os.Args[1])

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigCh
			log.Println("Shutdown requested. Dumping instrumentation logs...")
			instrumentation_export.DumpInstrumentationLogs() // <-- YOUR FUNCTION
			instrumentation_export.DumpQSizeLogs()
			instrumentation_export.DumpGStatusLogs()
			instrumentation_export.DumpInstrumentationLogsToFile("../json_results/instrumentation.json")
			instrumentation_export.DumpGStatusLogsToFile("../json_results/goroutine_status.json")
			instrumentation_export.DumpQSizeLogsToFile("../json_results/queue_size.json")
			os.Exit(0)
		}()

		// -----------------------------------------------------
		// Accept connections forever (until Ctrl-C is pressed)
		// -----------------------------------------------------
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Accept error:", err)
				continue
			}
			go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
		}

	} else {
		help()
	}
}
