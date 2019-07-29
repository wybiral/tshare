# tshare [![GoDoc](https://godoc.org/github.com/wybiral/tshare?status.svg)](https://godoc.org/github.com/wybiral/tshare)
This package implements (2,3) XOR threshold secret sharing for splitting secrets into three shares. None of the shares alone give away any information about the secret (other than the length) but any combination of two shares is able to fully recover the secret.
