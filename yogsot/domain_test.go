package yogsot

import (
	"io/ioutil"
	"testing"
)

func TestDomainCreate(t *testing.T) {
	template, err := ioutil.ReadFile("./fixtures/domain_test_TestDomainCreate.yaml")
	if err != nil {
		t.Fatal("unexpected error while opening fixture: ", err)
	}
	response, err := parseTemplate(template)
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	d := Domain{}
	err = d.buildRequest("TestStack", response.Resources["SkarlsoDomain"])
	if err != nil {
		t.Fatal("expected error to be nil. was: ", err)
	}
	if d.Request.Name != "skarlso.io" {
		t.Fatalf("name was incorrect: %s", d.Request.Name)
	}
	if d.Request.IPAddress != "127.0.0.1" {
		t.Fatalf("ip was incorrect: %s", d.Request.IPAddress)
	}

}
