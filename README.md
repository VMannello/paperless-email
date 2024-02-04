# PMail
A [paperless-ngx](https://docs.paperless-ngx.com/post-consumption) post-consumption emailer.

### Features

- Send emails with attachments to specified recipients based on tags.
- Configurable SMTP settings for each email account.
- Support for secure connections with TLS and server certificate verification.

### Setup

**Configuration File**

   Create a `config.yaml` file with your email accounts and tag mappings. Here's an example:

   ```yaml
# define accounts by giving them a name
# to be used in the tag mapping
accounts:
  - name: gmail
    email: account1@gmail.com
    smtp:
      server: smtp.gmail.com
      port: 587
      username: your@gmail.com
      # For gmail generate an "app password" and ensure 2FA:
      # https://support.google.com/mail/answer/185833?hl=en
      password: your_gmail_app_password
      secureConnection: true
      verifyServerCertificate: true

  - name: hotmail
    email: account2@hotmail.com
    smtp:
      server: smtp.live.com
      port: 587
      username: account2@hotmail.com
      password: account2_password

# map tags to emails
tag_mapping:
  # any docs with "tag-one" attached will go to "gmail"
  tag-one:
    - gmail
  # any docs with "tag-two" attached will go to "hotmail"
  tag-two:
    - hotmail
  # any docs with "tag-three" attached will go
  # to both "gmail" & "hotmail"
  tag-three:
    - gmail
    - hotmail

message:
  # use {{ var_name }} in subject and body for templating
  # available variables found here:
  # https://docs.paperless-ngx.com/advanced_usage/
  subject: "Paperless Emails - {{ DOCUMENT_ID }}"
  body: |
    Attached files:
      - {{DOCUMENT_CREATED}} - {{ DOCUMENT_FILE_NAME }} - {{ DOCUMENT_ID }}
  include_attachment: true
```

### Run the Script
`pmail config.yaml`

### Paperless Environment Variables

This assumes the following variables as defined by paperless-ngx:  
https://docs.paperless-ngx.com/advanced_usage/

Required by the script:
```shell
DOCUMENT_TAGS: <Comma separated list of tags applied (if any)>
DOCUMENT_SOURCE_PATH: <Path to the file to attach to the email>
```

All other variables can be used in subject and body of email.