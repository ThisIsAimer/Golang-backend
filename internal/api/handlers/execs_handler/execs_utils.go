package execs

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"simpleapi/pkg/utils"
	"strings"

	"golang.org/x/crypto/argon2"
)

func getModelTags(model any) []string {
	tags := []string{}
	modelType := reflect.TypeOf(model)
	for i := range modelType.NumField() {
		tag := modelType.Field(i).Tag.Get("json")
		tag = strings.TrimSuffix(tag, `,omitempty`)
		tags = append(tags, tag)
	}

	return tags
}

func passEncoder(password string, salt []byte) (string, error) {
	if password == "" {
		return "", utils.ErrorHandler(fmt.Errorf("password is empty"), "password is required")
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodedHash := saltBase64 + "." + hashBase64

	return encodedHash, nil
}