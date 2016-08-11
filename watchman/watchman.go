package watchman

import (
	"encoding/json"
	"net"
	"os/exec"
)

type Trigger struct {
	Name       string   `json:"name"`
	Expression []string `json:"expression"`
	Command    []string `json:"command"`
}

type Client struct {
	conn net.Conn
}

func (c *Client) Close() {
	c.conn.Close()
}

func Connect() *Client {
	sockname, err := getSockname()
	if err != nil {
		panic(err)
	}
	conn, err := net.Dial("unix", sockname)
	if err != nil {
		panic(err)
	}
	return &Client{
		conn: conn,
	}
}

func getSockname() (string, error) {
	var sock struct {
		Version  string `json:"version"`
		Sockname string `json:"sockname"`
	}
	out, err := exec.Command("watchman", "get-sockname").Output()
	if err != nil {
		return "", err
	}
	json.Unmarshal(out, &sock)
	return sock.Sockname, nil
}

type TriggerResult struct {
	Version     string `json:"version"`
	TriggerID   string `json:"triggerid"`
	Disposition string `json:"disposition"`
}

func (c *Client) Trigger(rootpath string, trigger Trigger) (*TriggerResult, error) {
	err := json.NewEncoder(c.conn).Encode([]interface{}{"trigger", rootpath, trigger})
	if err != nil {
		return nil, err
	}
	var result TriggerResult
	err = json.NewDecoder(c.conn).Decode(&result)
	return &result, err
}
