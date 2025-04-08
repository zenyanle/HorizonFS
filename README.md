# HorizonFS

HorizonFS 是一个基于 Go 语言构建的分布式文件存储系统原型，采用 Raft 共识协议和 gRPC 实现高可用架构。项目聚焦元数据水平扩展与高效数据流处理，为构建大规模分布式存储系统提供基础框架。

---

## 核心特性 ✨

### 🚀 元数据架构

#### 为扩展设计
- 基于键值哈希的元数据分片机制，支持将元数据分布到多个 Raft Group
- 客户端可以通过哈希计算直接定位目标 Raft Group，支持向任意节点（Leader/Follower）发起查询，实现去中心化查询，防止单点故障，提供了高可用性

#### 基础实现
- 提供完整的单 Raft Group 实现 (cmd/node)，包含元数据一致性保证、分布式存储逻辑和通信协议栈

### 🔒 强一致性保障
- 基于 etcd/raft 实现 Raft 共识协议
- 内存状态机 (Metastore) 实时应用已提交日志，提供快速元数据查询响应

### ⚡ 高性能数据管道

#### 智能流式传输
- 基于 gRPC 双向流实现分块传输协议（默认 64KB），支持大文件的高效传输
- 动态内存管理：分块处理机制降低了内存峰值消耗

#### 持久化引擎
- 采用 BadgerDB 嵌入式存储引擎，读写性能高

### 🧩 模块化架构

| 模块            | 功能描述               |
|-----------------|------------------------|
| internal/client | 智能客户端 SDK         |
| internal/node   | 分布式节点核心逻辑     |
| pkg/proto       | gRPC 服务接口定义       |
| cmd/*           | 命令行工具集           |

---

## 快速部署 🚀

### 环境要求
- Go 1.23+
- Make
- protoc 编译器

### 编译指南

```bash
# 构建节点程序
go build -o horizonfs-node ./cmd/node/main.go

# 构建管理客户端
go build -o hfs-cli ./cmd/client/main.go
```

---

## 集群启动示例

```bash
# 启动 3 节点集群（单 Raft Group 模式）

# 节点 1 (Leader)
./horizonfs-node -id 1 -cluster http://127.0.0.1:9021,http://127.0.0.1:9022,http://127.0.0.1:9023 -port 9121

# 节点 2 (Follower)
./horizonfs-node -id 2 -cluster ... -port 9122 -join  # 新终端执行

# 节点 3 (Follower)
./horizonfs-node -id 3 -cluster ... -port 9123 -join  # 新终端执行
```

---

## 客户端操作示例

```bash
# 上传文件到集群
./hfs-cli -addr=localhost:9121 upload -local=./sample.txt -remote=remote_sample.txt

# 从集群下载文件
./hfs-cli -addr=localhost:9121 download -remote=remote_sample.txt -local=./downloaded.txt

# 获取文件元数据（任意节点可查）
./hfs-cli -addr=localhost:9122 get-meta -remote=remote_sample.txt
```

---

## 架构规划

### 当前
- ✅ 单 Raft Group 元数据一致性
- ✅ 流式数据传输基础框架
- ✅ 高密度元数据存储

### 愿景
- 🛠 多 Raft Group 元数据分片
- 🛠 数据多副本机制 (3x 冗余)
- 🛠 基于一致性哈希的请求路由

---

## 技术矩阵

| 领域       | 技术选型                      |
|------------|-------------------------------|
| 共识协议   | etcd/raft v3                  |
| 存储引擎   | BadgerDB v4     |
| 通信协议   | gRPC with streaming           |
| 序列化     | Protocol Buffers v3           |
| 性能优化   | Sonic JSON + SIMD 加速        |

---

## 📁 项目结构总览

```text
HorizonFS/
├── cmd/                         # 命令行工具入口
│   ├── client/                  # 客户端主程序
│   │   └── main.go
│   └── node/                    # 节点主程序
│       └── main.go
├── configs/                     # 配置文件目录（预留）
├── go.mod                       # Go 依赖管理文件
├── go.sum
├── internal/                    # 项目核心逻辑
│   ├── client/                  # 客户端逻辑
│   │   ├── remote/              # 与节点交互逻辑
│   │   └── store/               # 本地文件操作封装
│   └── node/                    # 分布式节点核心模块
│       ├── grpcserver/          # gRPC 服务注册与启动
│       ├── kvstore/             # Badger 封装与文件处理
│       ├── metadata/            # 元数据结构与存储
│       ├── raft/                # Raft 共识实现封装
│       ├── service/             # gRPC 接口服务实现
│       └── types/               # 通用结构定义
├── pkg/                         # 公共库
│   ├── logger/                  # 日志组件
│   └── proto/                   # Proto 文件与生成代码
├── scripts/                     # 构建与测试脚本
│   ├── client_test_build.sh
│   ├── generate_datanode_proto.sh
│   ├── generate_metadata_proto.sh
│   └── node_test_build.sh
├── tools/                       # 调试与开发辅助工具
│   ├── grpc-debug-client/
│   │   └── main.go
│   └── grpc-debug-node/
│       └── main.go
└── README.md                    # 项目说明文档
```

## Inspiration

- [etcd-io/raftexample](https://github.com/etcd-io/etcd/blob/main/contrib/raftexample/README.md) — A minimal Raft usage example in Go from the etcd team
- [zenyanle/Gflow](https://github.com/zenyanle/Gflow) — A lightweight and efficient edge computing data interaction platform


