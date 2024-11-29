package helper

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func GetVNGCloudToken(clientID, clientSecret string) string {
	url := "https://iam.api.vngcloud.vn/accounts-api/v2/auth/token"
	method := "POST"

	payload := strings.NewReader(`{
	  "grant_type": "client_credentials"
  }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	CheckError(err)
	req.Header.Add("Content-Type", "application/json")
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Add("Authorization", "Basic "+auth)
	res, err := client.Do(req)
	CheckError(err)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Panic("Cannot get token with status code: ", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	CheckError(err)
	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	CheckError(err)
	if token, ok := data["access_token"]; ok {
		return token.(string)
	}
	log.Panic("Cannot get token")
	return ""
}
