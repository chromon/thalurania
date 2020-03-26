package sequence

import (
	"errors"
	"sync"
	"time"
)

// 因为 snowFlake 目的是解决分布式下生成唯一 id 所以 id 中是包含集群和节点编号在内的
// 64 位整形
// 第一位 bit：表示正负数，所需 id 为正数，所以值为 0
// 41位 bit：生成 id 时的毫秒时间戳，范围 0 ~ 2^41 -1
// 10位 bit：工作机器的 id，所以允许分布式最大节点数为 1024 个
// 后 12 位 bit：表示单台机器每毫秒生成的 id 序号 0 ~ 2^12 - 1

const (
	// 每台机器(节点)的 ID 位数 10 位最大可以有 2^10=1024 个节点
	workerBits uint8 = 10
	// 表示每个集群下的每个节点，1毫秒内可生成的 id 序号的二进制位数 即每毫秒可生成 2^12-1=4096 个唯一 ID
	numberBits uint8 = 12

	// 这里求最大值使用了位运算，-1 的二进制表示为 1 的补码
	// 节点 Id 的最大值，用于防止溢出
	workerMax   int64 = -1 ^ (-1 << workerBits)
	// 同上，用来表示生成 id 序号的最大值
	numberMax   int64 = -1 ^ (-1 << numberBits)
	// 时间戳向左的偏移量
	timeShift   uint8 = workerBits + numberBits
	// 节点ID向左的偏移量
	workerShift uint8 = numberBits
	// 41位字节作为时间戳数值的话 大约68年就会用完
	// 假如2010年1月1日开始开发系统 如果不减去2010年1月1日的时间戳 那么会浪费40年的时间戳！
	// 这个一旦定义且开始生成 ID 后千万不要改了 不然可能会生成相同的 ID
	// 这个是当前生成 epoch 这个变量时的时间戳(毫秒)
	// time.Now()：2020-03-26 09:36:34.073587 +0800 CST m=+0.028981701
	// time.Now().UnixNano() / 1e6：1585186594049
	epoch int64 = 1585186594049
)

// 定义一个 worker 工作节点所需要的基本参数
type Worker struct {
	// 添加互斥锁 确保并发安全
	mu        sync.Mutex
	// 记录时间戳
	timestamp int64
	// 该节点的ID
	workerId  int64
	// 当前毫秒已经生成的id序列号(从0开始累加) 1毫秒内最多生成4096个ID
	number    int64
}

// 实例化一个工作节点
func NewWorker(workerId int64) (*Worker, error) {
	// 要先检测workerId是否在上面定义的范围内
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("worker id excess of quantity")
	}
	// 生成一个新节点
	return &Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

// 生成 id，生成方法一定要挂载在某个 worker 下，这样逻辑会比较清晰 指定某个节点生成id
func (w *Worker) GetId() int64 {
	// 获取id最关键的一点 加锁 加锁 加锁
	w.mu.Lock()
	// 生成完成后记得 解锁 解锁 解锁
	defer w.mu.Unlock()

	// 获取生成时的时间戳
	// 纳秒转毫秒
	now := time.Now().UnixNano() / 1e6
	if w.timestamp == now {
		w.number++

		// 这里要判断，当前工作节点是否在1毫秒内已经生成 numberMax 个ID
		if w.number > numberMax {
			// 如果当前工作节点在1毫秒内生成的 ID 已经超过上限 需要等待1毫秒再继续生成
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 如果当前时间与工作节点上一次生成 ID 的时间不一致 则需要重置工作节点生成 ID 的序号
		w.number = 0
		w.timestamp = now // 将机器上一次生成 ID 的时间更新为当前时间
	}

	// now - epoch 为该算法目前已经奔跑了xxx毫秒
	// 如果在程序跑了一段时间修改了 epoch 这个值 可能会导致生成相同的 ID
	ID := int64((now-epoch)<<timeShift | (w.workerId << workerShift) | (w.number))

	return ID
}