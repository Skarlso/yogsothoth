package yogsot

import (
	"testing"
)

func TestDropletResource(t *testing.T) {
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
      Slug: "ubuntu-14-04-x64"
    SSHKeys:
      - Fingerprint: "Fingerprint1"
      - Fingerprint: "Fingerprint2"`)
	response, err := parseTemplate(template)
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	for _, v := range response.Resources {
		req := getDropletRequest("StackName", v)
		if req.Name != "MyDroplet" {
			t.Errorf("'Name' did not equal MyDroplet. Was instead: %s", req.Name)
		}
		if req.Image.Slug != "ubuntu-14-04-x64" {
			t.Errorf("'Slug' did not equal ubuntu-14-04-x64. Was instead: %s", req.Image.Slug)
		}
		if req.SSHKeys[0].Fingerprint != "Fingerprint1" {
			t.Errorf("'SSHKey 0' did not equal Fingerprint1. Was instead: %s", req.SSHKeys[0].Fingerprint)
		}
		if req.SSHKeys[1].Fingerprint != "Fingerprint2" {
			t.Errorf("'SSHKey 1' did not equal Fingerprint2. Was instead: %s", req.SSHKeys[1].Fingerprint)
		}
	}
}
