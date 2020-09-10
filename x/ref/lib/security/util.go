// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package security

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"os"
)

// DefaultSSHAgentSockNameFunc can be overridden to return the address of a custom
// ssh agent to use instead of the one specified by SSH_AUTH_SOCK. This is
// primarily intended for tests.
var DefaultSSHAgentSockNameFunc = func() string {
	return os.Getenv("SSH_AUTH_SOCK")
}

// NewPrivateKey creates a new private key of the requested type.
// keyType must be one of ecdsa256, ecdsa384, ecdsa521 or ed25519.
func NewPrivateKey(keyType string) (interface{}, error) {
	switch keyType {
	case "ecdsa256":
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "ecdsa384":
		return ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "ecdsa521":
		return ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	case "ed25519":
		_, privateKey, err := ed25519.GenerateKey(rand.Reader)
		return privateKey, err
	default:
		return nil, fmt.Errorf("unsupported key type: %T", keyType)
	}
}
