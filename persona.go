// A single function to verify an assertion
package persona

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
)

type Identity struct {
        Email string
        Audience string
        Expires string
        Issuer string
}

func Verify(parameters map[string]interface{}) (*Identity, error) {
        b, err := json.Marshal(parameters)
        if err != nil {
                return nil, err
        }
        resp, err := http.Post("https://verifier.login.persona.org/verify", "application/json", bytes.NewBuffer(b))
        defer resp.Body.Close();
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                return nil, err
        }
        response := make(map[string]string);
        json.Unmarshal(body, response);
        fmt.Println(response["status"])
        return nil, nil
}
