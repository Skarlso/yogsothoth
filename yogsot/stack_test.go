package yogsot

import (
	"testing"
)

func TestCreateStack(t *testing.T) {
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
	request := CreateStackRequest{TemplateBody: template, StackName: "TestStack"}
	yogClient := NewClient("testToken")
	response, err := yogClient.CreateStack(request)
	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
	}
	if len(response.Resources) < 1 {
		t.Fatal("should have contained one created resource")
	}
}
