package main

import (
	"bytes"
	"errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
)

type IStore interface {
	SetData() *User
}

type Neo4j struct{}

func (n Neo4j) SetData() *User {
	url := "http://localhost:7474/db/data/transaction/commit"

	var jsonStr = []byte(`{ "statements": [ { "statement" : "MATCH (n) WHERE n.username = 'Ivanov' RETURN n.username" }]}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("err")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can not parse response")
	}
	value := gjson.Get(string(body), "results.0.data.0.row.0")
	return &User{Username: value.String()}
}

func SetData(s IStore) *User {
	user := s.SetData()
	return user
}

// ##########
// USER MODULE
type User struct {
	Username string
}

type userKey string

func (u *User) GetUser(r *http.Request) (*User, error) {
	ctx := r.Context()
	if ctx == nil {
		return nil, errors.New("No context")
	}

	u, ok := ctx.Value(userKey("user")).(*User)
	if !ok {
		return nil, errors.New("Can not parse user")
	}
	return u, nil
}

// END USER MODULE
// ##########

func main() {
	s := Neo4j{}
	user := SetData(s)
	log.Println(user.Username)
}
