// Package pocket @author K·J Create at 2019-04-09 15:10
package pocket

import (
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	mutex sync.Mutex
)

// GetUUID gene uuid
func GetUUID() uuid.UUID {
	mutex.Lock()
	defer mutex.Unlock()
	return uuid.NewV4()
}

// UnixSecond unix time second
func UnixSecond() int64 {
	return time.Now().Unix()
}

// UnixMillisecond unix time millisecond
func UnixMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}
