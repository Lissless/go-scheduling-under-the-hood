package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"sort"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func checkFile(filePath string) error {
	check := strings.Split(filePath, ".")
	if len(check) < 2 {
		// log.Fatalf("File was not of correct type: .jsonl")
		return fmt.Errorf("file was not of correct type: .jsonl")
	} else {
		if check[len(check)-1] != "jsonl" {
			// log.Fatalf("File was not of correct type: .jsonl")
			return fmt.Errorf("file was not of correct type: .jsonl")
		}
	}
	return nil
}

func makeThrouputAnalysis(grouped map[string][]Summary, colors []color.Color) {
	p2 := plot.New()
	p2.Title.Text = "Throughput Analysis"
	p2.X.Label.Text = "Offered Load (req/s)"
	p2.Y.Label.Text = "Achieved Throughput (req/s)"
	p2.Add(plotter.NewGrid())

	i := 0
	for op, list := range grouped {
		points := make(plotter.XYs, len(list))
		for j, s := range list {
			points[j].X, points[j].Y = float64(s.Rate), s.Throughput
		}
		line, _ := plotter.NewLine(points)
		line.Color = colors[i%len(colors)]
		p2.Add(line)
		p2.Legend.Add(op, line)
		i++
	}

	p2.Legend.Top = true
	p2.Legend.Left = true

	if err := p2.Save(9*vg.Inch, 5*vg.Inch, "throughput_analysis.png"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saved: throughput_analysis.png")
}

func makeAvgLoadLatency(grouped map[string][]Summary, colors []color.Color) {
	p := plot.New()
	p.Title.Text = "Load-Latency Curve (Average)"
	p.X.Label.Text = "Request Rate (req/s)"
	p.Y.Label.Text = "Latency (ms)"
	p.Add(plotter.NewGrid())

	i := 0
	for op, list := range grouped {
		avg := make(plotter.XYs, len(list))
		for j, s := range list {
			avg[j].X, avg[j].Y = float64(s.Rate), s.AvgLatency
		}

		col := colors[i%len(colors)]
		// Average line (solid)
		lineAvg, _ := plotter.NewLine(avg)
		lineAvg.Color = col
		lineAvg.Width = vg.Points(1.2)
		lineAvg.Dashes = []vg.Length{}
		p.Add(lineAvg)
		p.Legend.Add(op, lineAvg)

		i++
	}
	p.Legend.Top = true
	p.Legend.Left = true

	if err := p.Save(9*vg.Inch, 5*vg.Inch, "load_latency_avg_curve.png"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saved: load_latency_avg_curve.png")
}

func makeP95LoadLatency(grouped map[string][]Summary, colors []color.Color) {
	p := plot.New()
	p.Title.Text = "Load-Latency Curve (P95)"
	p.X.Label.Text = "Request Rate (req/s)"
	p.Y.Label.Text = "Latency (ms)"
	p.Add(plotter.NewGrid())

	i := 0
	for op, list := range grouped {
		p95 := make(plotter.XYs, len(list))
		for j, s := range list {
			p95[j].X, p95[j].Y = float64(s.Rate), s.P95Latency
		}

		col := colors[i%len(colors)]
		// P95
		lineP95, _ := plotter.NewLine(p95)
		lineP95.Color = col
		lineP95.Width = vg.Points(1.2)
		lineP95.Dashes = []vg.Length{}
		p.Add(lineP95)
		p.Legend.Add(op, lineP95)

		i++
	}
	p.Legend.Top = true
	p.Legend.Left = true

	if err := p.Save(9*vg.Inch, 5*vg.Inch, "load_latency_P95_curve.png"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saved: load_latency_P95_curve.png")

}

func makeP99LoadLatency(grouped map[string][]Summary, colors []color.Color) {
	p := plot.New()
	p.Title.Text = "Load-Latency Curve (P99)"
	p.X.Label.Text = "Request Rate (req/s)"
	p.Y.Label.Text = "Latency (ms)"
	p.Add(plotter.NewGrid())

	i := 0
	for op, list := range grouped {
		p99 := make(plotter.XYs, len(list))
		for j, s := range list {
			p99[j].X, p99[j].Y = float64(s.Rate), s.P99Latency
		}

		col := colors[i%len(colors)]
		// P99
		lineP99, _ := plotter.NewLine(p99)
		lineP99.Color = col
		lineP99.Width = vg.Points(1.2)
		lineP99.Dashes = []vg.Length{}
		p.Add(lineP99)
		p.Legend.Add(op, lineP99)

		i++
	}
	p.Legend.Top = true
	p.Legend.Left = true

	if err := p.Save(9*vg.Inch, 5*vg.Inch, "load_latency_P99_curve.png"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saved: load_latency_P99_curve.png")

}

func makeP50LoadLatency(grouped map[string][]Summary, colors []color.Color) {
	p := plot.New()
	p.Title.Text = "Load-Latency Curve (P50)"
	p.X.Label.Text = "Request Rate (req/s)"
	p.Y.Label.Text = "Latency (ms)"
	p.Add(plotter.NewGrid())

	i := 0
	for op, list := range grouped {
		p50 := make(plotter.XYs, len(list))
		for j, s := range list {
			p50[j].X, p50[j].Y = float64(s.Rate), s.P50Latency
		}

		col := colors[i%len(colors)]
		// P99
		lineP50, _ := plotter.NewLine(p50)
		lineP50.Color = col
		lineP50.Width = vg.Points(1.2)
		lineP50.Dashes = []vg.Length{}
		p.Add(lineP50)
		p.Legend.Add(op, lineP50)

		i++
	}
	p.Legend.Top = true
	p.Legend.Left = true

	if err := p.Save(9*vg.Inch, 5*vg.Inch, "load_latency_P50_curve.png"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saved: load_latency_P50_curve.png")

}

func getFileData(filePath string) ([]Summary, error) {
	var data []Summary
	check := checkFile(filePath)
	if check != nil {
		return nil, fmt.Errorf("the File was invalid type, needs to be: .jsonl")
	}

	f, err := os.Open(filePath)
	if err != nil {
		// log.Fatalf("failed to open file %s: %v", filePath, err)
		return nil, fmt.Errorf("failed to Open file")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// 1️⃣ Read file line by line
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var s Summary
		if err := json.Unmarshal(line, &s); err != nil {
			log.Printf("Skipping invalid line: %v", err)
			continue
		}
		data = append(data, s)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	if len(data) == 0 {
		log.Fatal("No valid records found in results.jsonl")
	}

	return data, nil
}

func makeGraphs(data []Summary) {
	// Group by operation
	grouped := make(map[string][]Summary)
	for _, s := range data {
		grouped[s.Operation] = append(grouped[s.Operation], s)
	}

	// Sort each group by rate
	for _, list := range grouped {
		sort.Slice(list, func(i, j int) bool { return list[i].Rate < list[j].Rate })
	}

	colors := []color.Color{
		color.RGBA{255, 99, 132, 255},  // red
		color.RGBA{54, 162, 235, 255},  // blue
		color.RGBA{75, 192, 192, 255},  // teal
		color.RGBA{255, 206, 86, 255},  // yellow
		color.RGBA{153, 102, 255, 255}, // purple
		color.RGBA{255, 159, 64, 255},  // orange
	}

	makeAvgLoadLatency(grouped, colors)
	makeP95LoadLatency(grouped, colors)
	makeP99LoadLatency(grouped, colors)
	makeP50LoadLatency(grouped, colors)

	makeThrouputAnalysis(grouped, colors)
}

func printSummary(data []Summary) {
	grouped := make(map[string][]Summary)
	for _, s := range data {
		grouped[s.Operation] = append(grouped[s.Operation], s)
	}

	// Sort by rate for each operation
	for op, list := range grouped {
		sort.Slice(list, func(i, j int) bool { return list[i].Rate < list[j].Rate })
		grouped[op] = list
	}

	// Print summary per operation
	fmt.Println("\n Summary by Operation:")
	for op, list := range grouped {
		fmt.Printf("\nOperation: %s\n", op)
		fmt.Println("Seed\tRate\tAvg(ms)\tP50(ms)\tP95(ms)\tP99(ms)\tThroughput\tErrors")
		fmt.Println("-------------------------------------------------------------------------")
		for _, s := range list {
			fmt.Printf("%d\t%d\t%.2f\t%.2f\t%.2f\t%.2f\t%.1f\t\t%d\n",
				s.Seed, s.Rate, s.AvgLatency, s.P50Latency, s.P95Latency, s.P99Latency, s.Throughput, s.Errors)
		}
	}
}
