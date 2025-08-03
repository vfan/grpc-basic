package main

import (
	"context"
	"fmt"
	"log"
	"time"

	// 导入生成的protobuf代码
	pb "grpc-basic-client/pb"

	// 导入gRPC相关包
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// BookClient 图书管理客户端
type BookClient struct {
	client pb.BookServiceClient
	conn   *grpc.ClientConn
}

// NewBookClient 创建新的图书客户端
func NewBookClient(serverAddr string) (*BookClient, error) {
	// 建立到服务器的连接
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("连接服务器失败: %v", err)
	}

	// 创建客户端
	client := pb.NewBookServiceClient(conn)

	return &BookClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭客户端连接
func (c *BookClient) Close() error {
	return c.conn.Close()
}

// CreateBook 创建图书
func (c *BookClient) CreateBook(title, author string, price float32, description string, publishYear int32) (string, error) {
	// 创建上下文，设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建图书信息
	book := &pb.Book{
		Title:       title,
		Author:      author,
		Price:       price,
		Description: description,
		PublishYear: publishYear,
	}

	// 发送创建图书请求
	resp, err := c.client.CreateBook(ctx, &pb.CreateBookRequest{Book: book})
	if err != nil {
		return "", fmt.Errorf("创建图书失败: %v", err)
	}

	log.Printf("✅ 图书创建成功，ID: %s", resp.Id)
	return resp.Id, nil
}

// GetBook 获取图书信息
func (c *BookClient) GetBook(bookID string) (*pb.Book, error) {
	// 创建上下文，设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 发送获取图书请求
	resp, err := c.client.GetBook(ctx, &pb.GetBookRequest{Id: bookID})
	if err != nil {
		return nil, fmt.Errorf("获取图书失败: %v", err)
	}

	log.Printf("✅ 成功获取图书: %s", resp.Book.Title)
	return resp.Book, nil
}

// UpdateBook 更新图书信息
func (c *BookClient) UpdateBook(bookID, title, author string, price float32, description string, publishYear int32) error {
	// 创建上下文，设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建更新的图书信息
	book := &pb.Book{
		Id:          bookID,
		Title:       title,
		Author:      author,
		Price:       price,
		Description: description,
		PublishYear: publishYear,
	}

	// 发送更新图书请求
	resp, err := c.client.UpdateBook(ctx, &pb.UpdateBookRequest{Book: book})
	if err != nil {
		return fmt.Errorf("更新图书失败: %v", err)
	}

	log.Printf("✅ 图书更新成功: %s", resp.Message)
	return nil
}

// DeleteBook 删除图书
func (c *BookClient) DeleteBook(bookID string) error {
	// 创建上下文，设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 发送删除图书请求
	resp, err := c.client.DeleteBook(ctx, &pb.DeleteBookRequest{Id: bookID})
	if err != nil {
		return fmt.Errorf("删除图书失败: %v", err)
	}

	log.Printf("✅ 图书删除成功: %s", resp.Message)
	return nil
}

// ListBooks 列出所有图书
func (c *BookClient) ListBooks(page, pageSize int32) ([]*pb.Book, int32, error) {
	// 创建上下文，设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 发送列出图书请求
	resp, err := c.client.ListBooks(ctx, &pb.ListBooksRequest{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("列出图书失败: %v", err)
	}

	log.Printf("✅ 成功列出图书，总数: %d, 当前页: %d", resp.Total, page)
	return resp.Books, resp.Total, nil
}

// SearchBooksByPrice 按价格区间查询图书
func (c *BookClient) SearchBooksByPrice(minPrice, maxPrice float32) ([]*pb.Book, error) {
	// 创建上下文，设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 发送按价格查询请求
	resp, err := c.client.SearchBooksByPrice(ctx, &pb.SearchBooksByPriceRequest{
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	})
	if err != nil {
		return nil, fmt.Errorf("按价格查询图书失败: %v", err)
	}

	log.Printf("✅ 按价格查询完成，找到 %d 本图书", len(resp.Books))
	return resp.Books, nil
}

// printBookInfo 打印图书信息
func printBookInfo(book *pb.Book) {
	fmt.Printf("📚 图书信息:\n")
	fmt.Printf("   ID: %s\n", book.Id)
	fmt.Printf("   标题: %s\n", book.Title)
	fmt.Printf("   作者: %s\n", book.Author)
	fmt.Printf("   价格: ¥%.2f\n", book.Price)
	fmt.Printf("   描述: %s\n", book.Description)
	fmt.Printf("   出版年份: %d\n", book.PublishYear)
	fmt.Println()
}

// printBookList 打印图书列表
func printBookList(books []*pb.Book) {
	if len(books) == 0 {
		fmt.Println("📚 暂无图书")
		return
	}

	fmt.Printf("📚 图书列表 (共 %d 本):\n", len(books))
	for i, book := range books {
		fmt.Printf("%d. %s - %s (¥%.2f)\n", i+1, book.Title, book.Author, book.Price)
	}
	fmt.Println()
}

func main() {
	// 创建客户端
	client, err := NewBookClient("localhost:50051")
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	defer client.Close()

	log.Println("🚀 开始演示图书管理服务...")
	log.Println("==================================================")

	// 演示1: 创建图书
	log.Println("📝 演示1: 创建图书")
	bookID1, err := client.CreateBook(
		"The Go Programming Language",
		"Alan A. A. Donovan",
		45.99,
		"Go语言的权威指南，适合初学者和有经验的开发者",
		2015,
	)
	if err != nil {
		log.Printf("❌ 创建图书失败: %v", err)
	}

	_, err = client.CreateBook(
		"Design Patterns",
		"Erich Gamma",
		39.99,
		"面向对象设计模式的经典著作",
		1994,
	)
	if err != nil {
		log.Printf("❌ 创建图书失败: %v", err)
	}

	bookID3, err := client.CreateBook(
		"Clean Code",
		"Robert C. Martin",
		29.99,
		"编写可维护代码的最佳实践",
		2008,
	)
	if err != nil {
		log.Printf("❌ 创建图书失败: %v", err)
	}

	// 演示2: 获取图书信息
	log.Println("📖 演示2: 获取图书信息")
	book, err := client.GetBook(bookID1)
	if err != nil {
		log.Printf("❌ 获取图书失败: %v", err)
	} else {
		printBookInfo(book)
	}

	// 演示3: 更新图书信息
	log.Println("✏️ 演示3: 更新图书信息")
	err = client.UpdateBook(
		bookID1,
		"The Go Programming Language (Updated)",
		"Alan A. A. Donovan",
		49.99,
		"Go语言的权威指南，适合初学者和有经验的开发者 (更新版)",
		2015,
	)
	if err != nil {
		log.Printf("❌ 更新图书失败: %v", err)
	}

	// 验证更新结果
	updatedBook, err := client.GetBook(bookID1)
	if err != nil {
		log.Printf("❌ 获取更新后的图书失败: %v", err)
	} else {
		printBookInfo(updatedBook)
	}

	// 演示4: 列出所有图书
	log.Println("📋 演示4: 列出所有图书")
	books, total, err := client.ListBooks(1, 10)
	if err != nil {
		log.Printf("❌ 列出图书失败: %v", err)
	} else {
		fmt.Printf("总共有 %d 本图书\n", total)
		printBookList(books)
	}

	// 演示5: 按价格区间查询
	log.Println("🔍 演示5: 按价格区间查询 (¥30-50)")
	priceBooks, err := client.SearchBooksByPrice(30, 50)
	if err != nil {
		log.Printf("❌ 按价格查询失败: %v", err)
	} else {
		printBookList(priceBooks)
	}

	// 演示6: 删除图书
	log.Println("🗑️ 演示6: 删除图书")
	err = client.DeleteBook(bookID3)
	if err != nil {
		log.Printf("❌ 删除图书失败: %v", err)
	}

	// 验证删除结果
	log.Println("📋 删除后的图书列表:")
	booksAfterDelete, _, err := client.ListBooks(1, 10)
	if err != nil {
		log.Printf("❌ 列出图书失败: %v", err)
	} else {
		printBookList(booksAfterDelete)
	}

	log.Println("🎉 演示完成!")
}
