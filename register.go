package mangoplus

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
)

const RegisterPath = "register"

type RegisterationData struct {
	DeviceSecret string `json:"deviceSecret"`
}

func (c *PlusClient) Register(id string) error {
	u, _ := url.Parse(BaseAPI)
	u = u.JoinPath(RegisterPath)

	// TODO: revise if the extra part of the key is required
	deviceToken := md5Hex(id)
	securityKey := md5Hex(deviceToken + "4Kin9vGg")
	p := map[string]string{
		"device_token": deviceToken,
		"security_key": securityKey,
	}

	res, err := c.Request(context.Background(), http.MethodPut, *u, p, nil, nil)
	if err != nil {
		return err
	}
	if res.Success.RegisterationData == nil {
		return fmt.Errorf("Unexpected error: register data response is empty")
	}
	if res.Success.RegisterationData.DeviceSecret == "" {
		return fmt.Errorf("Unexpected error: device secret response is empty")
	}

	c.secret = &res.Success.RegisterationData.DeviceSecret
	return nil
}

func md5Hex(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}
