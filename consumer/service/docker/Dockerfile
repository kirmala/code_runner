FROM ubuntu:latest

RUN apt-get update && \
    apt-get install -y \
    clang \
    g++ \
    python3 \
    && rm -rf /var/lib/apt/lists/*

RUN update-alternatives --install /usr/bin/python python /usr/bin/python3 1

WORKDIR /app

CMD ["python3", "--version"]