// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"time"
)

func mama() {
	for {
		fmt.Println("dfljs")
		time.Sleep(1 * time.Second)
	}
}

func nigger() {
	go mama()
}

func main() {
	fmt.Println("Hello, 世界")

	go nigger()

	time.Sleep(5 * time.Minute)
}
