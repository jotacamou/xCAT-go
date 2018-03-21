package xCAT

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client - contains xCAT API connection settings
type Client struct {
	Master   string
	Token    string
	Insecure bool
}

// NewRequest makes a custom request to the xCAT API given a URI string.
func (c *Client) NewRequest(uri string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.Insecure},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", c.requestString(uri), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", c.Token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// GetNodes returns a list of all registered nodes within the xCAT inventory.
func (c *Client) GetAllNodes() []byte {
	body, err := c.NewRequest("/nodes")
	if err != nil {
		fmt.Println(err)
	}
	return body
}

// NodeRange returns a list of nodes in a given xCAT noderange
func (c *Client) NodeRange(nodeRange string, args ...string) (body []byte, err error) {
	switch len(args) {
	case 0:
		nodeRange = fmt.Sprintf("/nodes/%s", nodeRange)
		body, err = c.NewRequest(nodeRange)
		result := make(map[string]interface{})
		//var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, err
		}
		nodes := []string{}
		body, err = json.MarshalIndent(nodes, "", "   ")
		return
	default:
		uri := strings.Join(args, "/")
		nodeRange = fmt.Sprintf("/nodes/%s/%s?pretty=1", nodeRange, uri)
	}
	body, err = c.NewRequest(nodeRange)
	return
}

// GetNetworks returns a list of all registered networks within the xCAT inventory.
func (c *Client) GetNetworks() ([]byte, error) {
	return c.NewRequest("/networks?pretty=1")
}

// GetNetworkObjects returns the network objects requested on a CSV value.
func (c *Client) GetNetworkObjects() ([]byte, error) {
	var networkNames interface{}
	networkNames, err := c.GetNetworks()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(networkNames.([]byte), &networkNames)
	if err != nil {
		return nil, err
	}
	var csv string
	for _, n := range networkNames.([]interface{}) {
		if len(csv) == 0 {
			csv = n.(string)
			continue
		}
		csv = csv + "," + n.(string)
	}
	return c.NewRequest(fmt.Sprintf("/networks/%s?pretty=1", csv))
}

func (c *Client) requestString(uri string) string {
	return fmt.Sprintf("%s%s", c.Master, uri)
}
