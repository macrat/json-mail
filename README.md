JSON To Mail
============

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

$ cat ./mails.csv | json2mail -server smtp.example.com -username you@example.com -password your_p@ssword
EOS
```

```
$ export JSON2MAIL_SERVER="smtp.example.com" JSON2MAIL_USERNAME="you@example.com" JSON2MAIL_PASSWORD="your_p@ssword"

$ cat ./mails.csv | json2mail
EOS
```
