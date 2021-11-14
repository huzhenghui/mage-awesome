// +build mage

package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"yaegi-import-local/mypackage"
)

var Default = HelloWorld

func HelloWorld() {
	fmt.Println(mypackage.Hello("World"))
}

func YaegiRun() {
	cmd := exec.Command("yaegi", "run", "yaegi.go")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)
	fmt.Println(buf.String())
}
