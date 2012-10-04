// A single function to verify an assertion
package persona

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

//	"time"
)

type Parameters struct {
	Assertion string `json:"assertion"`
	Audience  string `json:"audience"`
}

type Identity struct {
	Reason  string
	Email    string
	Audience string
	Expires  int64
	Issuer   string
}

func Verify(parameters *Parameters) (*Identity, error) {
	b, err := json.Marshal(parameters)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post("https://verifier.login.persona.org/verify", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	i := new(Identity)
	json.Unmarshal(body, i)
	if i.Reason != "" {
		return nil, errors.New(i.Reason)
	}
	return i, nil
}
