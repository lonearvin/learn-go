
# go 教程
# docker命令
## sudo docker inspect <contain_id>
这是一个 Docker 命令，用于获取指定容器（这里是 zjt-mysql）的详细信息（以 JSON 格式输出）。这个 JSON 包含容器的所有元信息，例如状态、网络、挂载、端口等。
```bash
sudo docker inspect zjt-mysql
```
# 创建一个redis
在 Docker 中配置 Redis 非常简单，你可以通过以下几种方式运行和配置 Redis。

---

## 1. **快速启动 Redis 容器**

### 1.1 拉取 Redis 镜像
首先从 Docker Hub 拉取最新的 Redis 镜像：
```bash
docker pull redis:latest
```

### 1.2 启动 Redis 容器
运行以下命令启动 Redis 容器：
```bash
docker run --name my-redis -d redis
```

- `--name my-redis`：指定容器名称为 `my-redis`。
- `-d`：让容器以后台方式运行。

Redis 会运行在容器的默认端口 `6379`。

---

## 2. **使用端口映射访问 Redis**

如果希望从主机访问 Redis，可以映射 Redis 的端口到主机：
```bash
docker run --name my-redis -d -p 6379:6379 redis
```

- `-p 6379:6379`：将容器的 `6379` 端口映射到主机的 `6379` 端口。

现在可以通过以下方式连接 Redis：
```bash
redis-cli -h 127.0.0.1 -p 6379
```

---

## 3. **挂载本地配置文件**

如果需要自定义 Redis 的配置，可以挂载本地的配置文件到容器中。

### 3.1 创建自定义 Redis 配置文件
在本地创建一个 `redis.conf` 文件，例如：
```bash
mkdir ~/redis-config
nano ~/redis-config/redis.conf
```

添加你需要的 Redis 配置，例如：
```conf
bind 0.0.0.0
protected-mode no
```

- `bind 0.0.0.0`：允许所有 IP 地址访问（默认只允许本地访问）。
- `protected-mode no`：关闭保护模式，允许远程连接（仅供测试用，生产环境需小心）。

### 3.2 挂载配置文件到容器
启动 Redis 时，将本地的配置文件挂载到容器中：
```bash
docker run --name my-redis -d -p 6379:6379 -v ~/redis-config/redis.conf:/usr/local/etc/redis/redis.conf redis redis-server /usr/local/etc/redis/redis.conf
```

- `-v ~/redis-config/redis.conf:/usr/local/etc/redis/redis.conf`：将本地的 `redis.conf` 文件挂载到容器中。
- `redis redis-server /usr/local/etc/redis/redis.conf`：指定 Redis 使用自定义配置文件启动。

---

## 4. **使用 Docker Compose 配置 Redis**

如果你使用 Docker Compose 管理多个容器，可以使用以下配置：

### 4.1 创建 `docker-compose.yml`
在当前目录下创建一个 `docker-compose.yml` 文件：
```yaml
version: '3.8'
services:
  redis:
    image: redis:latest
    container_name: my-redis
    ports:
      - "6379:6379"
    volumes:
      - ./redis-data:/data
      - ./redis.conf:/usr/local/etc/redis/redis.conf
    command: ["redis-server", "/usr/local/etc/redis/redis.conf"]
```

- `ports`：映射容器和主机的端口。
- `volumes`：
  - 挂载本地目录 `./redis-data` 用于持久化数据。
  - 挂载自定义的 `redis.conf` 文件。
- `command`：指定使用配置文件启动 Redis。

### 4.2 启动服务
运行以下命令启动 Redis：
```bash
docker-compose up -d
```

---

## 5. **Redis 数据持久化**

Redis 默认使用内存存储数据，但支持持久化到磁盘。你可以使用 Docker 挂载卷的方式来实现数据持久化。

### 5.1 挂载数据目录
运行容器时挂载本地目录以保存 Redis 数据：
```bash
docker run --name my-redis -d -p 6379:6379 -v ~/redis-data:/data redis
```

- `-v ~/redis-data:/data`：将 Redis 的 `/data` 目录挂载到主机的 `~/redis-data`。

Redis 会将持久化文件（如 `dump.rdb`）保存在 `~/redis-data` 中。

### 5.2 开启持久化模式
确保 `redis.conf` 中启用了持久化模式（默认已启用）：
```conf
save 900 1
save 300 10
save 60 10000
dir /data
```

---

## 6. **测试连接 Redis**

如果 Redis 容器运行正常，可以通过以下方式测试连接：

### 6.1 使用 Redis CLI 连接
在主机上运行：
```bash
redis-cli -h 127.0.0.1 -p 6379
```

### 6.2 进入容器内部测试
你还可以进入 Redis 容器内部测试：
```bash
docker exec -it my-redis redis-cli
```

---

## 7. **安全性配置**

在生产环境中，建议开启以下安全配置：

1. **设置密码**
   在 `redis.conf` 中添加：
   ```conf
   requirepass yourpassword
   ```

   然后重启 Redis 容器，并通过以下命令连接：
   ```bash
   redis-cli -h 127.0.0.1 -p 6379 -a yourpassword
   ```

2. **限制绑定地址**
   配置 `bind`，例如仅允许特定 IP 访问：
   ```conf
   bind 127.0.0.1
   ```

3. **关闭无保护模式**
   确保生产环境中启用保护模式：
   ```conf
   protected-mode yes
   ```

---

通过这些步骤，你可以在 Docker 中轻松地运行、配置和管理 Redis。
## redis分布式锁
### 实现redlock算法-多redis实现分布式锁
#### 使用使用 Docker Compose 创建多个 Redis 实例
docker compose 是用于定义和运行多容器 Docker 应用程序的工具。通过 Compose，您可以使用 YML 文件来配置应用程序需要的所有服务。然后，使用一个命令，就可以从 YML 文件配置中创建并启动所有服务。
```yaml
version: '3.8'
services:
  redis1:
    image: redis:latest
    container_name: redis1
    ports:
      - "6379:6379"
  redis2:
    image: redis:latest
    container_name: redis2
    ports:
      - "6380:6379"
  redis3:
    image: redis:latest
    container_name: redis3
    ports:
      - "6381:6379"
```
启动：
    将上述 docker-compose.yml 保存到当前目录。
    在终端运行以下命令：
    ```bash
        docker-compose up -d
    ```














