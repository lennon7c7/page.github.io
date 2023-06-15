package util

import (
	"strings"
)

// 身份证部分信息改成星号
func MaskIDCard(IDCard string) (maskIDCard string) {
	idCardAfterArray := strings.Split(IDCard, "")
	idCardAfterArray[8] = "*"
	idCardAfterArray[9] = "*"
	idCardAfterArray[10] = "*"
	idCardAfterArray[11] = "*"
	maskIDCard = strings.Join(idCardAfterArray, "")

	return
}
