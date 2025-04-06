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

package raft

import (
	log "HorizonFS/pkg/logger"
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"

	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
)

// serveRaft 启动一个HTTP服务器处理raft消息
func (rc *raftNode) serveRaft() {
	url, err := url.Parse(rc.peers[rc.id-1])
	if err != nil {
		log.Fatalf("raftexample: 解析URL失败 (%v)", err)
	}

	ln, err := newStoppableListener(url.Host, rc.httpstopc)
	if err != nil {
		log.Fatalf("raftexample: 监听rafthttp失败 (%v)", err)
	}

	err = (&http.Server{Handler: rc.transport.Handler()}).Serve(ln)
	select {
	case <-rc.httpstopc:
	default:
		log.Fatalf("raftexample: 服务rafthttp失败 (%v)", err)
	}
	close(rc.httpdonec)
}

// stopHTTP 关闭HTTP服务器
func (rc *raftNode) stopHTTP() {
	rc.transport.Stop()
	close(rc.httpstopc)
	<-rc.httpdonec
}

// Process 实现raft.Transport接口
func (rc *raftNode) Process(ctx context.Context, m raftpb.Message) error {
	return rc.node.Step(ctx, m)
}

// IsIDRemoved 实现raft.Transport接口
func (rc *raftNode) IsIDRemoved(id uint64) bool { return false }

// ReportUnreachable 实现raft.Transport接口
func (rc *raftNode) ReportUnreachable(id uint64) { rc.node.ReportUnreachable(id) }

// ReportSnapshot 实现raft.Transport接口
func (rc *raftNode) ReportSnapshot(id uint64, status raft.SnapshotStatus) {
	rc.node.ReportSnapshot(id, status)
}

// newStoppableListener 创建一个可停止的TCP监听器
func newStoppableListener(addr string, stopc <-chan struct{}) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &stoppableListener{ln, stopc}, nil
}

type stoppableListener struct {
	net.Listener
	stopc <-chan struct{}
}

func (ln stoppableListener) Accept() (net.Conn, error) {
	connc := make(chan net.Conn, 1)
	errc := make(chan error, 1)
	go func() {
		conn, err := ln.Listener.Accept()
		if err != nil {
			errc <- err
			return
		}
		connc <- conn
	}()
	select {
	case <-ln.stopc:
		return nil, errors.New("server stopped")
	case err := <-errc:
		return nil, err
	case conn := <-connc:
		return conn, nil
	}
}
