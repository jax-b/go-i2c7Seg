package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jax-b/go-i2c7Seg"
)

func main() {
	fmt.Println("Hello world")
	fmt.Println("7 segment I2C Test")
	seg, err := i2c7Seg.NewSevenSegI2C12(0x70, 1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := 0; i < 2; i++ {
		seg.WriteAsciiChar(0, 'D', false)
		seg.WriteAsciiChar(1, 'E', true)
		seg.WriteAsciiChar(3, 'A', false)
		seg.WriteAsciiChar(4, 'D', true)
		seg.DrawColon(false)
		seg.WriteDisplay()
		time.Sleep(time.Second * 2)
		seg.WriteAsciiChar(0, 'B', true)
		seg.WriteAsciiChar(1, 'E', false)
		seg.WriteAsciiChar(3, 'E', true)
		seg.WriteAsciiChar(4, 'F', false)
		seg.DrawColon(true)
		seg.WriteDisplay()
		time.Sleep(time.Second * 2)
	}
	seg.WriteAsciiChar(0, '1', false)
	seg.WriteAsciiChar(1, '2', false)
	seg.WriteAsciiChar(3, '3', false)
	seg.WriteAsciiChar(4, '4', false)
	seg.DrawColon(true, 0xFF)
	seg.WriteDisplay()
	time.Sleep(time.Second * 2)
	seg.Clear()
	seg.WriteDisplay()

	seg.Close()
}
