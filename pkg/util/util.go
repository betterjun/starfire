package util

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

func init() {
	// 以时间作为初始化种子
	rand.Seed(time.Now().UnixNano())
}

const (
	ascstr = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numstr = "1123456789"
)

// 获取随机字符TOKEN，输出字符数由参数size指定。
func GenNToken(size int) string {
	var bytes = make([]byte, size)
	rand.Read(bytes)
	length := byte(len(ascstr))
	for k, v := range bytes {
		bytes[k] = ascstr[v%length]
	}
	return string(bytes)
}

// 获取随机数字TOKEN，输出字符数由参数size指定。
func GenNumberToken(size int) string {
	var bytes = make([]byte, size)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = numstr[v%10]
	}
	return string(bytes)
}

//
//// GetStartAndEndOf 获取给定时间的一天的开始和结束时间
//func GetStartAndEndOf(t1 time.Time) (start, end time.Time, err error) {
//	if t1.IsZero() {
//		err = fmt.Errorf("参数有误！")
//		return
//	}
//	start, err = time.ParseInLocation(consts.YyyyMmDdHhMmss, t1.Format(consts.YyyyMmDd)+" 00:00:00", time.Local)
//	if err != nil {
//		return
//	}
//	end = start.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
//	return
//}

func MD5(key string) string {
	ha := md5.New()
	ha.Reset()
	ha.Write([]byte(key))

	return hex.EncodeToString(ha.Sum(nil))
}
