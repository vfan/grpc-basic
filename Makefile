# gRPC 图书管理服务 Makefile

# 项目配置
PROJECT_NAME := grpc-basic
SERVER_BINARY := server/bookstore-server
CLIENT_BINARY := client/bookstore-client
PROTO_DIR := protos
PB_DIR := pb/bookstore

# Go 相关配置
GO := go
GOFLAGS := -v

# 默认目标
.PHONY: all
all: proto build

# 生成 protobuf 代码
.PHONY: proto
proto:
	@echo "🔧 生成 protobuf 代码..."
	@mkdir -p $(PB_DIR)
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/bookstore.proto
	@mv $(PROTO_DIR)/*.pb.go $(PB_DIR)/
	@echo "✅ protobuf 代码生成完成"

# 构建项目
.PHONY: build
build: proto
	@echo "🔨 构建项目..."
	@$(GO) build $(GOFLAGS) -o $(SERVER_BINARY) ./server
	@$(GO) build $(GOFLAGS) -o $(CLIENT_BINARY) ./client
	@echo "✅ 构建完成"

# 运行服务器
.PHONY: server
server: build
	@echo "🚀 启动图书管理服务器..."
	@./$(SERVER_BINARY)

# 运行客户端
.PHONY: client
client: build
	@echo "📱 启动图书管理客户端..."
	@./$(CLIENT_BINARY)

# 同时运行服务器和客户端（需要两个终端）
.PHONY: run
run: build
	@echo "🎯 启动完整演示..."
	@echo "请在第一个终端运行: make server"
	@echo "请在第二个终端运行: make client"

# 清理构建文件
.PHONY: clean
clean:
	@echo "🧹 清理构建文件..."
	@rm -f $(SERVER_BINARY) $(CLIENT_BINARY)
	@rm -rf $(PB_DIR)/*.pb.go
	@echo "✅ 清理完成"

# 安装依赖
.PHONY: deps
deps:
	@echo "📦 安装项目依赖..."
	@$(GO) mod tidy
	@$(GO) mod download
	@echo "✅ 依赖安装完成"

# 运行测试
.PHONY: test
test:
	@echo "🧪 运行测试..."
	@$(GO) test ./...
	@echo "✅ 测试完成"

# 格式化代码
.PHONY: fmt
fmt:
	@echo "🎨 格式化代码..."
	@$(GO) fmt ./...
	@echo "✅ 代码格式化完成"

# 代码检查
.PHONY: lint
lint:
	@echo "🔍 代码检查..."
	@$(GO) vet ./...
	@echo "✅ 代码检查完成"

# 显示帮助信息
.PHONY: help
help:
	@echo "📚 gRPC 图书管理服务 - 可用命令:"
	@echo ""
	@echo "  make all      - 生成代码并构建项目"
	@echo "  make proto    - 生成 protobuf 代码"
	@echo "  make build    - 构建项目"
	@echo "  make server   - 运行服务器"
	@echo "  make client   - 运行客户端"
	@echo "  make run      - 显示运行说明"
	@echo "  make clean    - 清理构建文件"
	@echo "  make deps     - 安装依赖"
	@echo "  make test     - 运行测试"
	@echo "  make fmt      - 格式化代码"
	@echo "  make lint     - 代码检查"
	@echo "  make help     - 显示此帮助信息"
	@echo ""
	@echo "🚀 快速开始:"
	@echo "  1. make deps    # 安装依赖"
	@echo "  2. make server  # 启动服务器"
	@echo "  3. make client  # 在另一个终端启动客户端" 