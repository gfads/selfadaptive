package main

import (
	"fmt"
	"selfadaptive/shared"
)

func Behaviour() {

	fmt.Println("Sent Message:", shared.EncryptMessage(shared.Keys32[0], shared.PlainText)+" ")

}
