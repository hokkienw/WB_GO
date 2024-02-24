package main

import (
	"io/ioutil"
	"log"
	"level_0/pkg/stanpkg"
)

func main() {
	filePath := "data/model.json"
	
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	stanpkg.PublishMessage(content)
}