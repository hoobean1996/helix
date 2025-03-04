package main

import (
	"flag"
	"path/filepath"

	"helix.io/helix/package/blob"
	"helix.io/helix/package/rds"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// 创建 clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	mysqlManager := rds.NewMysqlManager(clientset)
	mysqlSpecification := &rds.MySQLSpecification{
		ProjectName:   "test-app",
		Env:           "test",       // test, staging, prod
		Password:      "test123456", // 测试环境的密码
		Database:      "testdb",
		RequestCPU:    "100m",  // 0.1 核心
		RequestMemory: "256Mi", // 256 MB
		LimitCPU:      "200m",  // 0.2 核心
		LimitMemory:   "512Mi", // 512 MB
	}
	mysqlManager.Deploy(mysqlSpecification)

	blobManager := blob.NewBlobManager(clientset)
	blobSpecification := &blob.BlobSpecification{
		ProjectName:   "test-app",
		Env:           "test",
		AccessKey:     "minioadmin",
		SecretKey:     "minioadmin",
		RequestCPU:    "100m",
		RequestMemory: "256Mi",
		LimitCPU:      "200m",
		LimitMemory:   "512Mi",
	}
	blobManager.Deploy(blobSpecification)
}
