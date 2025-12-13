package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"runtime/instrumentation_export"
)

func help() {
	fmt.Println(helpMessage)
}

func main() {
	argsLen := len(os.Args)
	if argsLen == 2 {
		log.SetOutput(os.Stdout)
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
		rpc.Register(new(GetHash))
		rpc.Register(new(ArraySort))
		rpc.Register(new(MatrixMultiply))
		rpc.Register(new(Zlib))
		rpc.Register(new(Shutdown))

		listener, err := net.Listen("tcp", os.Args[1])
		if err != nil {
			log.Fatal("Listen error:", err)
		}
		log.Println("Goroutine instrumentation enabled: ", instrumentation_export.ReturnSchedulerType())
		log.Println("JSON-RPC server listening on: ", os.Args[1])

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
