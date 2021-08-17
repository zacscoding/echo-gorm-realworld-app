FROM golang:1.15-alpine AS build

RUN mkdir -p /go/src/github.com/zacscoding/echo-gorm-realworld-app ~/.ssh && \
    apk add --no-cache git openssh-client make gcc libc-dev
WORKDIR /go/src/github.com/zacscoding/echo-gorm-realworld-app
COPY . .
RUN make build

FROM alpine:3
COPY --from=build /go/src/github.com/zacscoding/echo-gorm-realworld-app/app-server /bin/app-server
CMD /bin/app-server