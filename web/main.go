package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle(
		"/",
		http.FileServer(http.Dir("./static")),
	)
	fmt.Println(
		"MovieMood web запущен: http://localhost:8080",
	)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска:", err)
	}
}
