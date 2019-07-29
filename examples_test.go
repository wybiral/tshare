package tshare_test

import (
	"fmt"
	"log"

	"github.com/wybiral/tshare"
)

var shares [][]byte

// Split secret into three tagged shares.
func ExampleSplitBytes() {
	secret := []byte("Hello world!")
	shares, err := tshare.SplitBytes(secret)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(shares[0])
	fmt.Println(shares[1])
	fmt.Println(shares[2])
}

// Recover secret by joining two tagged shares.
func ExampleJoinBytes() {
	// any combination of two different shares will work
	secret, err := tshare.JoinBytes(shares[0], shares[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(secret)
}
