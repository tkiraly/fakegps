# build stage
FROM golang:1.13 AS build
ADD go.mod /app/
WORKDIR /app
RUN go mod download

ARG GOPROXY
ARG MAJOR
ARG MINOR
ARG COMMITCOUNT
ARG GONOSUMDB

RUN go mod download
ADD . /app
RUN ./scripts/build.sh

# final stage
FROM ubuntu:18.04
COPY --from=build /app/fakegps /fakegps
ENTRYPOINT ["/fakegps"]
CMD [ "version" ]
