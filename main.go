package main

import (
	"os"

	"github.com/neko233-com/linuxsafe/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
