package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type Options struct {
	Server   string
	Username string
	Password string
}

var (
	options = Options{}
)

func init() {
	flag.StringVar(&options.Server, "server", "", "SMTP server address")
	flag.StringVar(&options.Username, "username", "", "Username for login SMTP server")
	flag.StringVar(&options.Password, "password", "", "Password for login SMTP server")
}

func (opts *Options) ParseEnv() {
	if opts.Server == "" {
		opts.Server = os.Getenv("JSON2MAIL_SERVER")
	}
	if opts.Username == "" {
		opts.Username = os.Getenv("JSON2MAIL_USERNAME")
	}
	if opts.Password == "" {
		opts.Password = os.Getenv("JSON2MAIL_PASSWORD")
	}
}

func (opts *Options) Assert(out io.Writer) (ok bool) {
	var msgs []string
	if opts.Server == "" {
		msgs = append(msgs, "--server is required.")
	}
	if opts.Username == "" {
		msgs = append(msgs, "--username is required.")
	}
	if opts.Password == "" {
		msgs = append(msgs, "--password is required.")
	}

	if msgs != nil {
		fmt.Fprintln(out, "error:")
		for _, m := range msgs {
			fmt.Fprintln(out, " ", m)
		}
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Please see --help to more detail.")
	}

	return msgs == nil
}

func main() {
	flag.Parse()
	options.ParseEnv()
	if !options.Assert(os.Stderr) {
		os.Exit(2)
	}

	m, err := NewMailer(options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	s := NewMailScanner(os.Stdin)

	for {
		if !s.Scan() {
			if s.Err() == nil {
				break
			}
			fmt.Fprintln(os.Stderr, "failed to parse:", s.Err())
			continue
		}

		err := m.Send(s.Mail())
		if err != nil {
			fmt.Println(err)
		}
	}
}
