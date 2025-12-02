package main

import "runtime"

func main() {
	print("hello world\n")
	runtime.DumpCreationLogs()
}