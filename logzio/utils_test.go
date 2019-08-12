package logzio

import (
	"fmt"
	"io/ioutil"
	"log"
)

func ReadFixtureFromFile(fileName string) string {
	content, err := ioutil.ReadFile("testdata/fixtures/"+fileName)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s", content)
}

func ReadResourceFromFile(resourceName string, fileName string) string {
	return fmt.Sprintf(ReadFixtureFromFile(fileName), resourceName)
}
