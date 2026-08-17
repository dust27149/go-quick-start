# 示例程序
## 简介
目录下的示例程序是基于`go-quick-start`框架的示例应用，展示了如何使用该框架进行快速开发。每个示例程序都包含了完整的代码和配置文件，便于开发者参考和学习。

目录结构
```text
samples/
├── internal
|   ├── components/     # 组件示例
|   |  ├── https/       # HTTP组件示例
|   |  ├── logger/      # 日志组件示例
|   |  └── writer/      # 写入组件示例
│   └── utils/          # 工具函数示例
|      ├── compress/    # 压缩工具示例
|      ├── config/      # 配置工具示例
|      └── merge/       # 合并工具示例
├── logs/               # 日志目录
├── modules/            # 模块示例
|  └── api/             # API模块示例
├── .dockerignore       # Docker忽略文件
├── .env                # Docker环境变量配置文件
├── config.json         # 配置文件
├── docker-compose.yml  # Docker Compose配置文件
├── Dockerfile          # Dockerfile文件
├── go.mod              # Go模块文件
├── main.go             # 示例程序入口
└── README.md           # 示例程序说明文档
```

项目配置文件见`config.json`，可以根据需要进行修改。示例程序使用了`go-quick-start`框架的默认配置，开发者可以根据实际需求进行调整。

项目在`internal/components/http`目录下提供了两个http接口，程序启动后可以通过访问以下接口进行测试：

curl示例：
```bash
curl -X GET "http://localhost:8080/health?a=1&b=2"
curl -X POST "http://localhost:8080/health" \
  -H "Content-Type: application/json" \
  -d '{"a":1,"b":2}'
```

## Docker交付
### 镜像导出/导入
修改`.env`文件中的`HOST_PORT`和`CONTAINER_PORT`变量，可以自定义主机端口和容器端口的映射关系。

在项目根目录下执行以下命令，可以将构建好的镜像导出为`go-quick-start-samples.tar`文件，或者从该文件导入镜像。
```bash
docker-compose build # 构建镜像
docker save go-quick-start-samples:latest -o go-quick-start-samples.tar # 导出镜像
docker load -i go-quick-start-samples.tar # 导入镜像
```

### 运行
在项目根目录下执行以下命令，可以启动示例程序。
```bash
docker-compose up -d
```

### 停止
在项目根目录下执行以下命令，可以停止示例程序。
```bash
docker-compose down
```