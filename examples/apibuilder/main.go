package main

import (
	"os"
	"io/ioutil"
	"github.com/youpenglai/apix/apibuilder"
	"fmt"
)

func main() {
	f, _ := os.Open("test_api.yaml")
	testApi, _ := ioutil.ReadAll(f)
	doc := apibuilder.NewApiDoc()
	if err := doc.Parse(testApi); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(doc)
}
