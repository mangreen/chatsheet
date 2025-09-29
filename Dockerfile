# chatsheet/Dockerfile

# ==================================
# 階段 1: Golang 後端編譯 (Builder Stage)
# ==================================
FROM golang:1.23.12-alpine AS builder

# 設定 Go 模組代理以加速下載
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct

# 安裝 Git, CGO 相關的工具
RUN apk add --no-cache git gcc musl-dev

# 設定工作目錄
WORKDIR /app

# 複製 Go 模組文件
COPY go.mod .
COPY go.sum .

# 下載所有依賴，並緩存 (如果 go.mod/go.sum 未改變)
RUN go mod download

# 複製所有專案原始碼
COPY . .

# 編譯 Go 應用程式
# CGO_ENABLED=0 是為了建立一個完全靜態連結的執行檔
# -o myapp 是輸出檔名
# ./cmd/myapp 是 main 函式所在的路徑
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /myapp ./cmd/myapp

# ==================================
# 階段 2: Svelte 前端打包 (Frontend Stage)
# ==================================
FROM node:20-alpine AS frontend

WORKDIR /app/chatsheet/web/myapp

# 複製前端配置和依賴文件
COPY ./web/myapp/package.json .
COPY ./web/myapp/package-lock.json .
COPY ./web/myapp/vite.config.js .
COPY ./web/myapp/svelte.config.js .

# 安裝前端依賴
RUN npm install

# 複製 Svelte 原始碼
COPY ./web/myapp/src ./src
COPY ./web/myapp/index.html .

# 執行前端打包，輸出到 /app/chatsheet/web/myapp/dist
RUN npm run build


# ==================================
# 階段 3: 最終映像檔 (Final Stage)
# ==================================
FROM alpine:latest

# 設定時區 (確保應用程式日誌時間正確)
ENV TZ=Asia/Taipei
RUN apk add --no-cache tzdata

# 設定工作目錄
WORKDIR /root/

# 複製 Go 執行檔
# myapp 是後端伺服器執行檔
COPY --from=builder /myapp /root/myapp
COPY ./config /root/config

# 複製前端靜態檔案
# chatsheet/web/myapp/dist 是 main.go 中 r.Static 和 r.NoRoute 所需的路徑
COPY --from=frontend /app/chatsheet/web/myapp/dist /root/web/myapp/dist

# 暴露服務埠
EXPOSE 8080

# 設定容器啟動命令
# 運行後端伺服器執行檔
CMD ["/root/myapp"]