FROM golang:1.16-alpine3.13 AS build

WORKDIR /application
COPY ./application .

RUN go mod tidy
RUN mkdir dist
RUN go build -o dist

FROM alpine:3.13 AS deploy

COPY --from=build /application/dist /application

CMD ["/application/go-docker-learn"]
