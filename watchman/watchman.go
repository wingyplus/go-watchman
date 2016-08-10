package watchman

import (
	"encoding/json"
	"log"
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

func (c *Client) Trigger(rootpath string, trigger Trigger) {
	err := json.NewEncoder(c.conn).Encode([]interface{}{"trigger", rootpath, trigger})
	if err != nil {
		log.Println(err)
	}
}
