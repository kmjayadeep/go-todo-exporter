package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	todoGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "xtodo_todos",
		Help: "Todo items in xTodo",
	}, []string{"status"})
)

func getTodoDir() string {
	todoDir := os.Getenv("TODO_HOME")
	if todoDir != "" {
		return todoDir
	}

	homeDir, _ := os.UserHomeDir()
	return homeDir + "/workspace/todos"
}

func main() {
	prometheus.MustRegister(todoGauge)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		panic(http.ListenAndServe(":18080", nil))
	}()

	for {
		done, total := getTodoCount()
		fmt.Println(done, total)

		todoGauge.WithLabelValues("done").Set(float64(done))
		todoGauge.WithLabelValues("notdone").Set(float64(total - done))

		time.Sleep(1 * time.Second)
	}
}

func getTodoCount() (int, int) {
	todoDir := getTodoDir()

	fmt.Println("opening todos from :", todoDir)

	now := time.Now().Local()
	month := now.Month().String()[:3]
	year := now.Local().Year()
	day := now.Day()

	file := fmt.Sprintf("%s/data/%d%s/README.md", todoDir, year, month)

	fmt.Println("opening file : ", file)

	f, err := os.Open(file)
	if err != nil {
		fmt.Println("unable to open file", err)
		return 0, 0
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	inside := false
	items := []string{}

	prefix := fmt.Sprintf("## %s %02d", month, day)
	fmt.Println("today: ", prefix)

	for scanner.Scan() {
		line := scanner.Text()

		if inside && (line == "" || line[0] == '#') {
			break
		}

		if inside {
			items = append(items, line)
		} else if strings.HasPrefix(line, prefix) {
			inside = true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	total := len(items)
	done := 0
	for _, item := range items {
		if strings.HasPrefix(item, "- [x]") {
			done++
		}
	}
	return done, total
}
