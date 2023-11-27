# v0.0.1
FROM node:latest AS frontend
RUN rm -rf ./admin/reminder-admin/.cache
RUN rm -rf ./admin/reminder-admin/dist
RUN rm -rf ./admin/reminder-admin/node_modules
COPY ./admin/reminder-admin/package.json /frontend/package.json
COPY ./admin/reminder-admin/yarn.lock /frontend/yarn.lock
WORKDIR /frontend
RUN yarn install
ADD ./admin/reminder-admin /frontend
RUN WPATH='/admin' yarn run build

FROM --platform=linux/amd64 golang:alpine AS builder
RUN apk add --no-cache git gcc g++
#RUN apk add --no-cache git gcc g++ openssh-client build-base musl-dev
ADD . /go/src
WORKDIR /go/src

RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64  go build -ldflags="-w -s -extldflags=-static" -o go-reminder-bot .

FROM scratch
WORKDIR /root/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src .
COPY --from=frontend /frontend/dist ./admin/reminder-admin/dist/
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/

EXPOSE 2909
CMD [ "./go-reminder-bot" ]
