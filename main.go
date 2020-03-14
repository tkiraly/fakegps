package main

import (
	"math/rand"
	"time"

	"github.com/tkiraly/fakegps/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute()
}
