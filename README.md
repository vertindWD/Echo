# Echo - 高性能社区后端系统

Echo 是一款基于 Go 语言构建的高性能社区论坛后端服务。项目采用标准的 Controller-Logic-DAO 三层架构，并针对社交场景下的“实时排行”与“高频投票”进行了专项性能优化。

---

## 🚀 核心技术亮点

### 1. 高性能排行榜方案 (Redis ZSet)
* **技术方案**：利用 Redis Sorted Set 存储帖子分值，将分页排序复杂度降至 O(logN)。
* **应用场景**：支持“最新发布”与“最高得分”双维度实时切换，避免了 MySQL 磁盘排序的性能瓶颈。

### 2. 批量状态聚合优化 (Redis Pipeline)
* **技术方案**：在列表页展示时，通过 Redis Pipeline 将多次查询指令合并为一次批量操作。
* **优化效果**：彻底消除 N+1 网络往返问题，高并发下接口吞吐量提升约 40%。

### 3. 工程化与安全实践
* **鉴权安全**：基于 JWT 的无状态认证机制，配合中间件实现灵活的权限控制。
* **唯一ID**：集成 Snowflake 雪花算法生成分布式唯一 ID，确保主键有序性。
* **自动文档**：集成 Swagger (swaggo) 实现接口文档的自动化构建。

---

## 📊 性能压测报告 (Benchmark)

环境：i9-13900HX (24C/32T) | 32GB RAM | WSL2

* **测试接口**：/api/v1/post/list (带 JWT 鉴权)
* **吞吐量 (QPS)**：82,592.26 req/sec
* **平均延迟 (Latency)**：1.27 ms
* **错误率**：0.00%

> 压测工具：wrk (4 threads, 100 connections)

---

## 📂 项目结构说明

* `conf/` : 配置文件管理
* `controller/` : 接口入口与参数校验
* `logic/` : 核心业务逻辑处理
* `dao/` : 数据库与缓存交互 (MySQL/Redis)
* `models/` : 结构体与 DTO 定义
* `pkg/` : JWT、雪花算法等公共组件

---

## 🛠️ 快速开始

1. **配置环境**：将 `conf/config.example.yaml` 复制为 `conf/config.yaml` 并填入数据库连接信息。
2. **启动服务**：执行 `go run main.go`。
3. **查看文档**：访问 `http://127.0.0.1:8080/swagger/index.html`。
