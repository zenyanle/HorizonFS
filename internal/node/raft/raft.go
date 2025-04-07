// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// raft包提供了一个基于etcd/raft的分布式一致性实现
package raft

import (
	"fmt"

	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	"go.etcd.io/etcd/server/v3/wal"

	"HorizonFS/internal/node/types"
	"go.uber.org/zap"
)

// commit 表示已提交的数据条目
type commit = types.Commit

// raftNode 表示一个由raft支持的键值流
type raftNode struct {
	proposeC    <-chan string            // 提议的消息(k,v)
	confChangeC <-chan raftpb.ConfChange // 提议的集群配置变更
	commitC     chan<- *commit           // 提交到日志的条目(k,v)
	errorC      chan<- error             // raft会话中的错误

	id          int      // raft会话的客户端ID
	peers       []string // raft对等节点URL
	join        bool     // 节点是否正在加入现有集群
	waldir      string   // WAL目录的路径
	snapdir     string   // 快照目录的路径
	getSnapshot func() ([]byte, error)

	confState     raftpb.ConfState
	snapshotIndex uint64
	appliedIndex  uint64

	// 为commit/error通道提供raft支持
	node        raft.Node
	raftStorage *raft.MemoryStorage
	wal         *wal.WAL

	snapshotter      *snap.Snapshotter
	snapshotterReady chan *snap.Snapshotter // 表示snapshotter准备好时的信号

	snapCount uint64
	transport *rafthttp.Transport
	stopc     chan struct{} // 表示提议通道关闭的信号
	httpstopc chan struct{} // 表示http服务器关闭的信号
	httpdonec chan struct{} // 表示http服务器完成关闭的信号

	logger *zap.Logger
}

var defaultSnapshotCount uint64 = 10000

// NewRaftNode 初始化一个raft实例，返回已提交日志条目的通道和错误通道。
// 日志更新的提议通过proposeC通道发送。所有日志条目通过commitC重放，
// 随后是一个nil消息(表示通道是最新的)，然后是新的日志条目。
// 要关闭，请关闭proposeC并读取errorC。
func NewRaftNode(id int, peers []string, join bool, getSnapshot func() ([]byte, error), proposeC <-chan string,
	confChangeC <-chan raftpb.ConfChange) (<-chan *commit, <-chan error, <-chan *snap.Snapshotter) {

	commitC := make(chan *commit)
	errorC := make(chan error)

	rc := &raftNode{
		proposeC:    proposeC,
		confChangeC: confChangeC,
		commitC:     commitC,
		errorC:      errorC,
		id:          id,
		peers:       peers,
		join:        join,
		waldir:      fmt.Sprintf("raft-%d", id),
		snapdir:     fmt.Sprintf("raft-%d-snap", id),
		getSnapshot: getSnapshot,
		snapCount:   defaultSnapshotCount,
		stopc:       make(chan struct{}),
		httpstopc:   make(chan struct{}),
		httpdonec:   make(chan struct{}),

		logger: zap.NewExample(),

		snapshotterReady: make(chan *snap.Snapshotter, 1),
		// WAL重放后填充结构的其余部分
	}
	go rc.startRaft()
	return commitC, errorC, rc.snapshotterReady
}
