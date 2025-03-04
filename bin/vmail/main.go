package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/mail"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"helix.io/helix/ent"
)

const (
	// 服务器配置
	listenAddr = "0.0.0.0:25" // 监听地址和端口
)

// 处理SMTP会话
func handleSession(conn net.Conn, client *ent.Client) {
	defer conn.Close()

	// 设置读写超时
	conn.SetDeadline(time.Now().Add(5 * time.Minute))

	// 打印连接信息
	fmt.Printf("【新连接】来自: %s\n", conn.RemoteAddr())

	// 发送欢迎信息
	conn.Write([]byte("220 邮件服务器就绪\r\n"))

	// 读取邮件数据
	reader := bufio.NewReader(conn)

	// 解析SMTP会话并收集邮件信息
	var (
		from       string
		to         []string
		data       strings.Builder
		inDataMode bool
	)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("【错误】读取数据: %v\n", err)
			}
			return
		}

		line = strings.TrimRight(line, "\r\n")

		// 如果处于DATA模式，收集邮件内容
		if inDataMode {
			// 检查数据结束标记
			if line == "." {
				inDataMode = false
				conn.Write([]byte("250 邮件已接收\r\n"))

				// 处理收到的邮件
				processEmail(from, to, data.String(), client)

				// 重置数据，准备接收下一封邮件
				from = ""
				to = to[:0]
				data.Reset()
			} else {
				// 收集数据
				data.WriteString(line)
				data.WriteString("\n")
			}
			continue
		}

		// 处理SMTP命令
		if strings.HasPrefix(line, "EHLO") || strings.HasPrefix(line, "HELO") {
			conn.Write([]byte("250 欢迎\r\n"))
		} else if strings.HasPrefix(line, "MAIL FROM:") {
			from = strings.TrimPrefix(line, "MAIL FROM:")
			from = strings.Trim(from, "<>")
			conn.Write([]byte("250 发件人已接受\r\n"))
		} else if strings.HasPrefix(line, "RCPT TO:") {
			recipient := strings.TrimPrefix(line, "RCPT TO:")
			recipient = strings.Trim(recipient, "<>")
			to = append(to, recipient)
			conn.Write([]byte("250 收件人已接受\r\n"))
		} else if line == "DATA" {
			inDataMode = true
			conn.Write([]byte("354 请输入邮件内容，以\".\"结束\r\n"))
		} else if line == "QUIT" {
			conn.Write([]byte("221 再见\r\n"))
			return
		} else {
			conn.Write([]byte("500 未知命令\r\n"))
		}
	}
}

// 处理接收到的邮件
func processEmail(from string, to []string, data string, client *ent.Client) {
	fmt.Println("\n==========================================================")
	fmt.Println("         收到一封新邮件         ")
	fmt.Println("==========================================================")
	fmt.Printf("【时间】: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("【发件人】: %s\n", from)
	fmt.Printf("【收件人】: %s\n", strings.Join(to, ", "))

	// 解析邮件内容
	msg, err := mail.ReadMessage(strings.NewReader(data))
	if err != nil {
		fmt.Printf("【错误】解析邮件: %v\n", err)
		fmt.Println("【原始内容】:")
		fmt.Println("----------------------------------------------------------")
		fmt.Println(data)
		fmt.Println("----------------------------------------------------------")
		return
	}

	// 打印邮件头信息
	fmt.Printf("【主题】: %s\n", msg.Header.Get("Subject"))
	fmt.Printf("【日期】: %s\n", msg.Header.Get("Date"))

	// 读取并打印邮件正文
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		fmt.Printf("【错误】读取邮件正文: %v\n", err)
		return
	}

	fmt.Println("【邮件正文】:")
	fmt.Println("----------------------------------------------------------")
	fmt.Println(string(body))
	fmt.Println("----------------------------------------------------------")
	fmt.Println()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client.EntEmail.Create().SetFrom(from).SetTo(to).SetDate(msg.Header.Get("Date")).SetTopic(msg.Header.Get("Subject")).SetBody(string(body)).Exec(ctx)
}

func main() {
	client, err := ent.Open("mysql", "user:password@tcp(mysql:3306)/mydb?parseTime=True")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// 监听TCP端口
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("无法监听端口 %s: %v\n", listenAddr, err)
		return
	}
	defer listener.Close()
	fmt.Printf("邮件服务器正在监听 %s\n", listenAddr)
	// 接受连接并处理
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("接受连接失败: %v\n", err)
			continue
		}

		// 为每个连接创建一个新的goroutine
		go handleSession(conn, client)
	}
}
