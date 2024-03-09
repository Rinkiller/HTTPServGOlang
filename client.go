package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	dl, err := net.Dial("tcp", "Localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		text, _, err := bufio.NewReader(os.Stdout).ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		if _, err := dl.Write([]byte(text)); err != nil {
			log.Fatal(err)
		}
		text1, err := bufio.NewReader(dl).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if len(text) > 0 {
			fmt.Println("UpperText: ", string(text1))
		}
	}

}

// 2024/03/07 21:58:31 read /dev/stdout: The handle is invalid.
// exit status 1
