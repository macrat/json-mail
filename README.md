JSON Mail
=========

Send many emails from JSON data.

## Usage

``` shell
$ cat ./mails.csv
[{
  "to": "destination@example.com",
  "from": "you@example.com",
  "subject": "test mail",
  "body": "hello!\nthis is a test"
}]

$ cat ./mails.csv | json-mail -server smtp.example.com:465 -username you@example.com -password your_p@ssword
EOS
```

```
$ export JSON_MAIL_SERVER="smtp.example.com:465" JSON_MAIL_USERNAME="you@example.com" JSON_MAIL_PASSWORD="your_p@ssword"

$ cat ./mails.csv | json-mail
EOS
```
