package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	"github.com/gin-gonic/gin"
)

// 定义 push 事件的结构
type PushEvent struct {
	Ref        string `json:"ref"`
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Name string `json:"name"`
		} `json:"owner"`
	} `json:"repository"`
	Commits []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"commits"`
}

func main() {
	r := gin.Default()

	r.POST("/api/webhook", func(c *gin.Context) {
		// 读取请求体
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to read request body"})
			return
		}

		// 解析事件类型
		event := c.GetHeader("X-GitHub-Event")
		if event == "push" {
			var pushEvent PushEvent
			if err := json.Unmarshal(body, &pushEvent); err != nil {
				c.JSON(400, gin.H{"error": "Failed to parse push event"})
				return
			}

			// 设置工作目录
			repoPath := "/Users/binhe/saas/helix/example/app"

			// 执行 git pull
			cmd := exec.Command("git", "pull")
			cmd.Dir = repoPath
			var pullOutput bytes.Buffer
			cmd.Stdout = &pullOutput
			cmd.Stderr = &pullOutput
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error executing git pull: %v\n", err)
				fmt.Printf("Command output: %s\n", pullOutput.String())
				c.JSON(500, gin.H{"error": "Failed to pull repository"})
				return
			}

			// Docker build
			cmd1 := exec.Command("docker", "build", "-t", "app:latest", ".")
			cmd1.Dir = repoPath
			var buildOutput bytes.Buffer
			cmd1.Stdout = &buildOutput
			cmd1.Stderr = &buildOutput
			if err := cmd1.Run(); err != nil {
				fmt.Printf("Error executing docker build: %v\n", err)
				fmt.Printf("Command output: %s\n", buildOutput.String())
				c.JSON(500, gin.H{"error": "Failed to build Docker image"})
				return
			}

			// k3d image import
			cmd2 := exec.Command("k3d", "image", "import", "app:latest", "-c", "my-cluster")
			cmd2.Dir = repoPath
			var importOutput bytes.Buffer
			cmd2.Stdout = &importOutput
			cmd2.Stderr = &importOutput
			if err := cmd2.Run(); err != nil {
				fmt.Printf("Error executing k3d image import: %v\n", err)
				fmt.Printf("Command output: %s\n", importOutput.String())
				c.JSON(500, gin.H{"error": "Failed to import image to k3d"})
				return
			}

			// kubectl apply
			cmd3 := exec.Command("kubectl", "apply", "-f", "deploy.yml")
			cmd3.Dir = repoPath
			var applyOutput bytes.Buffer
			cmd3.Stdout = &applyOutput
			cmd3.Stderr = &applyOutput
			if err := cmd3.Run(); err != nil {
				fmt.Printf("Error executing kubectl apply: %v\n", err)
				fmt.Printf("Command output: %s\n", applyOutput.String())
				c.JSON(500, gin.H{"error": "Failed to apply deployment"})
				return
			}

			fmt.Printf("Deployment successful\n")
		}

		c.JSON(200, gin.H{"status": "success"})
	})

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("Run gin engine failed, err=%s\n", err.Error())
	}
}
