kubectl port-forward svc/postfix 2525:25
nc localhost 2525


EHLO test.local
MAIL FROM: sender@example.com
RCPT TO: test@vmail.today
DATA
Subject: Test Email

This is a test email sent to Postfix.
.
QUIT