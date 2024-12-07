package main

import (
	"log"
	"net/http"
)

func main() {
	// 静的ファイルの提供
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// サーバーの起動
	log.Println("サーバーを起動します: http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
