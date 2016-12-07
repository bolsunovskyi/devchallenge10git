GBot
======

Launch via docker:

- docker build -t gbot_image .
- docker run gbot_image gbot --help
- docker run gbot_image gbot -t [secret token] -r [repo name]


Manual install:

- Install golang 1.6
- Setup $GOPATH
- Place project to $GOPATH/src folder
- Run `go get`
- Run `go build`
- Execute gbot