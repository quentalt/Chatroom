package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pusher "github.com/pusher/pusher-http-go/v5"
)

var client = pusher.Client{
	AppID:   "1334984",
	Key:     "5b8e12f5ab64caf7dac4",
	Secret:  "ca6ca6f91ddcf5371b57",
	Cluster: "eu",
	Secure:  true,
}

type user struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

func registerNewUser(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var newUser user

	err = json.Unmarshal(body, &newUser)
	if err != nil {
		panic(err)
	}
	client.Trigger("update", "new-user", newUser)
	json.NewEncoder(rw).Encode(newUser)
}

func pusherAuth(res http.ResponseWriter, req *http.Request) {
	params, _ := ioutil.ReadAll(req.Body)
	response, err := client.AuthenticatePrivateChannel(params)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(res, string(response))

}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/new/user", registerNewUser)
	http.HandleFunc("/pusher/auth", pusherAuth)

	log.Fatal(http.ListenAndServe(":8090", nil))
}
