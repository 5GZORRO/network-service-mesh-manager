FROM ubuntu:18.04
RUN apt update && apt install -y python3 python3-pip software-properties-common
RUN apt install -y wireguard wireguard-tools iproute2 openresolv
COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt
CMD tail -f /dev/null