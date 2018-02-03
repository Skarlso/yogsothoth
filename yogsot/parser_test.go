package yogsot

import (
	"io/ioutil"
	"testing"
)

func TestParser(t *testing.T) {
	content, err := ioutil.ReadFile("template.yaml")
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	response, err := parseTemplate(content)

	for _, v := range response.Resources {
		if v["Type"] == "Droplet" {
			d := Droplet{}

		}
	}
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	if len(response.Parameters) == 0 {
		t.Fatal("parameters is empty")
	}
	if len(response.Resources) == 0 {
		t.Fatal("resources is empty")
	}
}
