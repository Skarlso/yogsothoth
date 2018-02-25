package yogsot

import (
	"github.com/digitalocean/godo"
)

// Firewall is a struct for firewall creation requests.
type Firewall struct {
	Response                *godo.Response
	Firewall                *godo.Firewall
	Request                 *godo.FirewallRequest
	DropletNames            []string
	InboundDropletNames     map[string][]string
	OutboundDropletNames    map[string][]string
	InboundRequestsForName  map[string]*godo.InboundRule
	OutboundRequestsForName map[string]*godo.OutboundRule
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
	fw.InboundRequestsForName = make(map[string]*godo.InboundRule, 0)
	fw.OutboundRequestsForName = make(map[string]*godo.OutboundRule, 0)
	req := new(godo.FirewallRequest)
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
				fw.InboundRequestsForName[inKey.(string)] = rule.InboundRule
				inboundRules = append(inboundRules, *rule.InboundRule)
			}
			req.InboundRules = inboundRules
		case "OutboundRules":
			outboundRules := []godo.OutboundRule{}
			for outKey, value := range v.(map[interface{}]interface{}) {
				rule := new(Rule)
				rule.generateOutbound(value.(map[interface{}]interface{}))
				fw.OutboundDropletNames[outKey.(string)] = rule.DropletNames
				fw.OutboundRequestsForName[outKey.(string)] = rule.OutboundRule
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

func (fw *Firewall) setInboundDropletIDs(inbound *godo.InboundRule, ids []int) {
	inbound.Sources.DropletIDs = append(inbound.Sources.DropletIDs, ids...)
}

func (fw *Firewall) setOutboundDropletIDs(outbound *godo.OutboundRule, ids []int) {
	outbound.Destinations.DropletIDs = append(outbound.Destinations.DropletIDs, ids...)
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
	r.InboundRule.Sources.Addresses,
		r.InboundRule.Sources.Tags,
		r.InboundRule.Sources.DropletIDs,
		r.InboundRule.Sources.LoadBalancerUIDs,
		r.DropletNames = generateSources(sources)
}

func (r *Rule) getDestinations(destinations map[interface{}]interface{}) {
	r.DropletNames = make([]string, 0)
	r.OutboundRule.Destinations = new(godo.Destinations)
	r.OutboundRule.Destinations.Addresses,
		r.OutboundRule.Destinations.Tags,
		r.OutboundRule.Destinations.DropletIDs,
		r.OutboundRule.Destinations.LoadBalancerUIDs,
		r.DropletNames = generateSources(destinations)
}

func generateSources(data map[interface{}]interface{}) (
	addresses []string,
	tags []string,
	ids []int,
	uids []string,
	dropletNames []string) {

	for k, v := range data {
		switch k {
		case "Addresses":
			for _, value := range v.([]interface{}) {
				addresses = append(addresses, value.(string))
			}
		case "Tags":
			for _, value := range v.([]interface{}) {
				tags = append(tags, value.(string))
			}
		case "DropletIDs":
			for _, value := range v.([]interface{}) {
				switch o := value.(type) {
				case string:
					dropletNames = append(dropletNames, o)
				case int:
					ids = append(ids, o)
				}
			}
		case "LoadBalancerUIDs":
			for _, value := range v.([]interface{}) {
				uids = append(uids, value.(string))
			}
		}
	}
	return
}
