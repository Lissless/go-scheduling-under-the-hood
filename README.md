# Scheduling Under the Hood: Go RPC Insights

In a disaggregated architecture, computation and data often reside on different nodes, requiring efficient remote communication. Remote Procedure Calls (RPCs) serve as the foundation of this interaction, making their performance central to system efficiency. 

The Go programming language is widely used in modern datacenter environments for building scalable cloud services. It provides a powerful yet abstracted concurrency model using lightweight goroutines. However, the runtime’s internal behavior can significantly affect RPC latency and throughput. Having a better understanding on how much time is spent on runtime mechanisms vs data processing could lead to developments that can improve the performance of Golang implementations.
