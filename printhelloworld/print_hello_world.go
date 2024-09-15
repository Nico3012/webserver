package printhelloworld

import (
	"fmt"

	"webserver/gethelloworld"
)

func PrintHelloWorld() {
	fmt.Println(gethelloworld.GetHelloWorld())
}
