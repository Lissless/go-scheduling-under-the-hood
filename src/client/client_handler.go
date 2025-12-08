package main

import (
	"time"
)

const helpMessage = `
Usage:
  ./main [option] [arguments]

Options:
  -s:
    Run a small batch of synchronous requests to the server.
    Format:  ./main -s <server:port>
    Example: ./main -s localhost:1234

  -a:
    Run a small batch of asynchronous requests to the server.
    Format:  ./main -a <server:port>
    Example: ./main -a localhost:1234

  -lt:
    Conduct a single load test and add the data to a file.
    Format:  ./main -lt <server:port> <Rate> <Duration> <Seed> <Mode> <HeavyMix%> <ResultFileName>

    Descriptions:
      <Rate>          The number of requests per second.
      <Duration>      The number of seconds to run the load test for.
      <Seed>          A randomness seed used to determine the mix of operations and/or request size.
      <Mode>          Changes the mix of the requests issued (values 0 - 4):
                        0 → Mixed Operations
                        1 → String Hashing Only
                        2 → Matrix Multiplication Only
                        3 → Zlib Compression Only
                        4 → Array Sort Only
      <HeavyMix%>     A value from 0 to 100, indicating the percentage chance of "heavy" requests.
      <ResultFileName>  The JSONL file where results will be stored.
                        (Created if it does not exist.)

    Example:
      ./main -lt localhost:1234 10 5 1 0 25 result
      → Runs a load test at localhost:1234 doing 10 requests/sec for 5 seconds.
        Uses seed=1, mode=0 (mixed operations), with 25% heavy requests.
        Results saved to result.jsonl.

  -g:
    Create graphs (Average, 50th, 95th, and 99th Percentiles) for a conducted load test.
    Format:  ./main -g <filename>

  -pg:
    Print the summary data used to create load test graphs to the console.
    Format:  ./main -pg <filename>

Pre-Prepared Load Tests:
   -lt1 --> Rates from 100 to 2000 requests per second increasing in intervals of 100 req/s. Lasts one second for every request mode, zero chance of large requests. 
   			Stores results in load_test_eg1.jsonl
   -lt2 --> Rates from 400 to 1200 requests per second increasing in intervals of 100 req/s. Lasts one second for every request mode, zero chance of large requests. 
			Stores results in load_test_eg2.jsonl
   -lt3 --> Rates from 400 to 1200 requests per second increasing in intervals of 100 req/s. Lasts one second for every request mode, fifty percent chance of large requests. 
			Stores results in load_test_eg3.jsonl
   -lt4 --> Rates from 400 to 1200 requests per second increasing in intervals of 100 req/s. Lasts one second for every request mode, 100 percent chance of large requests. 
			Stores results in load_test_eg4.jsonl

	`

const LARGE_TEXT = "I wanna be the very best\nLike no one ever was\nTo catch them is my real test\nTo train them is my cause\nI will travel across the land\nSearching far and wide\nTeach Pokémon to understand\nThe power that's inside\n\n[Chorus]\n(Pokémon\nGotta catch 'em all) It's you and me\nI know it's my destiny (Pokémon)\nOh, you're my best friend\nIn a world we must defend (Pokémon\nGotta catch 'em all) A heart so true\nOur courage will pull us through\nYou teach me and I'll teach you (Ooh, ooh)\nPokémon! (Gotta catch 'em all)\nGotta catch 'em all\nYeah\n\n[Verse 2]\nEvery challenge along the way\nWith courage, I will face\nI will battle every day\nTo claim my rightful place\nCome with me, the time is right\nThere's no better team\nArm in arm, we'll win the fight\nIt's always been our dream\nSee upcoming rock shows\n\nGotta catch 'em all) It's you and me\nI know it's my destiny (Pokémon)\nOh, you're my best friend\nIn a world we must defend (Pokémon\nGotta catch 'em all) A heart so true\nOur courage will pull us through\nYou teach me and I'll teach you (Ooh, ooh)\nPokémon! (Gotta catch 'em all)\nGotta catch 'em all\n[Bridge]\nGotta catch 'em all\nGotta catch 'em all\nGotta catch 'em all\nYeah\n[Guitar Solo]"

var LARGE_ARR1 = []float64{
	3.2, 87.6, 42.1, 19.9, 64.3, 55.8, 92.4, 11.7,
	76.9, 28.4, 35.6, 81.2, 47.3, 68.7, 24.5, 59.1,
	95.0, 14.8, 33.9, 72.4, 49.5, 61.7, 7.3, 85.6,
	38.2, 57.9, 93.8, 22.5, 66.1, 31.4, 78.7, 9.6,
	72, 55, 72, 1,
}

