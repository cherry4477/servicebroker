package api

import (
	"encoding/json"
	"fmt"
)

const dcos_Exchange_Token_Api = "/acs/api/v1/auth/login"

type Auth interface {
	Exchange(host string) (string, error)
}

func (c *authConfig) Exchange(host string) (string, error) {
	body := []byte("{\"token\": \"")
	body = append(body, []byte(string(*c))...)
	body = append(body, []byte("\"}")...)

	url := host + dcos_Exchange_Token_Api
	b, err := httpPost(url, ContentType_Json, body)
	if err != nil {
		return "", err
	}

	tk := struct {
		Token string
	}{}

	if err := json.Unmarshal(b, &tk); err != nil {
		return "", err

	}

	if len(tk.Token) == 0 {
		return "", fmt.Errorf("exchange nil token")
	}

	return tk.Token, nil
}
