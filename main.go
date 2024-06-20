package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"go.bug.st/serial"
)

func main() {
	flag.Parse()

	portPath := "/dev/ttyUSB0"
	if flag.NArg() > 0 {
		portPath = flag.Arg(0)
	}

	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(portPath, mode)
	if err != nil {
		log.Fatal(err)
	}

	printTimestamp := true

	for {
		b := make([]byte, 1)
		if _, err = port.Read(b); err != nil {
			log.Fatal(err)
		}

		if b[0] == 0xFF {
			fmt.Println("========== CONSOLE BOOTING ==========")
			printTimestamp = true
		} else if b[0] == 0x00 {
			fmt.Println("========== CONSOLE SHUTDOWN ==========")

			// discard 0xa2
			if _, err = port.Read(b); err != nil {
				log.Fatal(err)
			}

			if b[0] != 0xa2 {
				log.Fatal("invalid shutdown byte")
			}

			printTimestamp = true
		} else {
			if printTimestamp {
				fmt.Printf("[%s] ", time.Now().Format("03:04:05.000"))
				printTimestamp = false
			}

			if b[0] == 0x0a {
				printTimestamp = true
			} else if b[0] == 0x0d {
				continue
			}

			os.Stdout.Write(b)
		}
	}
}
