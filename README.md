mbox
====
[![Build Status](https://travis-ci.org/blabber/mbox.svg?branch=master)](https://travis-ci.org/blabber/mbox)
[![GoDoc](https://godoc.org/github.com/blabber/mbox?status.svg)](https://godoc.org/github.com/blabber/mbox)

Package mbox parses the mbox file format into messages and formats messages into
mbox files.

TODO
----
Curently there is no clean way to delete messages.

Documentation
-------------
Please refer to [GoDoc](https://godoc.org/github.com/blabber/mbox) for the API
documentation.

Notes
-----
This package does not provide advanced mbox manipulation functionalities, like
deleting, searching, etc.

If you get a `bufio.Scanner: token too long` error, you need to increase the
buffer size of the scanner. To double the buffer size you'd do something like
this:

    s := mbox.NewScanner(...)
    
    bufsize := 2*bufio.MaxScanTokenSize
    buf := make([]byte, bufsize)
    s.Buffer(buf, bufsize)

Maintainers Note
----------------
This project is in maintenance mode. I am not actively working on this package
any more, but I am responding to issues and merge the occasional pull request.
If you miss any functionality, feel free to send pull requests.

If you are interested in taking over this package, ping me.

Alternatives
------------
* github.com/sam-falvo/mbox
* ~~github.com/bthomson/mbox~~

Both packages provide ways to parse the mbox file format. After looking at the
code of both packages I decided to roll my own variant, leveraging the standard
library as much as possible. This implementation basically boils down to one
well tested split function for bufio.Scanner.
