package main

import (
	"net"
	"strconv"
	"strings"

	"gopkg.in/mail.v2"
)

type Mailer struct {
	options     Options
	conn        mail.SendCloser
	DefaultFrom string
}

func NewMailer(options Options) (*Mailer, error) {
	host, port, err := net.SplitHostPort(options.Server)
	if err != nil {
		host = options.Server
		port = "587"
	}
	p, err := strconv.Atoi(port)
	if err != nil || p <= 0 || p >= 65535 {
		p = 587
	}

	defaultFrom := options.Username
	if !strings.Contains(defaultFrom, "@") {
		defaultFrom = defaultFrom + "@" + host
	}

	conn, err := mail.NewDialer(host, p, options.Username, options.Password).Dial()
	return &Mailer{
		options:     options,
		conn:        conn,
		DefaultFrom: defaultFrom,
	}, err
}

func (mailer *Mailer) Close() error {
	return mailer.conn.Close()
}

func (mailer *Mailer) Send(m Mail) error {
	m2 := mail.NewMessage()

	m2.SetHeader("To", m.To.String())
	if m.Cc != nil {
		m2.SetHeader("Cc", m.Cc.String())
	}
	if m.Bcc != nil {
		m2.SetHeader("Bcc", m.Bcc.String())
	}
	if m.From != nil {
		m2.SetHeader("From", m.From.String())
	} else {
		m2.SetHeader("From", mailer.DefaultFrom)
	}
	if m.Subject != "" {
		m2.SetHeader("Subject", m.Subject)
	}

	m2.SetBody("text/plain", m.Body)

	for _, a := range m.Attachments {
		m2.Attach(a)
	}

	return mail.Send(mailer.conn, m2)
}
