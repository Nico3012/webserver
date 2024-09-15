package printhelloworld

import (
	"fmt"

	"github.com/Nico3012/webserver/gethelloworld"
)

func PrintHelloWorld() {
	fmt.Println(gethelloworld.GetHelloWorld())
}
