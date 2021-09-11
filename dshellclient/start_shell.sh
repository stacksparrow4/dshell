export TERM=xterm
if which python3; then
    python3 -c 'import pty;pty.spawn("/bin/bash")'
elif which python2; then
    python2 -c 'import pty;pty.spawn("/bin/bash")'
elif which python; then
    python -c 'import pty;pty.spawn("/bin/bash")'
elif which script; then
    script -q /dev/null /bin/bash
else
    echo "WARNING: no pty upgrade could be found"
    /bin/bash
fi
