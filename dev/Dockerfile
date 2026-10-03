FROM fedora:latest

RUN dnf install -y \
    curl \
    unzip \
    jq \
    go \
    && dnf clean all

RUN curl -fsSL https://opencode.ai/install | bash
ENV PATH="${PATH}:/root/.opencode/bin"

WORKDIR /home/clanker/workspace
