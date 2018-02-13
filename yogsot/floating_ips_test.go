package yogsot

import "testing"

func TestFloatingIPS(t *testing.T) {
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
  FloatingIP:
    Region: nyc3
    Type: FloatingIP
    DropletID: MyDroplet`)
	response, err := parseTemplate(template)
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	f := FloatingIP{}
	err = f.buildRequest("TestStack", response.Resources["FloatingIP"])
	if err != nil {
		t.Fatal("expected error to be nil. was: ", err)
	}
	if f.Request.Region != "nyc3" {
		t.Fatalf("region was incorrect: %s", f.Request.Region)
	}
	if f.DropletName != "MyDroplet" {
		t.Fatalf("droplet name was incorrect. expected: MyDroplet. was: %s", f.DropletName)
	}
	if f.Request.DropletID != -1 {
		t.Fatalf("in case name is given, id should be -1. was: %d", f.Request.DropletID)
	}
}

func TestFloatingIPInCaseDropletIDIsSetUseThat(t *testing.T) {
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
  FloatingIP:
    Region: nyc3
    Type: FloatingIP
    DropletID: 987`)
	response, err := parseTemplate(template)
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	f := FloatingIP{}
	err = f.buildRequest("TestStack", response.Resources["FloatingIP"])
	if err != nil {
		t.Fatal("expected error to be nil. was: ", err)
	}
	if f.Request.Region != "nyc3" {
		t.Fatalf("region was incorrect: %s", f.Request.Region)
	}

	if f.Request.DropletID != 987 {
		t.Fatalf("droplet id was incorrect. expected: 987. was: %d", f.Request.DropletID)
	}
}

func TestFloatingParseError(t *testing.T) {
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
  FloatingIP:
    Reg: nyc3
    Type: FloatingIP
    DropletID: MyDroplet`)
	response, err := parseTemplate(template)
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	f := FloatingIP{}
	err = f.buildRequest("TestStack", response.Resources["FloatingIP"])
	if err == nil {
		t.Fatal("expected to fail because of incorrect key value")
	}
}

func TestFloatingIPMissingDropletID(t *testing.T) {
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
  FloatingIP:
    Region: nyc3
    Type: FloatingIP`)
	response, err := parseTemplate(template)
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	f := FloatingIP{}
	err = f.buildRequest("TestStack", response.Resources["FloatingIP"])
	if err == nil {
		t.Fatal("expected to fail because of missing droplet id")
	}
}
