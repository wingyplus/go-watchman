package watchman

import (
	"encoding/json"
	"log"
	"net"
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
	conn, err := net.Dial("unix", "/usr/local/var/run/watchman/wingyplus-state/sock")
	if err != nil {
		panic(err)
	}
	return &Client{
		conn: conn,
	}
}

func (c *Client) Trigger(rootpath string, trigger Trigger) {
	err := json.NewEncoder(c.conn).Encode([]interface{}{"trigger", rootpath, trigger})
	if err != nil {
		log.Println(err)
	}
}
