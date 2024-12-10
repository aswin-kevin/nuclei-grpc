FROM golang:1.23

WORKDIR /go/src/nuclei-grpc
COPY . .

RUN go mod tidy
RUN go build -o /usr/local/bin/nuclei-grpc

RUN git clone https://github.com/projectdiscovery/nuclei-templates.git /root/nuclei-templates

EXPOSE 8555

CMD ["nuclei-grpc", "start", "-a", "0.0.0.0", "-p", "8555"]
