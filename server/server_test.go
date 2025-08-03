package main

import (
	"context"
	"testing"

	// 导入生成的protobuf代码
	pb "grpc-basic/pb/bookstore"
)

// TestCreateBook 测试创建图书功能
func TestCreateBook(t *testing.T) {
	// 创建服务器实例
	server := NewBookServer()

	// 创建测试图书
	book := &pb.Book{
		Title:       "测试图书",
		Author:      "测试作者",
		Price:       29.99,
		Description: "这是一本测试图书",
		PublishYear: 2023,
	}

	// 创建请求
	req := &pb.CreateBookRequest{Book: book}

	// 调用创建图书方法
	resp, err := server.CreateBook(context.Background(), req)

	// 验证结果
	if err != nil {
		t.Fatalf("创建图书失败: %v", err)
	}

	if resp.Id == "" {
		t.Error("返回的图书ID为空")
	}

	if resp.Message != "图书创建成功" {
		t.Errorf("期望消息为'图书创建成功'，实际为: %s", resp.Message)
	}

	// 验证图书是否已存储
	if storedBook, exists := server.books[resp.Id]; !exists {
		t.Error("图书未正确存储")
	} else if storedBook.Title != book.Title {
		t.Errorf("存储的图书标题不匹配，期望: %s, 实际: %s", book.Title, storedBook.Title)
	}
}

// TestGetBook 测试获取图书功能
func TestGetBook(t *testing.T) {
	// 创建服务器实例
	server := NewBookServer()

	// 先创建一本图书
	book := &pb.Book{
		Title:       "测试图书",
		Author:      "测试作者",
		Price:       29.99,
		Description: "这是一本测试图书",
		PublishYear: 2023,
	}

	createReq := &pb.CreateBookRequest{Book: book}
	createResp, err := server.CreateBook(context.Background(), createReq)
	if err != nil {
		t.Fatalf("创建图书失败: %v", err)
	}

	// 获取图书
	getReq := &pb.GetBookRequest{Id: createResp.Id}
	getResp, err := server.GetBook(context.Background(), getReq)

	// 验证结果
	if err != nil {
		t.Fatalf("获取图书失败: %v", err)
	}

	if getResp.Book.Title != book.Title {
		t.Errorf("图书标题不匹配，期望: %s, 实际: %s", book.Title, getResp.Book.Title)
	}

	if getResp.Book.Author != book.Author {
		t.Errorf("作者不匹配，期望: %s, 实际: %s", book.Author, getResp.Book.Author)
	}
}

// TestUpdateBook 测试更新图书功能
func TestUpdateBook(t *testing.T) {
	// 创建服务器实例
	server := NewBookServer()

	// 先创建一本图书
	book := &pb.Book{
		Title:       "原始图书",
		Author:      "原始作者",
		Price:       29.99,
		Description: "原始描述",
		PublishYear: 2023,
	}

	createReq := &pb.CreateBookRequest{Book: book}
	createResp, err := server.CreateBook(context.Background(), createReq)
	if err != nil {
		t.Fatalf("创建图书失败: %v", err)
	}

	// 更新图书
	updatedBook := &pb.Book{
		Id:          createResp.Id,
		Title:       "更新后的图书",
		Author:      "更新后的作者",
		Price:       39.99,
		Description: "更新后的描述",
		PublishYear: 2024,
	}

	updateReq := &pb.UpdateBookRequest{Book: updatedBook}
	updateResp, err := server.UpdateBook(context.Background(), updateReq)

	// 验证更新结果
	if err != nil {
		t.Fatalf("更新图书失败: %v", err)
	}

	if updateResp.Message != "图书更新成功" {
		t.Errorf("期望消息为'图书更新成功'，实际为: %s", updateResp.Message)
	}

	// 验证图书是否已更新
	if storedBook, exists := server.books[createResp.Id]; !exists {
		t.Error("图书不存在")
	} else if storedBook.Title != updatedBook.Title {
		t.Errorf("图书标题未正确更新，期望: %s, 实际: %s", updatedBook.Title, storedBook.Title)
	}
}

