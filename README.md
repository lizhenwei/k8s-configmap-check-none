# Kubernetes ConfigMap 空值检查工具

这个工具用于检查 Kubernetes 集群中指定 Namespace 下的 ConfigMap 是否存在空值。

## 功能特性

- 检查指定 Namespace 下的所有 ConfigMap
- 找出空值的 ConfigMap 及其对应的 Key
- 支持在集群内运行（InClusterConfig）
- 支持通过本地 kubeconfig 文件运行
- 支持通过命令行参数设置 Namespace
- 支持通过环境变量设置 Namespace
- 支持交叉编译为 Linux 64 位二进制文件

## 技术栈

- Go 1.24.2
- Kubernetes Client-go

## 构建步骤

### 1. 下载依赖

```bash
go mod tidy
```

### 2. 构建本地可执行文件

```bash
go build -o main .
```

### 3. 交叉编译 Linux 64 位二进制文件

```bash
# 在 Windows PowerShell 中运行
$env:GOOS="linux" ; $env:GOARCH="amd64" ; $env:CGO_ENABLED="0" ; go build -o k8s-configmap-check-none .

# 在 Linux/macOS 或 Windows CMD 中运行
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
go build -o k8s-configmap-check-none .
```

## 运行方式

### 1. 直接运行 Go 程序

```bash
# 使用默认 Namespace
go run main.go

# 通过环境变量设置 Namespace
export NAMESPACE=ess-cloud
go run main.go

# 通过命令行参数设置 Namespace
go run main.go --namespace=ess-cloud
```

### 2. 运行编译后的二进制文件

```bash
# 赋予执行权限
chmod +x k8s-configmap-check-none

# 运行二进制文件
./k8s-configmap-check-none

# 通过命令行参数设置 Namespace
./k8s-configmap-check-none --namespace=default

# 通过环境变量设置 Namespace
export NAMESPACE=ess-cloud
./k8s-configmap-check-none

# 使用--ignore参数忽略包含特定关键字的 Key
./k8s-configmap-check-none --namespace=default --ignore=TEST
```

## 命令行参数

| 参数 | 描述 | 默认值 |
|------|------|--------|
| `--namespace` | 要检查的 Kubernetes Namespace | default |
| `--ignore` | 要忽略的关键字，包含此关键字的 Key 将不会被检查 | 空字符串（不忽略任何 Key） |

## 环境变量说明

| 环境变量 | 描述 | 默认值 |
|---------|------|--------|
| NAMESPACE | 要检查的 Kubernetes Namespace | default |

## 优先级

设置 Namespace 的优先级从高到低为：
1. 命令行参数 `--namespace`
2. 环境变量 `NAMESPACE`
3. 默认值 `default`

## 示例输出

### 发现空值的情况

```
在 namespace 'default' 中发现空值的 ConfigMap:
- ConfigMap: test-cm, 空的 key: empty-key
- ConfigMap: another-cm, 空的 key: ALL
```

### 没有发现空值的情况

```
在 namespace 'default' 中没有发现空值的 ConfigMap.
```

## 注意事项

1. 当在集群内运行时，程序会自动使用 InClusterConfig，无需 kubeconfig 文件
2. 当在集群外运行时，程序会使用默认 kubeconfig 路径 `~/.kube/config`
3. 可以通过设置 `KUBECONFIG` 环境变量指定自定义 kubeconfig 路径

## 项目结构

```
.
├── go.mod           # Go 模块依赖
├── main.go          # 主程序代码
└── readme.md        # 项目说明文档
```
