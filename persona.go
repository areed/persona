// A single function to verify an assertion
package persona

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//Parameters is the type passed to the Verify function
type Parameters struct {
	Assertion string `json:"assertion"`
	Audience  string `json:"audience"`
}

//expiryTime is used as an embedded struct in Identity and inherits all the methods of time.Time
//except UnmarshalJSON
type expiryTime struct {
	time.Time
}

//UnmarshalJSON takes the milliseconds since 1/1/1970 and converts it into type time.Time
func (e *expiryTime) UnmarshalJSON(data []byte) (err error) {
	milliseconds, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	e.Time = time.Unix(milliseconds/1000, 0)
	return
}

//Identity is the type returned to the application if authentication succeeds
type Identity struct {
	Email    string
	Audience string
	Expires  *expiryTime
	Issuer   string
}

//failure is the type the response unmarshals into first to check for unsuccessful authentication
type failure struct {
	Reason string
}

//Verify sends the assertion to Persona for verifications
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
	f := new(failure)
	json.Unmarshal(body, f)
	if f.Reason != "" {
		return nil, errors.New(f.Reason)
	}
	i := new(Identity)
	json.Unmarshal(body, i)
	return i, nil
}

//a convenience function that allows assertion and audience to be passed as strings instead of as
//fields of Parameter
func VerifyArgs(assertion, audience string) (*Identity, error) {
	p := new(Parameters)
	p.Assertion = assertion
	p.Audience = audience
	return Verify(p)
}
