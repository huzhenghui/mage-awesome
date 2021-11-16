// +build mage

package main

import (
	"fmt"
)

func init() {
	fmt.Println("init函数先于main函数自动执行")
}

func init() {
	fmt.Println("每个包中可以有多个init函数，每个包中的源文件中也可以有多个init函数")
}

func init() {
	fmt.Println("init函数没有输入参数、返回值，也未声明，所以无法引用")
}

var (
	a = c + b
	b = f()
	c = f()
	d = 3
)

func f() int {
	d++
	return d
}

func init() {
	fmt.Println("a = c + b // == ", a)
	fmt.Println("b = f()   // == ", b)
	fmt.Println("c = f()   // == ", c)
	fmt.Println("d = 3     // == ", d)
}

func init() {
	fmt.Println("如果当前包下有多个init函数，首先按照源文件名的字典序从前往后执行。")
}

func init() {
	fmt.Println("若一个文件中出现多个init函数，则按照出现顺序从前往后执行。")
}
