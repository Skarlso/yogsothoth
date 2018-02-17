package yogsot

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/digitalocean/godo"
)

var requiredFields = []string{"Name", "Algorithm", "Region", "ForwardingRules", "HealthCheck", "StickySessions", "DropletIDs", "RedirectHttpToHttps"}

// LoadBalancer is a struct which creates a LoadBalancer create request.
type LoadBalancer struct {
	Response     *godo.Response
	LoadBalancer *godo.LoadBalancer
	Request      *godo.LoadBalancerRequest
	DropletNames []string
}

func (lb *LoadBalancer) build(yogClient *YogClient) error {
	context := NewContext()
	loadBalancer, response, err := yogClient.LoadBalancers.Create(context, lb.Request)
	if err != nil {
		return err
	}
	lb.Response = response
	lb.LoadBalancer = loadBalancer
	return nil
}

func (lb *LoadBalancer) buildRequest(stackname string, resource map[string]interface{}) error {
	// err := checkRequiredFields(resource)
	// if err != nil {
	// 	return err
	// }
	req := &godo.LoadBalancerRequest{}

	for k, v := range resource {
		if k == "Type" {
			continue
		}
		if k == "ForwardingRules" {
			forwardingRules := []godo.ForwardingRule{}
			for _, value := range v.(map[interface{}]interface{}) {
				fRule := &godo.ForwardingRule{}
				for key, val := range value.(map[interface{}]interface{}) {
					ref := reflect.ValueOf(fRule)
					refVal := reflect.Indirect(ref).FieldByName(key.(string))
					if refVal == reflect.ValueOf(nil) {
						return errors.New("field not found: " + key.(string))
					}
					refVal.Set(reflect.ValueOf(val))
				}
				forwardingRules = append(forwardingRules, *fRule)
			}
			req.ForwardingRules = forwardingRules
			continue
		}

		if k == "HealthCheck" {
			hck := &godo.HealthCheck{}
			for key, value := range v.(map[interface{}]interface{}) {
				ref := reflect.ValueOf(hck)
				refVal := reflect.Indirect(ref).FieldByName(key.(string))
				if refVal == reflect.ValueOf(nil) {
					return errors.New("field not found: " + key.(string))
				}
				refVal.Set(reflect.ValueOf(value))
			}
			req.HealthCheck = hck
			continue
		}

		if k == "StickySessions" {

			continue
		}

		if k == "DropletIDs" {

			continue
		}

		ref := reflect.ValueOf(req)
		val := reflect.Indirect(ref).FieldByName(k)
		if val == reflect.ValueOf(nil) {
			return errors.New("field not found: " + k)
		}
		val.Set(reflect.ValueOf(v))
	}
	lb.Request = req
	return nil
}

func (lb *LoadBalancer) setDropletIDs(IDs []int) {
	if len(lb.Request.DropletIDs) < 1 {
		lb.Request.DropletIDs = IDs
	}
}

func checkRequiredFields(resource map[string]interface{}) error {
	for _, v := range requiredFields {
		if _, ok := resource[v]; !ok {
			s := fmt.Sprintf("missing required fields: %s", v)
			return errors.New(s)
		}
	}
	return nil
}
