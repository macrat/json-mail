package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	server   = flag.String("server", "", "SMTP server address")
	username = flag.String("username", "", "Username for login SMTP server")
	password = flag.String("password", "", "Password for login SMTP server")
)

func parseEnvironmentVariable() {
	if *server == "" {
		*server = os.Getenv("JSON_MAIL_SERVER")
	}
	if *username == "" {
		*username = os.Getenv("JSON_MAIL_USERNAME")
	}
	if *password == "" {
		*password = os.Getenv("JSON_MAIL_PASSWORD")
	}
}

func checkFlagsIfOK(out io.Writer) (ok bool) {
	var msgs []string
	if *server == "" {
		msgs = append(msgs, "--server is required.")
	}
	if *username == "" {
		msgs = append(msgs, "--username is required.")
	}
	if *password == "" {
		msgs = append(msgs, "--password is required.")
	}

	if msgs != nil {
		fmt.Fprintln(os.Stderr, "error:")
		for _, m := range msgs {
			fmt.Fprintln(os.Stderr, " ", m)
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Please see --help to more detail.")
	}

	return msgs == nil
}

func main() {
	flag.Parse()
	parseEnvironmentVariable()
	if !checkFlagsIfOK(os.Stderr) {
		os.Exit(2)
	}

	fmt.Println("server:", *server)
	fmt.Println("username:", *username)
	fmt.Println("password:", *password)

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
