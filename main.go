package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Task struct {
	Project string `json:"project"`
	Status string `json:"status"`
}

var (
	todoGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "taskwarrior_todos",
		Help: "Todo items in taskwarrior",
	}, []string{"project", "status"})
)

func main() {
	prometheus.MustRegister(todoGauge)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		panic(http.ListenAndServe(":18080", nil))
	}()

	for {
		todos, err := getTasks()
		if err != nil {
			fmt.Println("unable to get tasks", err)
			time.Sleep(10 * time.Second)
			continue
		}
		fmt.Printf("got todos : %d\n", len(todos))

		countMap := make(map[Task]int)
		for _, todo := range todos {
			countMap[todo] += 1
		}

		for todo, count := range countMap {
			todoGauge.WithLabelValues(todo.Project, todo.Status).Set(float64(count))
		}

		time.Sleep(10 * time.Second)
	}
}

func getTasks() ([]Task, error) {
	cmd := exec.Command("task", "export")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	tasks := []Task{}
	err = json.Unmarshal(out, &tasks)
	return tasks, err
}
