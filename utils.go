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

// pow 次方
func pow(x, n int) int32 {
	ret := 1 // 结果初始为0次方的值，整数0次方为1。如果是矩阵，则为单元矩阵。
	for n != 0 {
		if n%2 != 0 {
			ret = ret * x
		}
		n /= 2
		x = x * x
	}
	return int32(ret)
}

// CreateCaptcha random
func CreateCaptcha(digit int) string {
	format := "%0" + fmt.Sprintf("%d", digit) + "v"
	return fmt.Sprintf(format, rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(pow(10, digit)))
}

// md5 md5 32 lowercase
func Md5(txt string) string {
	md5Hash := md5.New()
	io.WriteString(md5Hash, txt)
	md5Bytes := md5Hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(md5Bytes))
}
