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
	version = "HEAD"
	commit  = "unknown"

	options     = Options{}
	showVersion = flag.Bool("version", false, "Show version information and exit")
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
		fmt.Fprintln(out, "Please see -help to more detail.")
	}

	return msgs == nil
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Fprintf(os.Stdout, "json2mail %s (%s)\n", version, commit)
		os.Exit(0)
	}

	options.ParseEnv()
	if !options.Assert(os.Stderr) {
		os.Exit(2)
	}

	l := NewLogger(os.Stdout)
	s := NewMailScanner(os.Stdin)

	m, err := NewMailer(options)
	if err != nil {
		l.Error("failed to connect server: " + err.Error())
		os.Exit(1)
	}

	for {
		if !s.Scan() {
			if s.Err() == nil {
				break
			}
			l.Error("failed to parse: " + s.Err().Error())
			continue
		}

		l.Mail(s.Mail())

		err := m.Send(s.Mail())
		if err != nil {
			l.Error("failed to send: " + err.Error())
		}
	}
}
