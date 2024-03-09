package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	lis, err := net.Listen("tcp4", "Localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is runing")
	con, err := lis.Accept()
	if err != nil {
		log.Fatal(err)
	}
	for {
		line, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("line: ", string(line))
		upperline := strings.ToUpper(string(line))
		if _, err := con.Write([]byte(upperline)); err != nil {
			log.Fatal(err)
		}

	}

}

// 024/03/07 21:58:31 read tcp4 127.0.0.1:8081->127.0.0.1:52327: wsarecv: An existing connection was forcibly closed by the remote host.
// exit status 1
