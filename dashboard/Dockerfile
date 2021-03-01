ARG ARCH=
FROM ${ARCH}golang:1.16.0-alpine3.13 as builder

WORKDIR /app

COPY . .

RUN ./wizard.sh -c


FROM ${ARCH}alpine

WORKDIR /

EXPOSE 8080

RUN apk update && apk add bash

#COPY --from=builder /wait /wait
COPY --from=builder /app/dashboard /dashboard
COPY --from=builder /app/static /static
COPY --from=builder /app/index.html /index.html
COPY --from=builder /app/wait-for-it.sh /wait-for-it.sh

