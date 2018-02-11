package yogsot

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
)

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *godo.Client

	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = godo.NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

func newTestClient() *YogClient {
	yogClient := YogClient{client}
	return &yogClient
}

func TestCreateStack(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/droplets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"droplets": [{"id":1},{"id":2}]}`)
	})

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
	yogClient := newTestClient()
	response, err := yogClient.CreateStack(request)
	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
	}
	if len(response.Resources) < 1 {
		t.Fatal("should have contained one created resource")
	}
}

func TestCreateStackMultipleResources(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/droplets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"droplets": [{"id":1},{"id":2}]}`)
	})

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
    FloatingIP:
      Name: MyFloatingIP
      Type: FloatingIP
      ID: !Ref Droplet`)
	request := CreateStackRequest{TemplateBody: template, StackName: "TestStack"}
	yogClient := newTestClient()
	response, err := yogClient.CreateStack(request)
	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
	}
	if len(response.Resources) < 1 {
		t.Fatal("should have contained one created resource")
	}
}
