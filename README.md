mbox
====
[![Build Status](https://travis-ci.org/blabber/mbox.svg?branch=master)](https://travis-ci.org/blabber/mbox)
[![Coverage Status](https://coveralls.io/repos/blabber/mbox/badge.svg)](https://coveralls.io/r/blabber/mbox)
[![GoDoc](https://godoc.org/github.com/blabber/mbox?status.svg)](https://godoc.org/github.com/blabber/mbox)

Package mbox parses the mbox file format into messages and formats messages into mbox files.

Documentation
-------------
Please refer to [GoDoc](https://godoc.org/github.com/blabber/mbox) for the API
documentation.

Alternatives
------------
* github.com/sam-falvo/mbox

Both packages provide ways to parse the mbox file format. After looking at the
code of both packages I decided to roll my own variant, leveraging the standard
library as much as possible. This implementation basically boils down to one
well tested split function for bufio.Scanner.
