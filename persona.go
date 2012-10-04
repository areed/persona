// A single function to verify an assertion
package persona

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "time"
)

type Parameters struct {
  Assertion string  `json:"assertion"`
  Audience string `json:"audience"`
}

type IdentityOpinionated struct {
        Email string
//        Audience url.Url
        Expires time.Time
        Issuer string
}

type Identity struct {
        Email string
        Audience string
        Expires string
        Issuer string
}

func Verify(parameters *Parameters) (*Identity, error) {
        b, err := json.Marshal(parameters)
fmt.Println(string(b))
        if err != nil {
                return nil, err
        }
        resp, err := http.Post("https://verifier.login.persona.org/verify", "application/json", bytes.NewBuffer(b))
        defer resp.Body.Close();
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                return nil, err
        }
        fmt.Println(string(body))
        response := make(map[string]string);
        json.Unmarshal(body, response);
        return nil, nil
}
