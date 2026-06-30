package main

import (
	"os"

	"anigate/internal/anigate"
)

func main() {
	os.Exit(anigate.RunCLI(os.Args[1:], anigate.ProductLineMax))
}
