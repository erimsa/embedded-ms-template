ARG ARCH=
FROM ${ARCH}golang:1.16.0-alpine3.13 as builder

WORKDIR /app

COPY . .

RUN ./wizard.sh -c


FROM ${ARCH}alpine

WORKDIR /

EXPOSE 3030

RUN apk update && apk add bash

#COPY --from=builder /wait /wait
COPY --from=builder /app/pressure-sensor-service /pressure-sensor-service
COPY --from=builder /app/wait-for-it.sh /wait-for-it.sh
