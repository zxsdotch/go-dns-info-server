# DNS Info Server
DNS server which responds to TXT requests with various info:
- a counter
- some random value
- server's unix timestamp

The purpose of this server is primarily for integration testing of DNS
clients -- specifically to ensure that caching is working as expected and/or
that cache busting is working as expected.

The server is running at: `dns-info.zxs.ch`. Give it a try with e.g.:
```bash
dig +short TXT dns-info.zxs.ch

"time=1696792214"
"rand=8724171680981445161"
"counter=42"
```

Code mostly based off this [gist](https://gist.github.com/walm/0d67b4fb2d5daf3edd4fad3e13b162cb).
