// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package aesgcm

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key, err := Key("xVKAGlWQiY3sOp8JVc0nbuNId3PNCgWh")
	if err != nil {
		t.Error(err)
		return
	}

	message := []byte("top-secret")
	ciphertext, err := Encrypt(message, key)
	if err != nil {
		t.Error(err)
		return
	}

	plaintext, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(message, plaintext) {
		t.Errorf("Expect secret encrypted and decrypted")
	}
}

func TestInvalidKey(t *testing.T) {
	_, err := Key("xVKAGlWQiY3sOp8J")
	if err != errInvalidKeyLength {
		t.Errorf("Want Invalid Key Length error")
	}
}
