package persona

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

const host, port string = "http://localhost", ":8000"

var audience string = host + port

func TestVerify(t *testing.T) {
	//should not authenticate
	parameters := new(Parameters)
	parameters.Assertion = "ABCDEFG1234567"
	parameters.Audience = audience
	user, err := Verify(parameters)
	if err == nil {
		t.Error("Invalid assertion should have returned an error")
	}
	if user != nil {
		t.Error("user should be nil")
	}
}

func TestVerifyArgs(t *testing.T) {
	http.Handle("/", http.FileServer(http.Dir("./")))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		var v interface{}
		json.Unmarshal(body, &v)
		user, err := VerifyArgs(v.(map[string]interface{})["assertion"].(string), audience)
		if err != nil {
			t.Error(err)
			return
		}
		if user.Audience == "" || user.Email == "" || user.Issuer == "" {
			t.Error("Identity struct string fields not populated")
		}
		//test will break for a couple minutes on New Year's Eve every year
		if time.Now().Year() != user.Expires.Year() {
			t.Error("The expiry year was ", user.Expires.Year())
		}
		return
	})

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
