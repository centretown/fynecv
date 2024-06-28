package hass

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const host = "http://melon:8123"
const api = "/api/"

func Get(cmd string) error {

	log.Printf("\n\n%s\n", cmd)
	req, err := http.NewRequest("GET", host+api+cmd, nil)
	if err != nil {
		log.Println(err, "GET")
		return err
	}

	Request(req)
	return err
}

func Post(cmd string, body string) error {
	log.Printf("\n%s\n%s\n", cmd, body)
	buf := bytes.NewBuffer(([]byte)(body))
	req, err := http.NewRequest("POST", host+api+cmd, buf)
	if err != nil {
		log.Println(err, "POST")
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	Request(req)
	return err
}

func Request(req *http.Request) error {
	client := &http.Client{}
	req.Header.Add("Authorization", "Bearer "+Token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err, "PROCESS")
		return err
	}

	buf, err := io.ReadAll(resp.Body)

	if err != nil {
		if err.Error() != "EOF" {
			log.Println(err, "READ")
			return err
		}
	}

	var v any
	err = json.Unmarshal(buf, &v)
	if err != nil {
		log.Println(err, "UNMARSHAL")
		fmt.Println(string(buf))
		return err
	}

	pretty, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println(string(pretty))
	return err

}
