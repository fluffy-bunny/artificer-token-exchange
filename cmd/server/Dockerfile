FROM alpine:latest as stage1

RUN apk add -U --no-cache ca-certificates

RUN apk add bash curl && \
  curl -fSsL https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh -o /wait-for-it.sh && \
  chmod +x /wait-for-it.sh && \
  apk del curl

FROM alpine:latest as alpine
RUN apk add --no-cache bash

LABEL maintainer="maintainer@fluffy-bunny.com"
LABEL description="starter kit"

WORKDIR /
COPY --from=stage1 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=stage1 /etc/passwd /etc/passwd
 


WORKDIR /app
RUN mkdir /static
COPY static/ ./static/
COPY server ./


ENV PATH=/app:${PATH}
ENV PORT="80"
EXPOSE 80
 
ENTRYPOINT ["/app/server"]