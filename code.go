// Package pocket @author K·J Create at 2019-04-09 15:11
package pocket

import (
	"encoding/json"
	"fmt"
)

// Code object
type Code struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// NewCode get object
func NewCode(code int, msg string) *Code {
	return &Code{
		Code: code,
		Msg:  msg,
	}
}

// Code to json
func (c *Code) Error() string {
	b, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf(`{"code": %d, "msg": "%s"}`, c.Code, c.Msg)
	}
	return string(b)
}
