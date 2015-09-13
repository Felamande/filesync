package main

import (
	"net/http"

	"github.com/Felamande/filesync/syncer"
)

func main() {
	syncer.New().Run()
	http.ListenAndServe(":8070", nil)
}

