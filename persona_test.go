package persona

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestVerify(t *testing.T) {
	//should not authenticate
	parameters := new(Parameters)
	parameters.Assertion = "ABCDEFG1234567"
	parameters.Audience = "http://localhost:8000"
	user, err := Verify(parameters)
	if err == nil {
		t.Error("Invalid assertion should have returned an error")
	}
        if user != nil {
                t.Error("user should be nil")
        }

	http.Handle("/", http.FileServer(http.Dir("./")))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		parameters := new(Parameters)
		json.Unmarshal(body, parameters)
		parameters.Audience = "http://localhost:8000"
		user, err = Verify(parameters)
		if err != nil {
			t.Error(err)
			return
		}
                if user.Audience == "" || user.Email == "" || user.Expires == 0 || user.Issuer == "" {
                        t.Error("Identity struct not populated")
                }
                if user.Reason != "" {
                        t.Error("reason should be empty")
                }
		return
	})

	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
