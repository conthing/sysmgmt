package services

import (
	"log"
	"testing"
)

func TestReadYAML(t *testing.T) {
	err := ReadYAML()
	log.Println(MyItemCollection)
	if err != nil {

		t.Error(err)
	}
}

func TestUnzip(t *testing.T) {
	err := UnZip()
	if err != nil {
		log.Fatal(err)
	}
}
