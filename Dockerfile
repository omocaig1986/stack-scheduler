FROM golang:1.13.4-alpine3.10 as build
LABEL stage=builder

RUN apk update && apk add curl git

# install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go
COPY . .

# install deps and build
RUN cd src/scheduler && dep ensure
RUN go build scheduler

FROM alpine:3.10

LABEL org.label-schema.license="MIT" \
    org.label-schema.vcs-url="https://gitlab.com/p2p-faas/stack-scheduler" \
    org.label-schema.vcs-type="Git" \
    org.label-schema.name="p2p-faas/scheduler" \
    org.label-schema.vendor="gabrielepmattia" \
    org.label-schema.docker.schema-version="1.0"

WORKDIR /home/app
COPY --from=build /go/scheduler .

RUN mkdir -p /data

# set permissions
# RUN addgroup -S app && adduser -S -g app app
# RUN chown -R app:app ./
# USER app

EXPOSE 18080
# pprof
EXPOSE 16060 

CMD ["./scheduler"]
