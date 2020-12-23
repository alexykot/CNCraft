package mojang

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const mojangSessionServerURL = "https://sessionserver.mojang.com/session/minecraft/hasJoined"

type auth struct {
	UUID       uuid.UUID
	Name       string
	Properties []prop
}

type prop struct {
	Name string
	Data []byte
	Sign []byte
}

type authJson struct {
	ID   string     `json:"id"`
	Name string     `json:"name"`
	Prop []propJson `json:"properties"`
}

type propJson struct {
	Name string `json:"name"`
	Data string `json:"value"`
	Sign string `json:"signature"`
}

func RunMojangSessionAuth(sharedSecret []byte, publicKeyDER []byte, username string) (*auth, error) {
	jsonRes, err := getMojangSessionAuth(username, generateAuthSHAHex(sharedSecret, publicKeyDER))
	if err != nil {
		return nil, fmt.Errorf("failed to get Mojang session auth: %w", err)
	}

	id, err := uuid.Parse(jsonRes.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Mojang profile ID: %w", err)
	}

	auth := auth{
		UUID:       id,
		Name:       jsonRes.Name,
		Properties: make([]prop, 0, len(jsonRes.Prop)),
	}

	for i, jsonProp := range jsonRes.Prop {
		propData, err := base64.StdEncoding.DecodeString(jsonProp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Mojang profile property data b64: %w", err)
		}
		propSign, err := base64.StdEncoding.DecodeString(jsonProp.Sign)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Mojang profile property signature b64: %w", err)
		}

		auth.Properties[i] = prop{
			Name: jsonProp.Name,
			Data: propData,
			Sign: propSign,
		}
	}

	return &auth, nil
}

func getMojangSessionAuth(username, hash string) (*authJson, error) {
	mojangURL := fmt.Sprintf("%s?username=%s&serverId=%s", mojangSessionServerURL, username, hash)

	out, err := http.Get(mojangURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call session server: %w", err)
	}
	response, err := ioutil.ReadAll(out.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var auth authJson
	if err = json.Unmarshal(response, &auth); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response JSON: %w", err)
	}

	return &auth, nil
}

func generateAuthSHAHex(sharedSecret, publicKey []byte) string {
	sha := sha1.New()
	sha.Write(sharedSecret)
	sha.Write(publicKey)
	hash := sha.Sum(nil)

	// Below is Mojang's custom SHA1 encoding. See https://wiki.vg/Protocol_Encryption#Authentication for details.

	// Check for negative hashes
	negative := (hash[0] & 0x80) == 0x80

	if negative {
		carry := true

		for i := len(hash) - 1; i >= 0; i-- {
			hash[i] = ^hash[i]
			if carry {
				carry = hash[i] == 0xff
				hash[i]++
			}
		}
	}

	// Trim away zeroes
	res := strings.TrimLeft(fmt.Sprintf("%x", hash), "0")
	if negative {
		res = "-" + res
	}

	return res
}
