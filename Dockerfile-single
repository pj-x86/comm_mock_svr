FROM golang
WORKDIR /go/src
COPY ./ comm_mock_svr/
WORKDIR /go/src/comm_mock_svr
#RUN go get -d -v golang.org/x/crypto/ssh github.com/pkg/sftp github.com/kr/fs \
#	github.com/pkg/errors golang.org/x/text golang.org/x/text/transform
RUN go get -d -v ./...
RUN go install -v .
EXPOSE 6610
ENTRYPOINT [ "comm_mock_svr", "-output", "console"]
