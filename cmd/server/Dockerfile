FROM ubuntu:latest
LABEL maintainer="maintainer@fluffy-bunny.com"
LABEL description="starter kit"

WORKDIR /app
RUN mkdir /css
RUN mkdir /js
RUN mkdir /assets
RUN mkdir /templates
COPY server .
COPY css/ ./css/
COPY js/ ./js/
COPY assets/ ./assets/
COPY templates/ ./templates/

ENV PATH=/app:${PATH}
ENV PORT="8081"

 
ENTRYPOINT ["/app/server"]