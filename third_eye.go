/*
 * Copyright 2017 Yoshihiro Tanaka
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Author: Yoshihiro Tanaka <contact@cordea.jp>
 * date  : 2017-05-16
 */

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ThirdEye struct {
	Key       string
	IsDecrypt bool
}

func (t *ThirdEye) Rename(path string) error {
	var name string
	var err error
	if t.IsDecrypt {
		name, err = t.decrypt(filepath.Base(path))
	} else {
		name, err = t.encrypt(filepath.Base(path))
	}
	if err != nil {
		return err
	}
	newPath := filepath.Join(filepath.Dir(path), name)
	if _, err := os.Stat(newPath); err == nil {
		return errors.New("File already exists.")
	}
	return os.Rename(path, newPath)
}

func (t *ThirdEye) encrypt(name string) (string, error) {
	bytes := []byte(name)

	gcm, err := t.gcm()
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	result := gcm.Seal(nonce, nonce, bytes, nil)
	return base64.URLEncoding.EncodeToString(result), nil
}

func (t *ThirdEye) decrypt(name string) (string, error) {
	bytes, err := base64.URLEncoding.DecodeString(name)
	if err != nil {
		return "", err
	}

	gcm, err := t.gcm()
	if err != nil {
		return "", err
	}

	size := gcm.NonceSize()
	if len(bytes) < size {
		return "", errors.New(fmt.Sprintf("Invalid file name: %s.", name))
	}

	nonce, bytes := bytes[:size], bytes[size:]
	result, err := gcm.Open(nil, nonce, bytes, nil)
	return string(result), err
}

func (t *ThirdEye) gcm() (cipher.AEAD, error) {
	key := []byte(t.Key)
	cr, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(cr)
}
