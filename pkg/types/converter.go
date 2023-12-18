package types

import (
	"GoBlog/pkg/logger"
	"strconv"
)

//Uint64ToString 将 uint64 转换为 string
func Uint64ToString(num uint64) string {
    return strconv.FormatUint(num, 10)
}

//StringToUint64 奖 string 转换为 uint64
func StringToUint64(str string) uint64 {
    i, err := strconv.ParseUint(str, 10, 64)
    if err != nil {
        logger.LogError(err)
    }

    return i
}