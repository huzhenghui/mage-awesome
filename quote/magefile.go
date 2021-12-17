//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"daonao.com/quote/lib"
)

func Hello() {
	fmt.Println(lib.Hello())
}

func Test() error {
	cmd := exec.Command("go", "test", "daonao.com/quote/test")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
