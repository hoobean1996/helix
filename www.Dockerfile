# 多阶段构建：第一阶段 - 编译React应用
FROM node:20-alpine AS build

# 设置工作目录
WORKDIR /app

# 复制项目文件
COPY www/package.json www/package-lock.json ./
# 安装依赖
RUN npm ci

# 复制源代码
COPY www ./

# 构建应用
RUN npm run build

# 多阶段构建：第二阶段 - 设置Nginx服务器
FROM nginx:alpine

# 从构建阶段复制编译好的文件到Nginx服务目录
COPY --from=build /app/build /usr/share/nginx/html

# 暴露80端口
EXPOSE 80

# 启动Nginx
CMD ["nginx", "-g", "daemon off;"]