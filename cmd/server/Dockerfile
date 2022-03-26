FROM ubuntu:latest
LABEL maintainer="maintainer@fluffy-bunny.com"
LABEL description="starter kit"

WORKDIR /app
COPY server .
COPY css/ .
COPY js/ .
COPY assets/ .
COPY templates/ .

ENV PATH=/app:${PATH}
ENV PORT="8081"

 
ENTRYPOINT ["/app/server"]