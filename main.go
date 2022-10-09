package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	todoGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "taskwarrior_todos",
		Help: "Todo items in taskwarrior",
	}, []string{"status"})
)

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
		todoGauge.WithLabelValues("pending").Set(float64(total - done))

		time.Sleep(10 * time.Second)
	}
}

func getTodoCount() (int, int) {
	var done, pending int

	cmd := exec.Command("task", "status:completed", "count")
	out1, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return 0, 0
	}
	if _, err := fmt.Sscanf(string(out1), "%d", &done); err != nil {
		fmt.Println(err.Error())
		return 0, 0
	}

	cmd = exec.Command("task", "status:pending", "count")
	out2, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return 0, 0
	}
	if _, err := fmt.Sscanf(string(out2), "%d", &pending); err != nil {
		fmt.Println(err.Error())
		return 0, 0
	}

	return done, pending + done
}
