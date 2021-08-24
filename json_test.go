package main_test

import (
	"strings"
	"reflect"
	"testing"
	"encoding/json"

	"github.com/macrat/json-mail"
)

func TestStringList(t *testing.T) {
	tests := []struct{
		Input  string
		Output string
	} {
		{`""`, `[""]`},
		{`[]`, `[]`},
		{`"hello world"`, `["hello world"]`},
		{`"foo,bar"`, `["foo,bar"]`},
		{`["hello world"]`, `["hello world"]`},
		{`["hello","world"]`, `["hello","world"]`},
	}

	for _, tt := range tests {
		t.Run(tt.Input, func(t *testing.T) {
			var l main.StringList

			if err := json.Unmarshal([]byte(tt.Input), &l); err != nil {
				t.Fatalf("failed to unmarshal: %s", err)
			}

			j, err := json.Marshal(l)
			if err != nil {
				t.Fatalf("failed to marshal: %s", err)
			}

			if string(j) != tt.Output {
				t.Errorf("unexpected result:\nexpected: %s\n but got: %s", tt.Output, j)
			}
		})
	}
}

func TestAddressList(t *testing.T) {
	tests := []struct{
		Input string
		Output string
	} {
		{`""`, `[]`},
		{`["", ""]`, `[]`},
		{`[]`, `[]`},
		{`null`, `[]`},
		{`"test@example.com"`, `["test@example.com"]`},
		{`"hello <world@example.com>"`, `["\"hello\" \u003cworld@example.com\u003e"]`},
		{`" foo@example.com, hello <world@example.com>"`, `["foo@example.com","\"hello\" \u003cworld@example.com\u003e"]`},
		{`["foo@example.com, hello <world@example.com>"]`, `["foo@example.com","\"hello\" \u003cworld@example.com\u003e"]`},
		{`["bar@example.com", "foo@example.com, hello <world@example.com>"]`, `["bar@example.com","foo@example.com","\"hello\" \u003cworld@example.com\u003e"]`},
	}

	for _, tt := range tests {
		t.Run(tt.Input, func(t *testing.T) {
			var l main.AddressList

			if err := json.Unmarshal([]byte(tt.Input), &l); err != nil {
				t.Fatalf("failed to unmarshal: %s", err)
			}

			j, err := json.Marshal(l)
			if err != nil {
				t.Fatalf("failed to marshal: %s", err)
			}

			if string(j) != tt.Output {
				t.Errorf("unexpected result:\nexpected: %s\n but got: %s", tt.Output, j)
			}
		})
	}
}

func TestMailList(t *testing.T) {
	tests := []struct{
		Input string
		Output string
	} {
		{`[]`, `[]`},
		{`null`, `[]`},
		{`{}`, `[{"body":""}]`},
		{`[{}]`, `[{"body":""}]`},
		{`[{}, {}]`, `[{"body":""},{"body":""}]`},
		{`[{"subject":"hello"}]`, `[{"subject":"hello","body":""}]`},
		{
			`[{"to": "a@example.com,  b@example.com"},           {"to":["c@example.com","d@example.com"]}]`,
			`[{"to":["a@example.com","b@example.com"],"body":""},{"to":["c@example.com","d@example.com"],"body":""}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Input, func(t *testing.T) {
			var l main.MailList

			if err := json.Unmarshal([]byte(tt.Input), &l); err != nil {
				t.Fatalf("failed to unmarshal: %s", err)
			}

			j, err := json.Marshal(l)
			if err != nil {
				t.Fatalf("failed to marshal: %s", err)
			}

			if string(j) != tt.Output {
				t.Errorf("unexpected result:\nexpected: %s\n but got: %s", tt.Output, j)
			}
		})
	}
}

func TestMailScanner(t *testing.T) {
	tests := []struct{
		Input string
		Bodies []string
	} {
		{``, []string{}},
		{`{"body":"hello"}`, []string{"hello"}},
		{`{"body":"hello"} {"body": "world"}`, []string{"hello", "world"}},
		{`[{"body":"hello"},{"body":"world"}]`, []string{"hello", "world"}},
	}

	for _, tt := range tests {
		t.Run(tt.Input, func(t *testing.T) {
			s := main.NewMailScanner(strings.NewReader(tt.Input))

			count := 0
			xs := []string{}
			for s.Scan() {
				xs = append(xs, s.Mail().Body)
				count++
			}

			if s.Err() != nil {
				t.Fatalf("failed to scan: %s", s.Err())
			}

			if !reflect.DeepEqual(xs, tt.Bodies) {
				t.Errorf("unexpected bodies:\nexpected: %#v\n but got: %#v", tt.Bodies, xs)
			}
		})
	}
}
