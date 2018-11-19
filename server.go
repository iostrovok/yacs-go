package main

/*

Run YACS processing via the command line.

*/

import (
	"fmt"

	"github.com/iostrovok/yacs-go/httpserver"
)

var rootDir string = ""

func main() {
	fmt.Println("Start 0.0.2...")
	httpserver.Run(rootDir)
}
