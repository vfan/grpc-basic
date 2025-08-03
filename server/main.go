package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	// 导入生成的protobuf代码
	pb "grpc-basic-server/pb"

	// 导入gRPC相关包
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BookServer 实现图书管理服务
type BookServer struct {
	// 嵌入未实现的服务接口，确保向后兼容
	pb.UnimplementedBookServiceServer

	// 互斥锁，用于保护并发访问
	mu sync.RWMutex

	// 内存中的图书存储（实际项目中应该使用数据库）
	books map[string]*pb.Book

	// 用于生成唯一ID的计数器
	idCounter int64
}

// NewBookServer 创建新的图书服务器实例
func NewBookServer() *BookServer {
	return &BookServer{
		books: make(map[string]*pb.Book),
	}
}

// generateID 生成唯一的图书ID
func (s *BookServer) generateID() string {
	s.idCounter++
	return fmt.Sprintf("book-%d", s.idCounter)
}

// CreateBook 创建图书
func (s *BookServer) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	// 记录请求日志
	log.Printf("收到创建图书请求: %v", req.GetBook().GetTitle())

	// 获取请求中的图书信息
	book := req.GetBook()

	// 验证图书信息
	if book.GetTitle() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "图书标题不能为空")
	}
	if book.GetAuthor() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "作者不能为空")
	}
	if book.GetPrice() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "图书价格必须大于0")
	}

	// 加写锁保护并发访问
	s.mu.Lock()
	defer s.mu.Unlock()

	// 生成唯一ID
	bookID := s.generateID()
	book.Id = bookID

	// 存储图书信息
	s.books[bookID] = book

	log.Printf("成功创建图书，ID: %s", bookID)

	// 返回成功响应
	return &pb.CreateBookResponse{
		Id:      bookID,
		Message: "图书创建成功",
	}, nil
}

// GetBook 获取图书信息
func (s *BookServer) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	// 记录请求日志
	log.Printf("收到获取图书请求，ID: %s", req.GetId())

	// 验证请求参数
	if req.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "图书ID不能为空")
	}

	// 加读锁保护并发访问
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 查找图书
	book, exists := s.books[req.GetId()]
	if !exists {
		log.Printf("图书未找到，ID: %s", req.GetId())
		return nil, status.Errorf(codes.NotFound, "图书不存在，ID: %s", req.GetId())
	}

	log.Printf("成功获取图书，ID: %s", req.GetId())

	// 返回图书信息
	return &pb.GetBookResponse{
		Book: book,
	}, nil
}

// UpdateBook 更新图书信息
func (s *BookServer) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	// 记录请求日志
	log.Printf("收到更新图书请求，ID: %s", req.GetBook().GetId())

	// 获取要更新的图书信息
	book := req.GetBook()

	// 验证请求参数
	if book.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "图书ID不能为空")
	}
	if book.GetTitle() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "图书标题不能为空")
	}
	if book.GetAuthor() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "作者不能为空")
	}
	if book.GetPrice() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "图书价格必须大于0")
	}

	// 加写锁保护并发访问
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查图书是否存在
	if _, exists := s.books[book.GetId()]; !exists {
		log.Printf("图书不存在，无法更新，ID: %s", book.GetId())
		return nil, status.Errorf(codes.NotFound, "图书不存在，ID: %s", book.GetId())
	}

	// 更新图书信息
	s.books[book.GetId()] = book

	log.Printf("成功更新图书，ID: %s", book.GetId())

	// 返回成功响应
	return &pb.UpdateBookResponse{
		Message: "图书更新成功",
	}, nil
}

// DeleteBook 删除图书
func (s *BookServer) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	// 记录请求日志
	log.Printf("收到删除图书请求，ID: %s", req.GetId())

	// 验证请求参数
	if req.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "图书ID不能为空")
	}

	// 加写锁保护并发访问
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查图书是否存在
	if _, exists := s.books[req.GetId()]; !exists {
		log.Printf("图书不存在，无法删除，ID: %s", req.GetId())
		return nil, status.Errorf(codes.NotFound, "图书不存在，ID: %s", req.GetId())
	}

	// 删除图书
	delete(s.books, req.GetId())

	log.Printf("成功删除图书，ID: %s", req.GetId())

	// 返回成功响应
	return &pb.DeleteBookResponse{
		Message: "图书删除成功",
	}, nil
}

