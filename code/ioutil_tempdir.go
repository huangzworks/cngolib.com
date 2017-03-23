package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir) // clean up

	fmt.Println(dir)
}
