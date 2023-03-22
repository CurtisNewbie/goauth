FROM golang:1.18
LABEL author="yongjie.zhuang"
LABEL descrption="Goauth - Resource mangement & authorization service"

ENV TZ=Asia/Shanghai

WORKDIR /usr/src/goauth

# for golang env
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# dependencies
COPY go.mod .
COPY go.sum .

RUN go mod download

# build executable
COPY . .

RUN go build -o ./main

# script (for io redirection stuff)
# COPY run.sh ./ 
RUN chmod +x run.sh

CMD ["sh", "run.sh"]
