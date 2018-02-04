package yogsot

import (
	"testing"
)

func TestDropletRequestBuilder(t *testing.T) {
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
			d := buildResource("Droplet")
			req, err := d.(Droplet).buildRequest("TestStack", v)
			if err != nil {
				t.Fatal("expected error to be nil. was: ", err)
			}
			if req.Name != "MyDroplet" {
				t.Fatalf("droplet name was not MyDroplet. was: %s", req.Name)
			}
		}
	}
}

func TestDropletRequestBuilderUnknownField(t *testing.T) {
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
			d := buildResource("Droplet")
			_, err := d.(Droplet).buildRequest("TestStack", v)
			if err == nil && k == "Asdf" {
				t.Fatal("expected error to be not nil")
			}
		}
	}
}
