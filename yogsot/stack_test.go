package yogsot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
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
		expected := map[string]interface{}{
			"name":               "name",
			"region":             "region",
			"size":               "size",
			"image":              "ubuntu-14-04-x64",
			"ssh_keys":           nil,
			"backups":            false,
			"ipv6":               false,
			"private_networking": false,
			"monitoring":         false,
			"tags":               []interface{}{"TestStack"},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{"droplet":{"id":1}, "links":{"actions": [{"id": 1, "href": "http://example.com", "rel": "create"}]}}`)
	})

	template, err := ioutil.ReadFile("./fixtures/stack_test_TestCreateStack.yaml")
	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
	}
	request := CreateStackRequest{TemplateBody: template, StackName: "TestStack"}
	yogClient := newTestClient()
	response, ye := yogClient.CreateStack(request)
	if ye.Error != nil {
		t.Fatal("unexpected error: " + ye.Error.Error())
	}
	if len(response.Resources) < 1 {
		t.Fatal("should have contained one created resource")
	}
}

func TestCreateStackMoreThanFiveDroplets(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/droplets", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"name":               "name",
			"region":             "region",
			"size":               "size",
			"image":              "ubuntu-14-04-x64",
			"ssh_keys":           nil,
			"backups":            false,
			"ipv6":               false,
			"private_networking": false,
			"monitoring":         false,
			"tags":               []interface{}{"TestStack"},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{"droplet":{"id":1}, "links":{"actions": [{"id": 1, "href": "http://example.com", "rel": "create"}]}}`)
	})

	// Normally, the none uniqueness of these names should raise an error.
	// But for unit testing purposes, I'm ignore those for now.
	template, err := ioutil.ReadFile("./fixtures/stack_test_TestCreateStackMoreThanFiveDroplets.yaml")
	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
	}
	request := CreateStackRequest{TemplateBody: template, StackName: "TestStack"}
	yogClient := newTestClient()
	response, ye := yogClient.CreateStack(request)
	if ye.Error != nil {
		t.Fatal("unexpected error: " + ye.Error.Error())
	}
	if len(response.Resources) < 1 {
		t.Fatal("should have contained one created resource")
	}
}

func TestCreateStackMultipleResources(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/droplets", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"name":               "MyDroplet",
			"region":             "region",
			"size":               "size",
			"image":              "ubuntu-14-04-x64",
			"ssh_keys":           nil,
			"backups":            false,
			"ipv6":               false,
			"private_networking": false,
			"monitoring":         false,
			"tags":               []interface{}{"TestStack"},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{"droplet":{"id":987}, "links":{"actions": [{"id": 1, "href": "http://example.com", "rel": "create"}]}}`)
	})

	mux.HandleFunc("/v2/floating_ips", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"region":     "nyc3",
			"droplet_id": float64(987),
		}
		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body = %+v, expected = %+v", v, expected)
		}

		fmt.Fprint(w, `{"floating_ip":{"region":{"slug":"nyc3"},"droplet":{"id":987},"ip":"192.168.0.1"}}`)
	})

	template, err := ioutil.ReadFile("./fixtures/stack_test_TestCreateStackMultipleResources.yaml")
	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
	}
	request := CreateStackRequest{TemplateBody: template, StackName: "TestStack"}
	yogClient := newTestClient()
	response, ye := yogClient.CreateStack(request)
	if ye.Error != nil {
		t.Fatal("unexpected error: " + ye.Error.Error())
	}
	if len(response.Resources) < 1 {
		t.Fatal("should have contained one created resource")
	}
	for _, v := range response.Resources {
		if f, ok := v.(*FloatingIP); ok {
			if f.Request.DropletID != 987 {
				t.Fatalf("floatingip request droplet id should have equaled 987. Was instead: %d\n", f.Request.DropletID)
			}
		}
	}
}

var lbCreateJSONResponse = `
{
    "load_balancer":{
        "id":"8268a81c-fcf5-423e-a337-bbfe95817f23",
        "name":"example-lb-01",
        "ip":"",
        "algorithm":"round_robin",
        "status":"new",
        "created_at":"2016-12-15T14:19:09Z",
        "forwarding_rules":[
            {
                "entry_protocol":"https",
                "entry_port":443,
                "target_protocol":"http",
                "target_port":80,
                "certificate_id":"a-b-c"
            },
            {
                "entry_protocol":"https",
                "entry_port":444,
                "target_protocol":"https",
                "target_port":443,
                "tls_passthrough":true
            }
        ],
        "health_check":{
            "protocol":"http",
            "port":80,
            "path":"/index.html",
            "check_interval_seconds":10,
            "response_timeout_seconds":5,
            "healthy_threshold":5,
            "unhealthy_threshold":3
        },
        "sticky_sessions":{
            "type":"cookies",
            "cookie_name":"DO-LB",
            "cookie_ttl_seconds":5
        },
        "region":{
            "name":"New York 1",
            "slug":"nyc1",
            "sizes":[
                "512mb",
                "1gb",
                "2gb",
                "4gb",
                "8gb",
                "16gb"
            ],
            "features":[
                "private_networking",
                "backups",
                "ipv6",
                "metadata",
                "storage"
            ],
            "available":true
        },
        "droplet_ids":[
            2,
            21
        ],
        "redirect_http_to_https":true
    }
}
`

func TestCreateStackLoadBalancer(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/load_balancers", func(w http.ResponseWriter, r *http.Request) {
		v := new(godo.LoadBalancerRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, http.MethodPost)
		// assert.Equal(t, createRequest, v)

		fmt.Fprint(w, lbCreateJSONResponse)
	})

	mux.HandleFunc("/v2/droplets", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"name":               "MyDroplet",
			"region":             "region",
			"size":               "size",
			"image":              "ubuntu-14-04-x64",
			"ssh_keys":           nil,
			"backups":            false,
			"ipv6":               false,
			"private_networking": false,
			"monitoring":         false,
			"tags":               []interface{}{"TestStack"},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{"droplet":{"id":987}, "links":{"actions": [{"id": 1, "href": "http://example.com", "rel": "create"}]}}`)
	})

	template, err := ioutil.ReadFile("./fixtures/stack_test_TestCreateStackLoadBalancer.yaml")
	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
	}
	request := CreateStackRequest{TemplateBody: template, StackName: "TestStack"}
	yogClient := newTestClient()
	response, ye := yogClient.CreateStack(request)
	if ye.Error != nil {
		t.Fatal("unexpected error: " + ye.Error.Error())
	}
	if len(response.Resources) < 1 {
		t.Fatal("should have contained one created resource")
	}
	for _, v := range response.Resources {
		if ldb, ok := v.(*LoadBalancer); ok {
			if !reflect.DeepEqual(ldb.Request.DropletIDs, []int{12, 987}) {
				t.Fatalf("Droplet ids should have equaled [12, 987]. Was instead: %v", ldb.Request.DropletIDs)
			}
		}
	}
}

