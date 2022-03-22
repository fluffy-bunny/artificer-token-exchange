FROM scratch
LABEL maintainer="maintainer@fluffy-bunny.com"
LABEL description="starter kit"

WORKDIR /app
COPY server .

ENV PATH=/app:${PATH}
ENV PORT="50051"

ENTRYPOINT ["/app/server"]