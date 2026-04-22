 # Step 1: 빌드 단계
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o webhook-server main.go

# Step 2: 실행 단계 (가벼운 이미지)
FROM alpine:latest  
WORKDIR /root/
RUN apk --no-cache add ca-certificates tzdata

# 빌드된 실행 파일 복사
COPY --from=builder /app/webhook-server .
# ⭐ 설정 파일(.env)을 이미지 내부로 복사 (이게 핵심!)
COPY .env .

EXPOSE 8080
CMD ["./webhook-server"]
