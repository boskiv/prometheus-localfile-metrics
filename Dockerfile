FROM golang:1.11-stretch as builder

WORKDIR /go/src/github.com/boskiv/prometheus-localfile-metrics

COPY main.go .

RUN go get -d -v

# If you hit the following error:
#     standard_init_linux.go:190: exec user process caused "no such file or directory"
# It means you did not set CGO_ENABLED=0.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o prometheus-localfile-metrics .


FROM alpine
RUN apk --no-cache add ca-certificates

WORKDIR /app/
COPY --from=builder /go/src/github.com/boskiv/prometheus-localfile-metrics/prometheus-localfile-metrics .

CMD ["./prometheus-localfile-metrics"]