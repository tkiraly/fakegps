# build stage
FROM golang:1.13 AS build
ADD go.mod /app/
WORKDIR /app

ARG GOPROXY
ARG MAJOR
ARG MINOR
ARG COMMITCOUNT
ARG GONOSUMDB

RUN CGO_ENABLED=0
RUN go mod download
ADD . /app
RUN ./scripts/build.sh

# final stage
FROM ubuntu:18.04
COPY --from=build /app/fakegps /fakegps
ENTRYPOINT ["/fakegps"]
CMD [ "version" ]
