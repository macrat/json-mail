JSON To Mail
============

Send emails from JSON data.

## Get started

1. Download latest binary from [download page](https://github.com/macrat/json2mail/releases).

2. Write emails as JSON.

   ```json
   [{
     "to": "\"someone\" <destination@example.com>",
     "from": "\"your name\" <you@example.com>",
     "subject": "test mail",
     "body": "hello!\nthis is a test"
   }]
   ```

3. Send emails using `json2email` command

   ```shell
   $ cat emails.json | json2mail -server smtp.example.com -username you@example.com -password your_p@ssword
   ```


## Options

json2mail supports following options.

Some of them can be set by the environment variable.
If set by both ways, the command-line options will be used.

| Command-line               | Environment variable | Description                                                                         |
|----------------------------|----------------------|-------------------------------------------------------------------------------------|
| **-server** *[ADDRESS]*    | JSON2MAIL_SERVER     | SMTP server address.<br />e.g. smtp.example.com                                     |
| **-username** *[USERNAME]* | JSON2MAIL_USERNAME   | Your username to login to the SMTP server.                                          |
| **-password** *[PASSWORD]* | JSON2MAIL_PASSWORD   | Your password to login to the SMTP server.                                          |
| **-source** *[FILE]*       |                      | Path of JSON file that including email data.<br />("-" or "" means read from stdin) |
| **-interval** *[DURATION]* |                      | Interval to send each emails.<br />e.g. `100ms`, `1s`, or `1.5s`                    |
| **-allow-insecure**        |                      | Allow to connect without encryption. (NOT recommended)*                             |
| **-dry-run**               |                      | Run json2mail without server connection. It's convinient for testing JSON source.   |


## Source format

You can write emails as JSON data.

The JSON data should be an object, an array of objects, or objects delimited by new-lines.
Each object, means each emails, can include the following fields.

| Field           | Type                       | Description                                      |
|-----------------|----------------------------|--------------------------------------------------|
| **to**          | String or Array-of-Strings | To address. This is the only required field.     |
| **cc**          | String or Array-of-Strings | CC address.                                      |
| **bcc**         | String or Array-of-Strings | BCC address.                                     |
| **from**        | String                     | From address.                                    |
| **subject**     | String                     | Subject of the email.                            |
| **body**        | String                     | Body of the email. (HTML body is not supported.) |
| **attachments** | String or Array-of-String  | Attachment file paths.                           |

### Examples

A simplest email looks like this.

```json
{
  "to": "destination@example.com",
  "from": "you@example.com",
  "subject": "test mail",
  "body": "hello!\nthis is a test"
}
```

You can send multiple emails as array.

```json
[
  {
    "to": "\"someone\" <destination@example.com>",
    "from": "\"your name\" <you@example.com>",
    "subject": "test mail",
    "body": "hello!\nthis is a test"
  },
  {
    "to": ["alice@example.com", "bob@example.com"],
    "cc": ["charlie@example.com"],
    "subject": "hello",
    "body": "Hi!\n\nIt's a test!"
  }
]
```

Or, you can just write multiple emails in a file without array.

```json
{
  "to": "\"someone\" <destination@example.com>",
  "from": "\"your name\" <you@example.com>",
  "subject": "test mail",
  "body": "hello!\nthis is a test"
}
{
  "to": ["alice@example.com", "bob@example.com"],
  "cc": ["charlie@example.com"],
  "subject": "hello",
  "body": "Hi!\n\nIt's a test!"
}
```


## Hints

- In Windows, multi-byte emails can be broken if using stdin. So please try to write to a file in UTF-8, and then use `-source` option.
