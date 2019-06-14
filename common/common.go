package common

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"reflect"
	"sync"
	"time"
)

// IsEmpty 判读数据是否为空
func IsEmpty(a interface{}) bool {
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

var Locker = make(map[string]*sync.RWMutex)

func Lock(index string) {
	for {
		_, ok := Locker[index]
		if !ok {
			Locker[index] = &sync.RWMutex{}
			break
		}
		//100ms轮训一次状态
		time.Sleep(100 * time.Millisecond)
	}

	Locker[index].Lock()
}

func Unlock(index string) {
	Locker[index].Unlock()
	//删除使用过的锁，避免map无限增加
	delete(Locker, index)
}

//生成32位md5字串
func getMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateUniqueId 生成Guid字串
func GenerateUniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return getMd5String(base64.URLEncoding.EncodeToString(b))
}
