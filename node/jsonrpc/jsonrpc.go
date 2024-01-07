package jsonrpc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
	ID      int         `json:"id"`
}

type Client struct {
	endpoint   string
	httpClient *http.Client
}

func NewClient(endpoint string) (*Client, error) {
	return &Client{
		endpoint:   endpoint,
		httpClient: &http.Client{},
	}, nil
}

func (c *Client) Call(method string, params any) (*Response, error) {
	req := Request{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(c.endpoint, "application/json", bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcResp Response
	err = json.Unmarshal(respBytes, &rpcResp)
	if err != nil {
		return nil, err
	}

	return &rpcResp, nil
}

func (r *Response) Bind(result any) error {
	resBytes, err := json.Marshal(r.Result)
	if err != nil {
		return err
	}

	return json.Unmarshal(resBytes, result)
}