// ListBooks 列出所有图书（支持分页）
func (s *BookServer) ListBooks(ctx context.Context, req *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	// 记录请求日志
	log.Printf("收到列出图书请求，页码: %d, 每页大小: %d", req.GetPage(), req.GetPageSize())

	// 设置默认分页参数
	page := req.GetPage()
	if page <= 0 {
		page = 1
	}

	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大页面大小
	}

	// 加读锁保护并发访问
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 计算总数量
	total := int32(len(s.books))

	// 计算分页参数
	start := (page - 1) * pageSize
	end := start + pageSize

	// 收集图书列表
	var books []*pb.Book
	count := int32(0)
	for _, book := range s.books {
		if count >= start && count < end {
			books = append(books, book)
		}
		count++
	}

	log.Printf("成功列出图书，总数: %d, 当前页: %d", total, page)

	// 返回图书列表
	return &pb.ListBooksResponse{
		Books: books,
		Total: total,
	}, nil
}

// SearchBooksByPrice 按价格区间查询图书
func (s *BookServer) SearchBooksByPrice(ctx context.Context, req *pb.SearchBooksByPriceRequest) (*pb.SearchBooksByPriceResponse, error) {
	// 记录请求日志
	log.Printf("收到按价格查询图书请求，价格区间: %.2f - %.2f", req.GetMinPrice(), req.GetMaxPrice())

	// 验证价格参数
	minPrice := req.GetMinPrice()
	maxPrice := req.GetMaxPrice()

	if minPrice < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "最低价格不能为负数")
	}
	if maxPrice < minPrice {
		return nil, status.Errorf(codes.InvalidArgument, "最高价格不能小于最低价格")
	}

	// 加读锁保护并发访问
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 查找符合条件的图书
	var books []*pb.Book
	for _, book := range s.books {
		price := book.GetPrice()
		if price >= minPrice && price <= maxPrice {
			books = append(books, book)
		}
	}

	log.Printf("按价格查询完成，找到 %d 本图书", len(books))

	// 返回查询结果
	return &pb.SearchBooksByPriceResponse{
		Books: books,
	}, nil
}

// 日志拦截器 - 记录所有RPC调用的日志
func logInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// 记录请求开始
	log.Printf("开始处理RPC调用: %s", info.FullMethod)

	// 调用实际的处理器
	resp, err := handler(ctx, req)

	// 记录请求结束和耗时
	duration := time.Since(start)
	if err != nil {
		log.Printf("RPC调用失败: %s, 耗时: %v, 错误: %v", info.FullMethod, duration, err)
	} else {
		log.Printf("RPC调用成功: %s, 耗时: %v", info.FullMethod, duration)
	}

	return resp, err
}

func main() {
	// 设置监听地址和端口
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("启动监听失败: %v", err)
	}

	// 创建gRPC服务器，添加日志拦截器
	s := grpc.NewServer(
		grpc.UnaryInterceptor(logInterceptor),
	)

	// 注册图书服务
	bookServer := NewBookServer()
	pb.RegisterBookServiceServer(s, bookServer)

	// 打印启动信息
	log.Printf("图书管理服务启动成功，监听地址: %v", lis.Addr())
	log.Printf("服务提供以下功能:")
	log.Printf("- 创建图书 (CreateBook)")
	log.Printf("- 获取图书 (GetBook)")
	log.Printf("- 更新图书 (UpdateBook)")
	log.Printf("- 删除图书 (DeleteBook)")
	log.Printf("- 列出图书 (ListBooks)")
	log.Printf("- 按价格查询 (SearchBooksByPrice)")

	// 启动服务器
	if err := s.Serve(lis); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
