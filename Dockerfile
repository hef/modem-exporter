FROM golang:1.17.1 as builder
WORKDIR /build
RUN echo nobody:x:65534:65534:nobody:/nonexistent:/sbin/nologin > passwd
ENV CGO_ENABLED=0
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-w -s" -trimpath -o /dist/app .

FROM scratch
COPY --from=builder /build/passwd /etc/passwd
COPY --from=builder /dist/app /
USER nobody
CMD ["/app"]


