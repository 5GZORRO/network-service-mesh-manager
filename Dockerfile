FROM ubuntu:20.04
COPY main ./main
COPY config.yaml config.yaml
CMD ./main -config config.yaml