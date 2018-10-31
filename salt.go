package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Auth struct{
	User string  `json:"username"`
	Passwd string  `json:"password"`
	Eauth string `json:"eauth"`
}

type Cfg struct {
	Base string	`json:"base"`
	Auth Auth
}

type Client struct {
	token string
	cfg Cfg
}

// Minion
type Minion struct {
	ID            string   `json:"id"`
	Name          string   `json:"nodename"`
	Host          string   `json:"host"`
	Domain        string   `json:"domain"`
	OS            string   `json:"os"`
	OSRelease     string   `json:"osrelease"`
	OSName        string   `json:"osfullname"`
	Kernel        string   `json:"kernel"`
	KernelRelease string   `json:"kernelrelease"`
	Shell         string   `json:"shell"`
	ARCH          string   `json:"osarch"`
	CPUS          int      `json:"num_cpus"`
	RAM           int      `json:"mem_total"`
	CPUModel      string   `json:"cpu_model"`
	CPUFlags      []string `json:"cpu_flags"`
	Virtual       string   `json:"virtual"`
	IPv4          []string `json:"ipv4"`
	IPv6          []string `json:"ipv6"`
	Path          string   `json:"path"`
	ServerID      int      `json:"server_id"`
}

// MinionsResponse
type MinionsResponse struct {
	  Minions []map[string]Minion `json:"return"`
}

// JobsResponse
type JobsResponse struct {
	Jobs []map[string]Job `json:"return"`
}

// JobResponse
type JobResponse struct {
	Job []Job `json:"info"`
}


// ExecutionResponse
type ExecutionResponse struct {
	Job []Job `json:"return"`
}

type Result struct {
	PID     int    `json:"pid"`
	Retcode int    `json:"retcode"`
	Return	interface{} `json:"return"`
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
}

// Job
type Job struct {
	ID         string            `json:"jid"`
	Function   string            `json:"Function"`
	Target     string            `json:"Target"`
	User       string            `json:"User"`
	StartTime  string            `json:"StartTime"`
	TargetType string            `json:"Target-Type"`
	Arguments  []string            `json:"Arguments"`
	Minions    []string            `json:"Minions"`
	Result     map[string]Result `json:"Result"`
}

// Running
func (j *Job) Running() bool {
	if len(j.Minions) != len(j.Result) {
		return false
	}
	return true
}

// Successful
func (j *Job) Successful() bool {
	for _, r := range j.Result {
		if r.Retcode != 0 {
			return false
		}
	}
	return true
}

// Get Token
func (salt *Client) Auth() (error) {
	urls := salt.cfg.Base + "/login"
	data := fmt.Sprintf(`{ "username":"%s", "password":"%s", "eauth": "%s" }`, salt.cfg.Auth.User, salt.cfg.Auth.Passwd,salt.cfg.Auth.Eauth)
	req, err := http.NewRequest("POST", urls, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return  err
	}

	req.Header.Set("X-Auth-Token", salt.token)
	req.Header.Set("Content-Type", "application/json")

	r := &http.Client{}
	resp,err := r.Do(req)
	if err != nil {
		return  err
	}
	if resp.StatusCode != 200{
		return  errors.New("请求失败")
	}
	if err != nil {
		return err
	}
	salt.token = resp.Header.Get("X-Auth-Token")
	return nil
}

// New Client
func New(cfg *Cfg) (*Client, error) {
	client := Client{}
	client.cfg = *cfg
	if err := client.Auth(); err != nil {
		return nil, err
	}
	return &client, nil
}

// Get
func (salt *Client) Get(postfix string) (*http.Response, error) {
	urls := salt.cfg.Base + postfix

	req, err := http.NewRequest("GET",urls, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", salt.token)
	req.Header.Set("Content-Type", "application/json")

	r := &http.Client{}
	resp,err := r.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("请求失败")
	}
	return resp,nil
}

// POST
func (salt *Client) Post(postfix string, data []byte) (*http.Response, error) {
	urls := salt.cfg.Base + postfix
	req, err := http.NewRequest("POST", urls, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", salt.token)
	req.Header.Set("Content-Type", "application/json")

	r := &http.Client{}
	resp,err := r.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 202 {
		return nil, errors.New("请求失败")
	}
	return resp,nil
}

// Close
func ParseResponse(resp *http.Response) (*[]byte, error) {
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	return &data, err
}

// Minions
func (salt *Client) Minions() (map[string]Minion, error) {
	m := MinionsResponse{}
	resp, err := salt.Get("/minions")
	if err != nil {
		return m.Minions[0], err
	}

	data, err := ParseResponse(resp)
	if err != nil {
		return m.Minions[0], err
	}

	err = json.Unmarshal(*data, &m)
	return m.Minions[0], err
}

// Minion
func (salt *Client) Minion(id string) (Minion, error) {
	var m Minion
	uri := fmt.Sprintf("/minions/%s", id)
	resp, err := salt.Get(uri)
	if err != nil {
		return m, err
	}

	data, err := ParseResponse(resp)
	fmt.Println(string(*data))
	if err != nil {
		return m, err
	}

	err = json.Unmarshal(*data, &m)

	return m, err
}

// Jobs
func (salt *Client) Jobs() ([]map[string]Job, error) {
	jr := JobsResponse{}

	resp, err := salt.Get("/jobs")
	if err != nil {
		return jr.Jobs, err
	}

	data, err := ParseResponse(resp)
	if err != nil {
		return jr.Jobs, err
	}

	err = json.Unmarshal(*data, &jr)

	return jr.Jobs, err
}

// Job
func (salt *Client) Job(id string) (Job, error) {
	j := JobResponse{}

	uri := fmt.Sprintf("/jobs/%s", id)
	resp, err := salt.Get(uri)
	if err != nil {
		return Job{}, err
	}

	data, err := ParseResponse(resp)
	if err != nil {
		return Job{}, err
	}
	err = json.Unmarshal(*data, &j)
	return j.Job[0], err
}

// Execute
func (salt *Client) Execute(client, function, command, target, targetType string) (string, error) {
	er := ExecutionResponse{}
	var req string
	if function == "test.ping" {
		req = fmt.Sprintf(`{"client":"%s", "fun": "%s", "tgt": "%s", "expr_form": "%s"}`, client,function, target, targetType)
	}else {
		req = fmt.Sprintf(`{"client":"%s", "fun": "%s", "arg": "%s", "tgt": "%s", "expr_form": "%s"}`, client,function, command, target, targetType)
	}
	resp, err := salt.Post("/minions", []byte(req))
	if err != nil {
		return "", err
	}

	data, err := ParseResponse(resp)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(*data, &er)
	return er.Job[0].ID, err
}