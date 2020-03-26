package main

import (
	"chalurania/service/log"
	"chalurania/service/sequence"
	"testing"
)

var IdWorker *sequence.Worker

func TestSnowFlake(t *testing.T) {
	// 生成一个节点实例
	// 传入当前节点id 此id在机器集群中一定要唯一 且从0开始排最多1024个节点，可以根据节点的不同动态调整该算法每毫秒生成的id上限
	IdWorker, _ = sequence.NewWorker(0)

	// 获得唯一id
	id := IdWorker.GetId()
	log.Info.Println(id)
}

func BenchmarkSnowFlake(t *testing.B) {
	IdWorker, _ = sequence.NewWorker(0)

	for i := 0; i < t.N; i++ {
		IdWorker.GetId()
	}
}