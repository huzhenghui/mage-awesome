// +build mage

package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/guptarohit/asciigraph"
)

var Default = Plot

func Plot() {
	data := []float64{3, 4, 9, 6, 2, 4, 5, 8, 5, 10, 2, 7, 2, 5, 6}
	graph := asciigraph.Plot(data)

	fmt.Println(graph)
}

func YaegiRun() {
	cmd := exec.Command("yaegi", "run", "-unrestricted", "yaegi.go")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)
	fmt.Println(buf.String())
}
