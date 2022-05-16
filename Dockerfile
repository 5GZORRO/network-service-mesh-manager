FROM ubuntu:20.04
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY main ./main
COPY config.yaml config.yaml
CMD ./main -config config.yaml