# gRPC 图书管理服务 - 完整实战项目

## 📚 项目概述

这是一个完整的 gRPC 图书管理服务项目，展示了如何使用 gRPC 构建微服务。项目包含服务端、客户端、完整的测试用例和详细的文档。

### 🎯 项目特性

- ✅ 完整的 CRUD 操作（创建、读取、更新、删除）
- ✅ 分页查询功能
- ✅ 按价格区间搜索
- ✅ 详细的错误处理和日志记录
- ✅ 完整的单元测试
- ✅ 中文注释和文档
- ✅ 使用 Makefile 简化构建流程

### 📁 项目结构

```
grpc-basic/
├── protos/                    # Protocol Buffers 定义文件
│   └── bookstore.proto       # 图书服务接口定义
├── pb/                       # 生成的 protobuf 代码
│   └── bookstore/
│       ├── bookstore.pb.go   # 消息类型定义
│       └── bookstore_grpc.pb.go # 服务接口定义
├── server/                   # 服务端代码
│   ├── main.go              # 服务端主程序
│   └── server_test.go       # 服务端单元测试
├── client/                   # 客户端代码
│   └── main.go              # 客户端演示程序
├── Makefile                  # 构建和运行脚本
├── go.mod                    # Go 模块定义
└── README.md                 # 项目文档
```

## 🚀 快速开始

### 1. 安装依赖

```bash
# 安装项目依赖
make deps

# 或者手动安装
go mod tidy
```

### 2. 生成代码

```bash
# 生成 protobuf 代码
make proto

# 或者手动生成
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    protos/bookstore.proto
```
将自动生成的代码复制到server和client目录下