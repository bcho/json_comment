package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bcho/json_comment"
)

func main() {
	b, err := ioutil.ReadAll(json_comment.NewStrippedReader(os.Stdin))
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	fmt.Printf(string(b))
}
