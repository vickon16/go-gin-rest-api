package utils

import "strconv"

func ConstructRedisUserKey(userId int64) string {
	cacheKey := "tutorial:user:" + strconv.FormatInt(userId, 10)
	return cacheKey
}
