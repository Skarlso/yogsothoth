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
			d, _ := buildResource(DROPLET)
			err = d.(*Droplet).buildRequest("TestStack", v)
			if err != nil {
				t.Fatal("expected error to be nil. was: ", err)
			}
			if d.(*Droplet).Request.Name != "MyDroplet" {
				t.Fatalf("droplet name was not MyDroplet. was: %s", d.(*Droplet).Request.Name)
			}
		}
	}
}

func TestUnknownResourceType(t *testing.T) {
	_, err := buildResource(999)
	if err == nil {
		t.Fatal("should have failed with unknown resource type")
	}
}
