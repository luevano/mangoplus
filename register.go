package mangoplus

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
)

const RegisterPath = "register"

// ensureRegistered: Registers a new device if there is no secret yet.
func (c *PlusClient) ensureRegistered(ctx context.Context) error {
	c.registerMu.Lock()
	defer c.registerMu.Unlock()

	if c.secret != "" {
		return nil
	}
	return c.register(ctx)
}

// register: Registers a new device, same as the official app on a fresh install.
//
// The resulting secret only has access to free chapters.
func (c *PlusClient) register(ctx context.Context) error {
	deviceToken := md5Hex(c.options.AndroidID)
	securityKey := md5Hex(deviceToken + "4Kin9vGg")
	p := map[string]string{
		"device_token": deviceToken,
		"security_key": securityKey,
	}

	res, err := c.request(ctx, http.MethodPut, RegisterPath, p)
	if err != nil {
		return fmt.Errorf("failed to register device: %w", err)
	}
	secret := res.GetRegistrationData().GetDeviceSecret()
	if secret == "" {
		return fmt.Errorf("Unexpected error: device secret response is empty")
	}

	c.secret = secret
	if c.options.OnRegister != nil {
		c.options.OnRegister(secret)
	}
	return nil
}

func md5Hex(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}
