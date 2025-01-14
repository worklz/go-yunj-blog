package main

import (
	"github.com/worklz/yunj-blog-go/app/es"
	_ "github.com/worklz/yunj-blog-go/pkg/boot"
)

func main() {
	err := es.Reset()
	if err != nil {
		panic(err)
	}
}