var LARGE_ARR2 = []float64{
	45.5, 18.3, 97.9, 26.7, 63.2, 88.6, 53.4, 32.8,
	70.1, 11.5, 82.9, 39.4, 58.7, 94.2, 21.6, 75.8,
	28.9, 67.3, 49.1, 84.7, 15.2, 60.9, 34.5, 91.4,
	43.6, 79.8, 25.4, 56.2, 99.3, 12.7, 73.5, 41.9,
	3, 7, 9, 12,
}

var LARGE_ARR300 = []int32{
	12, 57, 893, 44, 670, 381, 952, 278, 135, 749,
	23, 417, 699, 84, 963, 502, 248, 731, 53, 819,
	650, 104, 927, 311, 569, 448, 239, 755, 642, 390,
	872, 501, 190, 978, 615, 322, 708, 452, 89, 937,
	581, 465, 236, 871, 320, 741, 667, 275, 902, 123,
	748, 692, 239, 864, 527, 307, 780, 62, 950, 488,
	815, 373, 561, 199, 832, 91, 771, 405, 286, 978,
	150, 365, 732, 620, 947, 308, 176, 812, 274, 493,
	590, 144, 802, 463, 336, 990, 126, 513, 677, 820,
	92, 260, 547, 194, 725, 301, 669, 158, 949, 786,
	572, 341, 189, 915, 732, 405, 67, 953, 214, 879,
	307, 641, 857, 120, 478, 776, 209, 584, 95, 839,
	458, 370, 699, 147, 954, 621, 499, 311, 730, 167,
	451, 688, 905, 273, 564, 132, 741, 419, 980, 348,
	805, 287, 932, 513, 182, 760, 96, 821, 468, 605,
	135, 993, 207, 728, 410, 599, 334, 879, 460, 655,
	237, 975, 570, 803, 290, 610, 471, 342, 964, 284,
	739, 155, 812, 624, 197, 952, 486, 763, 332, 819,
	497, 142, 957, 306, 888, 520, 119, 671, 450, 943,
	214, 398, 721, 501, 92, 869, 303, 695, 571, 410,
	940, 483, 805, 227, 639, 95, 748, 360, 176, 879,
	539, 668, 285, 742, 182, 953, 419, 893, 310, 504,
	231, 740, 621, 477, 155, 896, 319, 781, 241, 933,
	365, 823, 698, 270, 592, 455, 711, 185, 949, 334,
	890, 639, 235, 479, 764, 125, 904, 218, 567, 351,
	832, 687, 492, 312, 956, 147, 724, 590, 200, 842,
	403, 668, 219, 910, 486, 177, 798, 272, 945, 381,
	154, 857, 691, 502, 130, 733, 407, 576, 289, 996,
	615, 418, 943, 290, 812, 536, 173, 867, 392, 665,
	203, 949, 127, 703, 372, 591, 275, 879, 486, 741,
}

type HashArgs struct {
	Data []byte `json:"data"`
	Size int    `json:"size"`
}

type SortArgs struct {
	Data []int32 `json:"data"`
	Size int     `json:"size"`
}

type MatMutArgs struct {
	Arr1 []float64 `json:"arr1"`
	Arr2 []float64 `json:"arr2"`
	Size int       `json:"size"`
}

type ZlibArgs struct {
	Data []byte `json:"data"`
	Size int    `json:"size"`
}

type LoadConfig struct {
	Address    string        // server address, e.g. "localhost:1234"
	Rate       int           // requests per second
	Duration   time.Duration // how long to run
	Seed       int64         // randomness seed
	Mode       int           // what mix of requests to have
	HeavyMix   int           // val from 0 to 100, percentage chance of requests that are "heavy"
	ResultFile string        // the location where the results of the load test will go
}

type Result struct {
	Latency time.Duration
	Error   error
}

type Summary struct {
	Operation  string  `json:"operation"`
	Seed       int64   `json:"seed"`
	Rate       int     `json:"rate"`   // requests per second
	AvgLatency float64 `json:"avg_ms"` // average in ms
	P50Latency float64 `json:"p50_ms"` // median in ms
	P95Latency float64 `json:"p95_ms"`
	P99Latency float64 `json:"p99_ms"`
	Throughput float64 `json:"throughput"` // successful req/s
	Errors     int     `json:"errors"`
}
