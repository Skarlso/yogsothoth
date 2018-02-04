package yogsot

import (
	"testing"

	"github.com/digitalocean/godo"
)

func TestRequestBuilder(t *testing.T) {
	template := []byte(`
Parameters:
  StackName:
    Description: The name of the stack to deploy
    Type: String
    Default: FurnaceStack
  Port:
    Description: Test port
    Type: Number
    Default: 80

Resources:
  Droplet:
    Name: MyDroplet
    Type: Droplet
    Image:
      Slug: "ubuntu-14-04-x64"`)
	response, err := parseTemplate(template)
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	for _, v := range response.Resources {
		if v["Type"] == "Droplet" {
			d := buildResource("TestStack", "Droplet", v)
			req, err := d.build(v)
			if err != nil {
				t.Fatal("expected error to be nil. was: ", err)
			}
			godoReq := req.(*godo.DropletCreateRequest)
			if godoReq.Name != "MyDroplet" {
				t.Fatalf("droplet name was not MyDroplet. was: %s", godoReq.Name)
			}
		}
	}
}

func TestRequestBuilderUnknownField(t *testing.T) {
	template := []byte(`
Parameters:
  StackName:
    Description: The name of the stack to deploy
    Type: String
    Default: FurnaceStack
  Port:
    Description: Test port
    Type: Number
    Default: 80

Resources:
  Droplet:
    Name: MyDroplet
    Type: Droplet
    Asdf: Bla
    Image:
      Slug: "ubuntu-14-04-x64"`)
	response, err := parseTemplate(template)
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	for k, v := range response.Resources {
		if v["Type"] == "Droplet" {
			d := buildResource("TestStack", "Droplet", v)
			_, err := d.build(v)
			if err == nil && k == "Asdf" {
				t.Fatal("expected error to be not nil")
			}
		}
	}
}
