
### 编译启动
```code
mkdir -p /apps/repo
cd /apps/repo
git clone  https://github.com/lixiang4u/docker-lnmp.git

cd docker-lnmp/src/
docker build --no-cache --tag docker-lnmp:latest .
docker run --rm -v /apps/repo/docker-lnmp:/home docker-lnmp:latest
```

### 说明
- 容器定义在 src/api/model/docker-compose.go
- 虚拟主机配置变更后需要重建docker项目
- 修改nginx/php等配置需要重构项目，不要再容器中直接更改，避免重构后丢失