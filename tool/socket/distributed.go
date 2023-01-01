package socket

import "medicinal_share/tool/db/redis"

var sole = redis.GetSole("socket", 1)

func getSocketId() string {
	return sole.GetID()
}
