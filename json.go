package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/mail"
	"os"
	"strings"
)

type StringList []string

func (l *StringList) UnmarshalJSON(text []byte) error {
	if err := json.Unmarshal(text, (*[]string)(l)); err == nil {
		return nil
	}

	var s string
	if err := json.Unmarshal(text, &s); err != nil {
		return err
	}

	*l = []string{s}
	return nil
}

type Address mail.Address

func (a *Address) String() string {
	if a.Name != "" {
		return (*mail.Address)(a).String()
	} else {
		return a.Address
	}
}

func (a *Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *Address) UnmarshalJSON(text []byte) error {
	var s string
	if err := json.Unmarshal(text, &s); err != nil {
		return err
	}

	addr, err := mail.ParseAddress(s)
	if err != nil {
		return err
	}
	*a = Address(*addr)
	return nil
}

type AddressList []*Address

func (l *AddressList) UnmarshalJSON(text []byte) error {
	var xs StringList
	if err := json.Unmarshal(text, &xs); err != nil {
		return err
	}

	for _, x := range xs {
		addrs, err := mail.ParseAddressList(x)
		if err != nil {
			if err.Error() == "mail: no address" {
				continue
			}
			return err
		}

		for _, a := range addrs {
			*l = append(*l, (*Address)(a))
		}
	}

	return nil
}

func (l *AddressList) String() string {
	xs := make([]string, len(*l))
	for i, x := range *l {
		xs[i] = x.String()
	}
	return strings.Join(xs, ", ")
}

type Mail struct {
	To          AddressList `json:"to,omitempty"`
	Cc          AddressList `json:"cc,omitempty"`
	Bcc         AddressList `json:"bcc,omitempty"`
	From        *Address    `json:"from,omitempty"`
	Subject     string      `json:"subject,omitempty"`
	Body        string      `json:"body,omitempty"`
	Attachments StringList  `json:"attachments,omitempty"`
}

func (m Mail) Validate() error {
	if len(m.To) == 0 {
		return errors.New("field `to` is required")
	}
	for _, a := range m.Attachments {
		if _, err := os.Stat(a); os.IsNotExist(err) {
			return errors.New("attachment not found: " + a)
		}
	}
	return nil
}

type MailList []Mail

func (l MailList) Validate() error {
	for _, m := range l {
		if err := m.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (l MailList) MarshalJSON() ([]byte, error) {
	if l == nil {
		return []byte("[]"), nil
	} else {
		return json.Marshal([]Mail(l))
	}
}

func (l *MailList) UnmarshalJSON(text []byte) error {
	if err := json.Unmarshal(text, (*[]Mail)(l)); err == nil {
		return l.Validate()
	}

	var m Mail
	if err := json.Unmarshal(text, &m); err != nil {
		return err
	}

	*l = []Mail{m}
	return l.Validate()
}

type MailScanner struct {
	decoder    *json.Decoder
	readerCopy *bytes.Buffer
	buf        MailList
	err        error
}

func NewMailScanner(r io.Reader) *MailScanner {
	copied := bytes.NewBuffer([]byte{})

	return &MailScanner{
		readerCopy: copied,
		decoder:    json.NewDecoder(io.TeeReader(r, copied)),
	}
}

func (s *MailScanner) Scan() bool {
	s.err = nil

	if len(s.buf) >= 2 {
		s.buf = s.buf[1:]
		return true
	}

	for {
		s.readerCopy.Reset()
		s.err = s.decoder.Decode(&s.buf)
		if s.err == io.EOF {
			s.err = nil
			return false
		} else if s.err != nil {
			return false
		}

		if len(s.buf) > 0 {
			break
		}
	}

	return true
}

func (s *MailScanner) Err() error {
	return s.err
}

func (s *MailScanner) CurrentString() string {
	return s.readerCopy.String()
}

func (s *MailScanner) Mail() Mail {
	return s.buf[0]
}
