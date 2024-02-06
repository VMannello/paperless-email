# PMail
A [paperless-ngx](https://docs.paperless-ngx.com/post-consumption) post-consumption emailer.

### Features

- Send emails with attachments to specified recipients based on tags.
- Configurable SMTP settings for each email account.
- Support for secure connections with TLS and server certificate verification.

### Releases
Linux (Docker), Darwin, and Windows releases available:  
https://github.com/VMannello/paperless-email/releases

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
      tls: true

  - name: hotmail
    email: account2@hotmail.com
    smtp:
      server: smtp.live.com
      port: 587
      username: account2@hotmail.com
      password: account2_password
      tls: true

# map tags to messages
tags:
  # tags are case-sensitive and must be unique
  # fields default to empty or false

  # all options example
  my-tag:
    # account from above
    from: gmail
    # list of recipients
    to:
      - "recipient-one@email.com"
      - "recipient-two@email.com"
    cc:
      - "carbon-copy@email.com"
    bcc:
      - "blind@email.com"

    include_attachment: true
    # use {{ var_name }} in subject and body for templating
    # available variables found here:
    # https://docs.paperless-ngx.com/advanced_usage/
    subject: "pmail automated delivery - {{ DOCUMENT_ID }}"
    body: |
      Attached file details:
        - {{ DOCUMENT_FILE_NAME }}
        - {{ DOCUMENT_ID }}
        - {{ DOCUMENT_CREATED }} 
        - {{ DOCUMENT_DOWNLOAD_URL }}

  # each tag can have a unique message
  simple-no-attachment:
    from: hotmail
    to:
      - "recipient@email.com"
    subject: "paperless processed a fil {{ DOCUMENT_FILE_NAME }} | {{ DOCUMENT_CREATED }}"
    body: "Processed {{ DOCUMENT_FILE_NAME }} | Download at: {{ DOCUMENT_DOWNLOAD_URL }}"
```

### Set up Paperless Post-Consume In Docker
Using the linux build provided add `pmail` to your container

```yaml
webserver:
  # ...
  volumes:
    # ensure both pmail and the config yaml are the docker volume
    - /home/paperless-ngx/scripts:/path/in/container/scripts/
  environment:
    # ...
    PAPERLESS_POST_CONSUME_SCRIPT: "/path/in/container/scripts/pmail_linux config.yaml"
  # ... more settings ... #
```

### Paperless Environment Variables

This assumes the following variables as defined by paperless-ngx:  
https://docs.paperless-ngx.com/advanced_usage/

Required by the script:
```shell
DOCUMENT_TAGS: <Comma separated list of tags applied (if any)>
DOCUMENT_SOURCE_PATH: < Optional, path to the file to attach to the email>
```

All other variables can be used in subject and body of email.
