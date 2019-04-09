// Package pocket @author K·J Create at 2019-04-09 15:13
package pocket

import "encoding/json"

// Resp common response
type Resp struct {
	Code
	Data interface{} `json:"data"`
}

// NewResp param to Object
func NewResp(error *Code, data interface{}) Resp {
	return Resp{
		*error,
		data,
	}
}

// Marshal to json []byte
func (resp *Resp) Marshal() []byte {
	b, err := json.Marshal(resp)
	if err != nil {
		return []byte(err.Error())
	}
	return b
}
