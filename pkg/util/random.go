package util

import (
	"math/rand"
	"strconv"
	"time"
)

// 随机生成-订单编号（28位数）
// 返回例子：2021030416064025497821967038
func RandomOrderSn() (orderSn string) {
	orderSnInt64 := RandomVersion()

	orderSn = strconv.FormatInt(orderSnInt64, 10)
	orderSn = StampToString(NowTimeMs(), "20060102150405") + orderSn

	return
}

// 随机生成-乐观锁（14位数）
// 返回例子：25497821967038
func RandomVersion() (version int64) {
	min := 10000000000000
	max := 99999999999999
	rand.NewSource(time.Now().Unix())
	version = int64(min + rand.Intn(max-min))

	return
}
