package main

import (
	"fmt"
	"gee/gee"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.Default()
	r.GET("/", func(ctx *gee.Context) {
		panic(11111)
	})
	r.Run(":8080")
}
