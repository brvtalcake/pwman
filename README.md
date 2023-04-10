# PWMan

PWMan is a simple password manager in Go using [AES](https://pkg.go.dev/crypto/aes@go1.20.2) and [Cipher](https://pkg.go.dev/crypto/cipher@go1.20.2) Golang packages for encryption (based on a key needed to start the app and wich is not stored anywhere, except in the memory during the execution), and [bzip2](https://github.com/dsnet/compress/tree/master/bzip2) for efficient and lossless compression.

It implements an "in-terminal" interface, built on top of the [tview library](https://github.com/rivo/tview).

This is a(n unifinished) hobby project and I have NO particular security skills, which suggests not to expect too much from it, especially since this is the first project I'm coding in Go.

## TODO

TBD

## Special mention

Both of Decrypt and Encrypt functions are origated from [here](https://bruinsslot.jp/post/golang-crypto/)
