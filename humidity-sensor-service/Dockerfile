ARG ARCH=
FROM ${ARCH}golang:1.16.0-alpine3.13 as builder

WORKDIR /app

COPY . .

RUN ./wizard.sh -c


FROM ${ARCH}alpine

EXPOSE 3030

RUN apk update && apk add bash

#COPY --from=builder /wait /wait
COPY --from=builder /app/humidity-sensor-service /humidity-sensor-service

