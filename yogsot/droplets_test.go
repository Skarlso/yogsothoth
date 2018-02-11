package yogsot

import "testing"

func TestDropletRequestBuilderFingerprints(t *testing.T) {
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
      - Fingerprint1
      - Fingerprint2`)
	response, err := parseTemplate(template)
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	d := Droplet{}
	err = d.buildRequest("TestStack", response.Resources["Droplet"])
	if err != nil {
		t.Fatal("expected error to be nil. was: ", err)
	}
	if len(d.Request.SSHKeys) < 2 {
		t.Fatalf("fingerprint count was incorrect: %d", len(d.Request.SSHKeys))
	}
	if d.Request.SSHKeys[0].Fingerprint != "Fingerprint1" {
		t.Fatalf("expect: 'Fingerprint1' was: %s", d.Request.SSHKeys[0].Fingerprint)
	}
}

func TestDropletRequestBuilderVolumes(t *testing.T) {
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
    Volumes:
      - VolumeName1
      - VolumeName2`)
	response, err := parseTemplate(template)
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	d := Droplet{}
	err = d.buildRequest("TestStack", response.Resources["Droplet"])
	if err != nil {
		t.Fatal("expected error to be nil. was: ", err)
	}
	if len(d.Request.Volumes) < 2 {
		t.Fatalf("volumes count was incorrect: %d", len(d.Request.Volumes))
	}
	if d.Request.Volumes[0].Name != "VolumeName1" {
		t.Fatalf("expect: 'VolumeName1' was: %s", d.Request.Volumes[0].Name)
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
	d := Droplet{}
	err = d.buildRequest("TestStack", response.Resources["Droplet"])
	if err == nil {
		t.Fatal("expected error to be not nil")
	}
}