// TestDeleteBook 测试删除图书功能
func TestDeleteBook(t *testing.T) {
	// 创建服务器实例
	server := NewBookServer()

	// 先创建一本图书
	book := &pb.Book{
		Title:       "要删除的图书",
		Author:      "作者",
		Price:       29.99,
		Description: "描述",
		PublishYear: 2023,
	}

	createReq := &pb.CreateBookRequest{Book: book}
	createResp, err := server.CreateBook(context.Background(), createReq)
	if err != nil {
		t.Fatalf("创建图书失败: %v", err)
	}

	// 删除图书
	deleteReq := &pb.DeleteBookRequest{Id: createResp.Id}
	deleteResp, err := server.DeleteBook(context.Background(), deleteReq)

	// 验证删除结果
	if err != nil {
		t.Fatalf("删除图书失败: %v", err)
	}

	if deleteResp.Message != "图书删除成功" {
		t.Errorf("期望消息为'图书删除成功'，实际为: %s", deleteResp.Message)
	}

	// 验证图书是否已删除
	if _, exists := server.books[createResp.Id]; exists {
		t.Error("图书未被正确删除")
	}
}

// TestListBooks 测试列出图书功能
func TestListBooks(t *testing.T) {
	// 创建服务器实例
	server := NewBookServer()

	// 创建多本图书
	books := []*pb.Book{
		{Title: "图书1", Author: "作者1", Price: 29.99, Description: "描述1", PublishYear: 2023},
		{Title: "图书2", Author: "作者2", Price: 39.99, Description: "描述2", PublishYear: 2024},
		{Title: "图书3", Author: "作者3", Price: 49.99, Description: "描述3", PublishYear: 2025},
	}

	// 创建图书
	for _, book := range books {
		req := &pb.CreateBookRequest{Book: book}
		_, err := server.CreateBook(context.Background(), req)
		if err != nil {
			t.Fatalf("创建图书失败: %v", err)
		}
	}

	// 列出图书
	listReq := &pb.ListBooksRequest{Page: 1, PageSize: 10}
	listResp, err := server.ListBooks(context.Background(), listReq)

	// 验证结果
	if err != nil {
		t.Fatalf("列出图书失败: %v", err)
	}

	if listResp.Total != 3 {
		t.Errorf("期望总数为3，实际为: %d", listResp.Total)
	}

	if len(listResp.Books) != 3 {
		t.Errorf("期望图书数量为3，实际为: %d", len(listResp.Books))
	}
}

// TestSearchBooksByPrice 测试按价格查询图书功能
func TestSearchBooksByPrice(t *testing.T) {
	// 创建服务器实例
	server := NewBookServer()

	// 创建不同价格的图书
	books := []*pb.Book{
		{Title: "便宜图书", Author: "作者1", Price: 19.99, Description: "描述1", PublishYear: 2023},
		{Title: "中等图书", Author: "作者2", Price: 39.99, Description: "描述2", PublishYear: 2024},
		{Title: "昂贵图书", Author: "作者3", Price: 59.99, Description: "描述3", PublishYear: 2025},
	}

	// 创建图书
	for _, book := range books {
		req := &pb.CreateBookRequest{Book: book}
		_, err := server.CreateBook(context.Background(), req)
		if err != nil {
			t.Fatalf("创建图书失败: %v", err)
		}
	}

	// 按价格区间查询
	searchReq := &pb.SearchBooksByPriceRequest{MinPrice: 30, MaxPrice: 50}
	searchResp, err := server.SearchBooksByPrice(context.Background(), searchReq)

	// 验证结果
	if err != nil {
		t.Fatalf("按价格查询图书失败: %v", err)
	}

	// 应该只找到一本中等价格的图书
	if len(searchResp.Books) != 1 {
		t.Errorf("期望找到1本图书，实际找到: %d", len(searchResp.Books))
	}

	if searchResp.Books[0].Title != "中等图书" {
		t.Errorf("期望图书标题为'中等图书'，实际为: %s", searchResp.Books[0].Title)
	}
}
