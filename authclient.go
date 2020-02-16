package gortsplib

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func md5Hex(in string) string {
	h := md5.New()
	h.Write([]byte(in))
	return hex.EncodeToString(h.Sum(nil))
}

// AuthClient is an object that helps a client sending its credentials to a server.
type AuthClient struct {
	user  string
	pass  string
	realm string
	nonce string
}

// NewAuthClient allocates an AuthClient.
// header is the WWW-Authenticate header provided by the server.
func NewAuthClient(header []string, user string, pass string) (*AuthClient, error) {
	headerAuthDigest := func() string {
		for _, v := range header {
			if strings.HasPrefix(v, "Digest ") {
				return v
			}
		}
		return ""
	}()
	if headerAuthDigest == "" {
		return nil, fmt.Errorf("Authentication/Digest header not provided")
	}

	auth, err := ReadHeaderAuth(headerAuthDigest)
	if err != nil {
		return nil, err
	}

	nonce, ok := auth.Values["nonce"]
	if !ok {
		return nil, fmt.Errorf("nonce not provided")
	}

	realm, ok := auth.Values["realm"]
	if !ok {
		return nil, fmt.Errorf("realm not provided")
	}

	return &AuthClient{
		user:  user,
		pass:  pass,
		realm: realm,
		nonce: nonce,
	}, nil
}

// GenerateHeader generates an Authorization Header that allows to authenticate a request with
// the given method and path.
func (ac *AuthClient) GenerateHeader(method string, path string) []string {
	ha1 := md5Hex(ac.user + ":" + ac.realm + ":" + ac.pass)
	ha2 := md5Hex(method + ":" + path)
	response := md5Hex(ha1 + ":" + ac.nonce + ":" + ha2)

	return []string{fmt.Sprintf("Digest username=\"%s\", realm=\"%s\", nonce=\"%s\", uri=\"%s\", response=\"%s\"",
		ac.user, ac.realm, ac.nonce, path, response)}
}
