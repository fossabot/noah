/*
 * Copyright (c) 2018 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package security

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

// BcryptCost is the cost to use when hashing passwords. It is exposed for
// testing.
//
// BcryptCost should increase along with computation power.
// For estimates, see: http://security.stackexchange.com/questions/17207/recommended-of-rounds-for-bcrypt
// For now, we use the library's default cost.
var BcryptCost = bcrypt.DefaultCost

// ErrEmptyPassword indicates that an empty password was attempted to be set.
var ErrEmptyPassword = errors.New("empty passwords are not permitted")

// CompareHashAndPassword tests that the provided bytes are equivalent to the
// hash of the supplied password. If they are not equivalent, returns an
// error.
func CompareHashAndPassword(hashedPassword []byte, password string) error {
	h := sha256.New()
	return bcrypt.CompareHashAndPassword(hashedPassword, h.Sum([]byte(password)))
}

// HashPassword takes a raw password and returns a bcrypt hashed password.
func HashPassword(password string) ([]byte, error) {
	h := sha256.New()
	return bcrypt.GenerateFromPassword(h.Sum([]byte(password)), BcryptCost)
}

// PromptForPassword prompts for a password.
// This is meant to be used when using a password.
func PromptForPassword() (string, error) {
	fmt.Print("Enter password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	// Make sure stdout moves on to the next line.
	fmt.Print("\n")

	return string(password), nil
}

// PromptForPasswordTwice prompts for a password twice, returning the read string if
// they match, or an error.
// This is meant to be used when setting a password.
func PromptForPasswordTwice() (string, error) {
	fmt.Print("Enter password: ")
	one, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	if len(one) == 0 {
		return "", ErrEmptyPassword
	}
	fmt.Print("\nConfirm password: ")
	two, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	// Make sure stdout moves on to the next line.
	fmt.Print("\n")
	if !bytes.Equal(one, two) {
		return "", errors.New("password mismatch")
	}

	return string(one), nil
}
