//+build mage

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

var Default = ChooseTarget

func MageHelp() error {
	c := exec.Command(os.Args[0])
	c.Stdout = os.Stdout
	c.Env = append(c.Env, "MAGEFILE_HELP=1")
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}

func MageList() error {
	c := exec.Command(os.Args[0])
	c.Stdout = os.Stdout
	c.Env = append(c.Env, "MAGEFILE_LIST=1")
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}

func invokeList(commandListOutStringChannel chan string) error {
	commandList := exec.Command(os.Args[0])
	commandListReader, commandListWriter, commandListPipeErr := os.Pipe()
	if commandListPipeErr != nil {
		return commandListPipeErr
	}
	commandList.Stdout = commandListWriter
	commandList.Env = append(commandList.Env, "MAGEFILE_LIST=1")
	if commandListRunErr := commandList.Run(); commandListRunErr != nil {
		return commandListRunErr
	}
	commandListWriter.Close()
	commandListStdout, commandListReadErr := ioutil.ReadAll(commandListReader)
	commandListReader.Close()
	if commandListReadErr != nil {
		return commandListReadErr
	}
	commandListOutString := fmt.Sprintf("%s", commandListStdout)
	commandListOutStringChannel <- commandListOutString
	return nil
}

func extractTargetNames(listOut string, targetNameChannel chan string) error {
	defer close(targetNameChannel)
	listOutLines := strings.Split(listOut, "\n")
	var targetNames []string
	targetRegexp, err := regexp.Compile("^ +([^ *]+).*$")
	if err != nil {
		return err
	}
	for _, commandListOutLine := range listOutLines {
		targetMatch := targetRegexp.FindAllStringSubmatch(commandListOutLine, -1)
		if 1 == len(targetMatch) {
			target := targetMatch[0]
			if 2 == len(target) {
				targetName := target[1]
				targetNames = append(targetNames, targetName)

			}
		}
	}
	sort.Strings(targetNames)
	for _, targetName := range targetNames {
		targetNameChannel <- targetName
	}
	return nil
}

func chooseGui(targetNames []string, chooseTarget chan string) error {
	defer close(chooseTarget)
	choose := exec.Command("/usr/local/opt/choose-gui/bin/choose")
	stdin, commandListRunErr := choose.StdinPipe()
	if commandListRunErr != nil {
		return commandListRunErr
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, strings.Join(targetNames, "\n"))
	}()

	choose_out, commandListRunErr := choose.Output()
	if commandListRunErr != nil {
		return commandListRunErr
	}
	choose_target := fmt.Sprintf("%s", choose_out)
	chooseTarget <- choose_target
	return nil
}

func invokeSubCommand(choose_target string) error {
	command_sub := exec.Command(os.Args[0], choose_target)
	command_sub.Stdout = os.Stdout
	err := command_sub.Run()
	if err != nil {
		return err
	}
	return nil
}

func ChooseTarget() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(reflect.TypeOf(err).String())
			fmt.Println(err)
			panic(err)
		}
	}()
	listChannel := make(chan string)
	defer close(listChannel)
	go func() {
		if err := invokeList(listChannel); err != nil {
			panic(err)
		}
	}()
	i, ok := <-listChannel
	if ok == false {
		panic(ok)
	}
	targetNamesChannel := make(chan string)
	go func() {
		if err := extractTargetNames(i, targetNamesChannel); err != nil {
			panic(err)
		}
	}()
	var targetNames []string
	for targetName := range targetNamesChannel {
		targetNames = append(targetNames, targetName)
	}
	chooseTarget := make(chan string)
	go func() {
		if err := chooseGui(targetNames, chooseTarget); err != nil {
			panic(err)
		}
	}()
	choose_target := <-chooseTarget
	fmt.Println(choose_target)
	if 0 <= sort.SearchStrings(targetNames, choose_target) {
		if err := invokeSubCommand(choose_target); err != nil {
			return err
		}
	}
	return nil
}
