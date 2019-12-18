package services

import (
	"log"
	"testing"
)

func TestReadYAML(t *testing.T) {
	err := ReadYAML("/Users/naive/sysmgmt-next/update.yaml")
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
