package main

import (
	"fmt"

	"myCrypto"
)

func main() {
	encodedString := myCrypto.Encode([]byte("erer"))
	fmt.Println(encodedString)
}
