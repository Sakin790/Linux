package main

import (
	"fmt"
	"os"
)

func main() {
	fs, _ := os.Create("web-server-configuration")
	defer fs.Close()
	fs.WriteString("Hello world sakin")

	_, err := os.Stat("web-server-configuration")
	if os.IsNotExist(err) {
		fmt.Println("ফাইলটি নেই")
	}

}
