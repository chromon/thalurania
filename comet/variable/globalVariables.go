package variable

import (
	"chalurania/service/db/conn"
	"chalurania/service/pubsub"
	"chalurania/service/sequence"
)

// redis 连接池
var RedisPool *pubsub.RedisPool

// mysql 连接
var GoDB *conn.GoDB

// id 生成器
var IdWorker *sequence.Worker
