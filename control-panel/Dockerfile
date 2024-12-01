FROM golang:1.23 AS build
WORKDIR /app
COPY . .
RUN go build .

FROM node:22
WORKDIR /app
COPY --from=build /app/controlpanel .
COPY ./web/ .
RUN mkdir /data
RUN yarn install

CMD ["/bin/sh", "-c", "/app/controlpanel & yarn dev:host"]