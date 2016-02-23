package utils

import (
    "strings"

    "github.com/nu7hatch/gouuid"
)

func RandomKey() (string, error) {
    keyUuid, err := uuid.NewV4()
    if err != nil {
        return "", err
    }
    key := strings.Replace(keyUuid.String(), "-", "", -1)
    return key, nil
}
