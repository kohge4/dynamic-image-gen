FROM golang:1.23.4 as build-lambda
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main /go/src/app/cmd/lambda/main.go

FROM golang:1.23.4 as build-local
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main /go/src/app/cmd/local/main.go

FROM chromedp/headless-shell:latest as lambda
WORKDIR /app
RUN mkdir /app/tmp
# 日本語フォントとemoji用
RUN apt-get update && apt-get install -y \
    dumb-init \
    fonts-noto-cjk \
    fonts-noto-color-emoji \
    fontconfig \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean \
    && rm -rf /var/cache/apt/archives/*
# フォントキャッシュの再構築
RUN fc-cache -fv
ENTRYPOINT ["dumb-init", "--"]
COPY --from=build-lambda /go/src/app/main .
COPY --from=build-lambda /go/src/app/static ./static
CMD ["/app/main"]

FROM chromedp/headless-shell:latest as local
WORKDIR /app
# 日本語フォントとemoji用
RUN apt-get update && apt-get install -y \
    dumb-init \
    fonts-noto-cjk \
    fonts-noto-color-emoji \
    fontconfig \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean \
    && rm -rf /var/cache/apt/archives/*
# フォントキャッシュの再構築
RUN fc-cache -fv
ENTRYPOINT ["dumb-init", "--"]
COPY --from=build-local /go/src/app/main .
COPY --from=build-local /go/src/app/static ./static
CMD ["/app/main"]