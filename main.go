package main

import (
	"encoding/json"
	"os"
	"fmt"
)

func main() {
	s := NewMailScanner(os.Stdin)

	for s.Scan() {
		x, err := json.MarshalIndent(s.Mail(), "", "  ")
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Println(string(x))
		fmt.Println()
	}
	fmt.Println(s.Err())
}
