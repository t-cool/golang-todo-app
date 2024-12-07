package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

var tasks []Task

func main() {
	// 初期データの追加
	tasks = append(tasks, Task{
		ID:        1,
		Title:     "Go言語を学ぶ",
		Done:      false,
		CreatedAt: time.Now(),
	})

	// APIハンドラーの設定
	http.HandleFunc("/api/tasks", handleTasks)

	// 静的ファイルの提供
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Println("サーバーを起動します: http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(tasks)
	case "POST":
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		task.ID = len(tasks) + 1
		task.CreatedAt = time.Now()
		tasks = append(tasks, task)
		json.NewEncoder(w).Encode(task)
	}
}
