# Tests

This tests the amount of bytes after compressing the communist manifesto taken
from `https://www.gutenberg.org/cache/epub/61/pg61.txt`. And compares said
amount with the result of gzip, and the v1 and the v2 implementation of tcmpr.

## Run locally

> Requires bash, the go compiler tool-chain and awk

```bash
git clone https://github.com/xNaCly/tcmpr
cd tests
./tests.sh
```
