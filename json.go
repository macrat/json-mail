package main

import (
	"encoding/json"
	"io"
	"net/mail"
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

type AddressList []*mail.Address

func (l AddressList) MarshalJSON() ([]byte, error) {
	xs := make([]string, len(l))
	for i, x := range l {
		if x.Name != "" {
			xs[i] = x.String()
		} else {
			xs[i] = x.Address
		}
	}

	return json.Marshal(xs)
}

func (l *AddressList) UnmarshalJSON(text []byte) error {
	var xs StringList
	if err := json.Unmarshal(text, &xs); err != nil {
		return err
	}

	for _, x := range xs {
		if ys, err := mail.ParseAddressList(x); err != nil {
			if err.Error() == "mail: no address" {
				continue
			}
			return err
		} else {
			*l = append(*l, ys...)
		}
	}

	return nil
}

type Mail struct {
	To          AddressList   `json:"to,omitempty"`
	Cc          AddressList   `json:"cc,omitempty"`
	Bcc         AddressList   `json:"bcc,omitempty"`
	From        *mail.Address `json:"from,omitempty"`
	Subject     string        `json:"subject,omitempty"`
	Body        string        `json:"body"`
	Attachments StringList    `json:"attachments,omitempty"`
}

type MailList []Mail

func (l MailList) MarshalJSON() ([]byte, error) {
	if l == nil {
		return []byte("[]"), nil
	} else {
		return json.Marshal([]Mail(l))
	}
}

func (l *MailList) UnmarshalJSON(text []byte) error {
	if err := json.Unmarshal(text, (*[]Mail)(l)); err == nil {
		return nil
	}

	var m Mail
	if err := json.Unmarshal(text, &m); err != nil {
		return err
	}

	*l = []Mail{m}
	return nil
}

type MailScanner struct {
	decoder *json.Decoder
	buf     MailList
	err     error
}

func NewMailScanner(r io.Reader) *MailScanner {
	return &MailScanner{
		decoder: json.NewDecoder(r),
	}
}

func (s *MailScanner) Scan() bool {
	s.err = nil

	if len(s.buf) >= 2 {
		s.buf = s.buf[1:]
		return true
	}

	for {
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

func (s *MailScanner) Mail() Mail {
	return s.buf[0]
}
