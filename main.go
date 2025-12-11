package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 命令行参数
	var namespace string
	var ignore string
	flag.StringVar(&namespace, "namespace", "", "要检查的 Kubernetes Namespace")
	flag.StringVar(&ignore, "ignore", "", "要忽略的关键字，包含此关键字的 Key 将不会被检查")
	flag.Parse()

	// 优先级：命令行参数 > 环境变量 > 默认值
	if namespace == "" {
		namespace = os.Getenv("NAMESPACE")
		if namespace == "" {
			namespace = "default"
		}
	}

	// Kubernetes 配置
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = os.ExpandEnv("$HOME/.kube/config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	cms, err := clientset.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// 直接检查空值，不需要复杂的模式匹配
	// 改为使用简单的字符串包含检查，查找各种形式的空值
	checkEmpty := func(line string) bool {
		// 移除空格以便更准确地检查
		line = strings.ReplaceAll(line, " ", "")
		
		// 检查各种空值形式
		return strings.Contains(line, ":''") ||
			   strings.Contains(line, ":\"\"") ||
			   strings.Contains(line, ":``") ||
			   strings.Contains(line, "=''") ||
			   strings.Contains(line, "='\"\"") ||
			   strings.Contains(line, "='``") ||
			   strings.Contains(line, ":'',") ||
			   strings.Contains(line, ":\"\",") ||
			   strings.Contains(line, ":``,") ||
			   strings.Contains(line, "='',") ||
			   strings.Contains(line, "='\"\",") ||
			   strings.Contains(line, "='``,")
	}

	foundEmpty := false
	fmt.Printf("在 namespace '%s' 中发现空值的 ConfigMap:\n", namespace)

	for _, cm := range cms.Items {
		// 定义ANSI颜色转义序列
		blueColor := "\033[34m"
		resetColor := "\033[0m"

		if cm.Data == nil || len(cm.Data) == 0 {
			fmt.Printf("- %sConfigMap: %s, Key: %sALL (整个 ConfigMap 为空)\n", blueColor, cm.Name, resetColor)
			foundEmpty = true
			continue
		}

		for key, value := range cm.Data {
			// 跳过忽略关键字
			if ignore != "" && strings.Contains(key, ignore) {
				continue
			}

			// 检查整个 value 是否为空
			trimmed := strings.TrimSpace(value)
			if trimmed == "" || trimmed == "''" || trimmed == "\"\"" || trimmed == "``" {
				fmt.Printf("- %sConfigMap: %s, Key: %s%s (值为空)\n", blueColor, cm.Name, resetColor, key)
				foundEmpty = true
				continue
			}

			// 逐行检查空值
			lines := strings.Split(value, "\n")
			for lineNum, line := range lines {
				originalLine := line // 保存原始行
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}

				// 检查当前行是否包含空值
				if checkEmpty(line) {
					// 移除行尾逗号（如果有）以便显示更清晰
					cleanLine := strings.TrimSuffix(originalLine, ",")
					fmt.Printf("- %sConfigMap: %s, Key: %s%s, Line %d: %s\n", blueColor, cm.Name, resetColor, key, lineNum+1, cleanLine)
					foundEmpty = true
				}
			}
		}
	}

	if !foundEmpty {
		fmt.Printf("在 namespace '%s' 中没有发现空值的 ConfigMap.\n", namespace)
	}
}
