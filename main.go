package main

import (
	"math/rand"
	"time"

	"gitlab.com/loranna/fakegps/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute()
}