func TestCreateStackWithDomain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/domains", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"name":       "skarlso.io",
			"ip_address": "127.0.0.1",
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprint(w, `{"domain":{"name":"skarlso.io"}}`)
	})

	template, err := ioutil.ReadFile("./fixtures/stack_test_TestCreateStackWithDomain.yaml")
	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
	}
	request := CreateStackRequest{TemplateBody: template, StackName: "TestStack"}
	yogClient := newTestClient()
	response, ye := yogClient.CreateStack(request)
	if ye.Error != nil {
		t.Fatal("unexpected error: " + ye.Error.Error())
	}
	if len(response.Resources) < 1 {
		t.Fatal("should have contained one created resource")
	}
	for _, v := range response.Resources {
		if d, ok := v.(*Domain); ok {
			if d.Request.IPAddress != "127.0.0.1" {
				t.Fatalf("ip did not equal 127.0.0.1. was: %s", d.Request.IPAddress)
			}
			if d.Request.Name != "skarlso.io" {
				t.Fatal("name did not equal skarlso.io. was: ", d.Request.Name)
			}
		}
	}
}

func TestCreateStackNoType(t *testing.T) {
	setup()
	defer teardown()

	template, err := ioutil.ReadFile("./fixtures/stack_test_TestCreateStackNoType.yaml")
	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
	}
	request := CreateStackRequest{TemplateBody: template, StackName: "TestStack"}
	yogClient := newTestClient()
	_, ye := yogClient.CreateStack(request)
	if ye.Error == nil {
		t.Fatal("should have failed with no Type for fields.")
	}
}

var firewallCreateJSONBody = `
{
  "name": "f-i-r-e-w-a-l-l",
  "inbound_rules": [
    {
      "protocol": "icmp",
      "sources": {
        "addresses": ["0.0.0.0/0"],
        "tags": ["frontend"],
        "droplet_ids": [123, 456],
        "load_balancer_uids": ["lb-uid"]
      }
    },
    {
      "protocol": "tcp",
      "ports": "8000-9000",
      "sources": {
        "addresses": ["0.0.0.0/0"]
      }
    }
  ],
  "outbound_rules": [
    {
      "protocol": "icmp",
      "destinations": {
        "tags": ["frontend"]
      }
    },
    {
      "protocol": "tcp",
      "ports": "8000-9000",
      "destinations": {
        "addresses": ["::/1"]
      }
    }
  ],
  "droplet_ids": [123],
  "tags": ["frontend"]
}
`

