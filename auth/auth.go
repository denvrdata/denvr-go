package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/denvrdata/go-denvr/result"
)

type Auth struct {
	Server         string
	AccessToken    string
	RefreshToken   string
	AccessExpires  int64
	RefreshExpires int64
}

func NewAuth(server string, username string, password string) Auth {
	data := result.Wrap(
		json.Marshal(
			map[string]string{
				"userNameOrEmailAddress": username,
				"password":               password,
			},
		),
	).Unwrap()

	resp := result.Wrap(
		http.Post(
			fmt.Sprintf("%s/api/TokenAuth/Authenticate", server),
			"application/json",
			bytes.NewBuffer(data),
		),
	).Unwrap()

	defer resp.Body.Close()

	// A bit ugly, but we'll define our specific response content to decode
	var content struct {
		Result struct {
			AccessToken                 string `json:"accessToken"`
			RefreshToken                string `json:"refreshToken"`
			ExpireInSeconds             int64  `json:"expireInSeconds"`
			RefreshTokenExpireInSeconds int64  `json:"refreshTokenExpireInSeconds"`
		} `json:"result"`
	}
	result.Wrap(content, json.NewDecoder(resp.Body).Decode(&content)).Unwrap()

	return Auth{
		server,
		content.Result.AccessToken,
		content.Result.RefreshToken,
		content.Result.ExpireInSeconds,
		content.Result.RefreshTokenExpireInSeconds,
	}
}

func (auth Auth) Token() {
	time.Now().Unix()
}
