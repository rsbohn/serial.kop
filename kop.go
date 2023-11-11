package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"go.bug.st/serial"
)

// GOLANG:CLI
// Simple Serial Console
// the serial reader loop:
func readLoop(port serial.Port) {
	buf := make([]byte, 128)

	for {
		n, err := port.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		fmt.Printf("%v", string(buf[:n]))
	}
}

// the input handler:
func main() {
	fmt.Println("Today is a new day.")
	mode := &serial.Mode{BaudRate: 115200}
	port, err := serial.Open("COM15", mode)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	go readLoop(port)

	reader := bufio.NewReader(os.Stdin)
	helpText := ".? for help\n.d to send Ctrl-D\n.q to quit"
	// the input/send loop
	for {
		text, _ := reader.ReadString('\n')
		if strings.HasPrefix(text, ".q") {
			break
		}
		if strings.HasPrefix(text, ".d") {
			text = "\x04"
		}
		if strings.HasPrefix(text, ".?") {
			fmt.Println(helpText)
			continue
		}
		_, err = port.Write([]byte(text + "\n"))
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("bye!")
	os.Exit(0)
}
