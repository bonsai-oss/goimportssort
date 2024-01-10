FROM golang AS builder
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /bin/goimportssort -trimpath -ldflags '-s -w' .
RUN strip /bin/goimportssort

FROM golang:1.21.6-alpine
COPY --from=builder /bin/goimportssort /bin/goimportssort
RUN apk add --no-cache git
CMD ["/bin/goimportssort","-v", "-w", "."]
