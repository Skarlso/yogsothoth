package yogsot

import (
	"github.com/digitalocean/godo"
)

// Firewall is a struct for firewall creation requests.
type Firewall struct {
	Response             *godo.Response
	Firewall             *godo.Firewall
	Request              *godo.FirewallRequest
	DropletNames         []string
	InboundDropletNames  map[string][]string
	OutboundDropletNames map[string][]string
}

// Rule is an in-outboundrule struct wrapper for ease of use
type Rule struct {
	*godo.InboundRule
	*godo.OutboundRule
	In           bool
	DropletNames []string
}

func (fw *Firewall) buildRequest(stackname string, resource map[string]interface{}) error {
	fw.DropletNames = make([]string, 0)
	fw.InboundDropletNames = make(map[string][]string, 0)
	fw.OutboundDropletNames = make(map[string][]string, 0)
	req := &godo.FirewallRequest{}
	for k, v := range resource {
		if k == "Type" {
			continue
		}
		switch k {
		case "Name":
			req.Name = v.(string)
		case "InboundRules":
			inboundRules := []godo.InboundRule{}
			for inKey, value := range v.(map[interface{}]interface{}) {
				rule := new(Rule)
				rule.generateInbound(value.(map[interface{}]interface{}))
				fw.InboundDropletNames[inKey.(string)] = rule.DropletNames
				inboundRules = append(inboundRules, *rule.InboundRule)
			}
			req.InboundRules = inboundRules
		case "OutboundRules":
			outboundRules := []godo.OutboundRule{}
			for outKey, value := range v.(map[interface{}]interface{}) {
				rule := new(Rule)
				rule.generateOutbound(value.(map[interface{}]interface{}))
				fw.OutboundDropletNames[outKey.(string)] = rule.DropletNames
				outboundRules = append(outboundRules, *rule.OutboundRule)
			}
			req.OutboundRules = outboundRules
		case "Tags":
		case "DropletIDs":
			for _, value := range v.([]interface{}) {
				switch o := value.(type) {
				case string:
					fw.DropletNames = append(fw.DropletNames, o)
				case int:
					req.DropletIDs = append(req.DropletIDs, o)
				}
			}
		}
	}
	fw.Request = req
	return nil
}

func (fw *Firewall) build(yogClient *YogClient) error {
	context := NewContext()
	firewall, response, err := yogClient.Firewalls.Create(context, fw.Request)
	if err != nil {
		return err
	}
	fw.Response = response
	fw.Firewall = firewall
	return nil
}

func (fw *Firewall) setFirewallDropletIDs(ids []int) {
	fw.Request.DropletIDs = append(fw.Request.DropletIDs, ids...)
}

func (fw *Firewall) setInboundDropletIDs(ids []int) {
	// for _, in := range fw.Request.InboundRules {

	// }
}

func (fw *Firewall) setOutboundDropletIDs(ids []int) {
}

func (r *Rule) generateInbound(v map[interface{}]interface{}) {
	r.InboundRule = new(godo.InboundRule)
	for key, innerV := range v {
		if key == "Protocol" {
			r.InboundRule.Protocol = innerV.(string)
		}
		if key == "PortRange" {
			r.InboundRule.PortRange = innerV.(string)
		}
		if key == "Sources" {
			r.getSources(innerV.(map[interface{}]interface{}))
		}
	}
}

func (r *Rule) generateOutbound(v map[interface{}]interface{}) {
	r.OutboundRule = new(godo.OutboundRule)
	for key, innerV := range v {
		if key == "Protocol" {
			r.OutboundRule.Protocol = innerV.(string)
		}
		if key == "PortRange" {
			r.OutboundRule.PortRange = innerV.(string)
		}
		if key == "Destinations" {
			r.getDestinations(innerV.(map[interface{}]interface{}))
		}
	}
}

func (r *Rule) getSources(sources map[interface{}]interface{}) {
	r.DropletNames = make([]string, 0)
	r.InboundRule.Sources = new(godo.Sources)
	r.InboundRule.Sources.Addresses = make([]string, 0)
	r.InboundRule.Sources.Tags = make([]string, 0)
	r.InboundRule.Sources.DropletIDs = make([]int, 0)
	r.InboundRule.Sources.LoadBalancerUIDs = make([]string, 0)
	for k, v := range sources {
		switch k {
		case "Addresses":
			for _, value := range v.([]interface{}) {
				r.InboundRule.Sources.Addresses = append(r.InboundRule.Sources.Addresses, value.(string))
			}
		case "Tags":
			for _, value := range v.([]interface{}) {
				r.InboundRule.Sources.Tags = append(r.InboundRule.Sources.Tags, value.(string))
			}
		case "DropletIDs":
			for _, value := range v.([]interface{}) {
				switch o := value.(type) {
				case string:
					r.DropletNames = append(r.DropletNames, o)
				case int:
					r.InboundRule.Sources.DropletIDs = append(r.InboundRule.Sources.DropletIDs, o)
				}
			}
		case "LoadBalancerUIDs":
			for _, value := range v.([]interface{}) {
				r.InboundRule.Sources.LoadBalancerUIDs = append(r.InboundRule.Sources.LoadBalancerUIDs, value.(string))
			}
		}
	}
}

func (r *Rule) getDestinations(destinations map[interface{}]interface{}) {
	r.DropletNames = make([]string, 0)
	r.OutboundRule.Destinations = new(godo.Destinations)
	r.OutboundRule.Destinations.Addresses = make([]string, 0)
	r.OutboundRule.Destinations.Tags = make([]string, 0)
	r.OutboundRule.Destinations.DropletIDs = make([]int, 0)
	r.OutboundRule.Destinations.LoadBalancerUIDs = make([]string, 0)
	for k, v := range destinations {
		switch k {
		case "Addresses":
			for _, value := range v.([]interface{}) {
				r.OutboundRule.Destinations.Addresses = append(r.OutboundRule.Destinations.Addresses, value.(string))
			}
		case "Tags":
			for _, value := range v.([]interface{}) {
				r.OutboundRule.Destinations.Tags = append(r.OutboundRule.Destinations.Tags, value.(string))
			}
		case "DropletIDs":
			for _, value := range v.([]interface{}) {
				switch o := value.(type) {
				case string:
					r.DropletNames = append(r.DropletNames, o)
				case int:
					r.OutboundRule.Destinations.DropletIDs = append(r.OutboundRule.Destinations.DropletIDs, o)
				}
			}
		case "LoadBalancerUIDs":
			for _, value := range v.([]interface{}) {
				r.OutboundRule.Destinations.LoadBalancerUIDs = append(r.OutboundRule.Destinations.LoadBalancerUIDs, value.(string))
			}
		}
	}
}
