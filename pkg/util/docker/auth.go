package util

import (
	"encoding/base64"
	"encoding/json"

	"github.com/docker/docker/api/types"
)

func encodeBase64(auth types.AuthConfig) string {
	authConfigBytes, _ := json.Marshal(auth)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)

	return authConfigEncoded
}