var firewallJSONResponse = `
{
  "firewall": {
    "id": "fe6b88f2-b42b-4bf7-bbd3-5ae20208f0b0",
    "name": "f-i-r-e-w-a-l-l",
    "status": "waiting",
    "inbound_rules": [
      {
        "protocol": "icmp",
        "ports": "0",
        "sources": {
          "tags": ["frontend"]
        }
      },
      {
        "protocol": "tcp",
        "ports": "8000-9000",
        "sources": {
          "addresses": ["0.0.0.0/0"]
        }
      }
    ],
    "outbound_rules": [
      {
        "protocol": "icmp",
        "ports": "0"
      },
      {
        "protocol": "tcp",
        "ports": "8000-9000",
        "destinations": {
          "addresses": ["::/1"]
        }
      }
    ],
    "created_at": "2017-04-06T13:07:27Z",
    "droplet_ids": [
      123
    ],
    "tags": [
      "frontend"
    ],
    "pending_changes": [
      {
        "droplet_id": 123,
        "removing": false,
        "status": "waiting"
      }
    ]
  }
}
`

func TestCreateFirewall(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/droplets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		fmt.Fprintf(w, `{"droplet":{"id":951, "name":"MyDroplet"}, "links":{"actions": [{"id": 1, "href": "http://example.com", "rel": "create"}]}}`)
	})

	expectedFirewallRequest := &godo.FirewallRequest{
		Name: "f-i-r-e-w-a-l-l",
		InboundRules: []godo.InboundRule{
			{
				Protocol: "icmp",
				Sources: &godo.Sources{
					Addresses:        []string{"0.0.0.0/0"},
					Tags:             []string{"frontend"},
					DropletIDs:       []int{123, 456},
					LoadBalancerUIDs: []string{"lb-uid"},
				},
			},
			{
				Protocol:  "tcp",
				PortRange: "8000-9000",
				Sources: &godo.Sources{
					Addresses: []string{"0.0.0.0/0"},
				},
			},
		},
		OutboundRules: []godo.OutboundRule{
			{
				Protocol: "icmp",
				Destinations: &godo.Destinations{
					Tags: []string{"frontend"},
				},
			},
			{
				Protocol:  "tcp",
				PortRange: "8000-9000",
				Destinations: &godo.Destinations{
					Addresses: []string{"::/1"},
				},
			},
		},
		DropletIDs: []int{123},
		Tags:       []string{"frontend"},
	}

	mux.HandleFunc("/v2/firewalls", func(w http.ResponseWriter, r *http.Request) {
		v := new(godo.FirewallRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, http.MethodPost)

		var actualFirewallRequest *godo.FirewallRequest
		json.Unmarshal([]byte(firewallCreateJSONBody), &actualFirewallRequest)
		if !reflect.DeepEqual(actualFirewallRequest, expectedFirewallRequest) {
			t.Errorf("Request body = %+v, expected %+v", actualFirewallRequest, expectedFirewallRequest)
		}

		fmt.Fprint(w, firewallJSONResponse)
	})

	template, err := ioutil.ReadFile("./fixtures/stack_test_TestCreateFirewall.yaml")
	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
	}
	request := CreateStackRequest{TemplateBody: template, StackName: "TestStack"}
	yogClient := newTestClient()
	response, ye := yogClient.CreateStack(request)
	if ye.Error != nil {
		t.Fatal("unexpected error: " + ye.Error.Error())
	}
	if len(response.Resources) < 1 {
		t.Fatal("should have contained one created resource")
	}
	var checker *Firewall
	for _, f := range response.Resources {
		if fw, ok := f.(*Firewall); ok {
			checker = fw
			break
		}
	}

	expectedFWDropletIDs := []int{123, 951}
	if !reflect.DeepEqual(checker.Request.DropletIDs, expectedFWDropletIDs) {
		t.Fatalf("expected '%+v' did not equal actual: '%+v'", expectedFWDropletIDs, checker.Request.DropletIDs)
	}
	expectedInboundIDs := []int{123, 456, 951}
	for _, v := range checker.Request.InboundRules {
		if len(v.Sources.DropletIDs) > 2 {
			if !reflect.DeepEqual(v.Sources.DropletIDs, expectedInboundIDs) {
				t.Fatalf("expected '%+v' did not equal actual: '%+v'", expectedInboundIDs, v.Sources.DropletIDs)
			}
		}
	}

}
