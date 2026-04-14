# syntax=docker/dockerfile:1.7

# ----- tektona-motd build stage -----
FROM golang:1.25-bookworm@sha256:29e59af995c51a5bf63d072eca973b918e0e7af4db0e4667aa73f1b8da1a6d8c AS motd-build
WORKDIR /src
COPY tektona-motd/go.mod tektona-motd/go.sum ./
RUN go mod download
COPY tektona-motd/ ./
ARG IMAGE_VERSION=dev
ARG BUILD_DATE=unknown
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags="-s -w -X main.version=${IMAGE_VERSION} -X main.buildDate=${BUILD_DATE}" \
    -o /out/tektona-motd .

# ----- final image -----
FROM ubuntu:24.04@sha256:d1e2e92c075e5ca139d51a140fff46f84315c0fdce203eab2807c7e495eff4f9

LABEL org.opencontainers.image.source="https://github.com/tektona-ai/sandbox-images"
LABEL org.opencontainers.image.description="Tektona sandbox base image — Ubuntu 24.04 with Node.js, code-server, JetBrains Toolbox, and AI coding CLIs"
LABEL org.opencontainers.image.licenses="MIT"

ENV DEBIAN_FRONTEND=noninteractive
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# System packages
RUN apt-get update && apt-get install -y \
    build-essential \
    bzip2 \
    ca-certificates \
    cmake \
    curl \
    dnsutils \
    e2fsprogs \
    file \
    git \
    git-lfs \
    gnupg2 \
    htop \
    iproute2 \
    iputils-ping \
    jq \
    less \
    lsof \
    nano \
    net-tools \
    openssh-client \
    p7zip-full \
    patch \
    pkg-config \
    python3 \
    python3-pip \
    rsync \
    shellcheck \
    strace \
    sudo \
    tree \
    unzip \
    wget \
    xz-utils \
    zstd \
    && rm -rf /var/lib/apt/lists/*

# User setup — sandbox processes run as root inside the runtime, but
# /home/tektona is the user-facing home directory.
RUN useradd -m -s /bin/bash tektona \
    && echo "tektona ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/tektona

# Node.js 22 LTS (pinned binary with integrity check)
ARG NODE_VERSION=22.14.0
ARG NODE_SHA256=69b09dba5c8dcb05c4e4273a4340db1005abeafe3927efda2bc5b249e80437ec
RUN curl -fsSL -o /tmp/node.tar.xz \
    "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.xz" && \
    echo "${NODE_SHA256}  /tmp/node.tar.xz" | sha256sum -c - && \
    tar -xJf /tmp/node.tar.xz -C /usr/local --strip-components=1 && \
    rm /tmp/node.tar.xz

# code-server (pinned .deb with integrity check)
ARG CODE_SERVER_VERSION=4.114.1
ARG CODE_SERVER_SHA256=b3430ea87ac89b3199a1fad61e334922c60ff1cd0b87c1d5dba2ebe07fc38ae0
RUN curl -fsSL -o /tmp/code-server.deb \
    "https://github.com/coder/code-server/releases/download/v${CODE_SERVER_VERSION}/code-server_${CODE_SERVER_VERSION}_amd64.deb" && \
    echo "${CODE_SERVER_SHA256}  /tmp/code-server.deb" | sha256sum -c - && \
    dpkg -i /tmp/code-server.deb && \
    rm /tmp/code-server.deb

# JetBrains Toolbox (pinned with integrity check)
ARG JETBRAINS_TOOLBOX_VERSION=3.4.1.78303
ARG JETBRAINS_TOOLBOX_SHA256=98afad9070063c620d8b107e53a42c3edc63418b496e025bb2230e6aa6e23f99
RUN curl -fsSL -o /tmp/jetbrains-toolbox.tar.gz \
    "https://download.jetbrains.com/toolbox/jetbrains-toolbox-${JETBRAINS_TOOLBOX_VERSION}.tar.gz" && \
    echo "${JETBRAINS_TOOLBOX_SHA256}  /tmp/jetbrains-toolbox.tar.gz" | sha256sum -c - && \
    tar -xzf /tmp/jetbrains-toolbox.tar.gz -C /opt && \
    mv /opt/jetbrains-toolbox-* /opt/jetbrains-toolbox && \
    ln -s /opt/jetbrains-toolbox/bin/jetbrains-toolbox /usr/local/bin/jetbrains-toolbox && \
    rm /tmp/jetbrains-toolbox.tar.gz

# AI coding CLIs (claude-code, codex, opencode) — locked via package.json
COPY sandbox-base/package.json sandbox-base/package-lock.json /tmp/ai-tools/
RUN cd /tmp/ai-tools && npm install --global-style && \
    ln -s /tmp/ai-tools/node_modules/.bin/claude /usr/local/bin/claude && \
    ln -s /tmp/ai-tools/node_modules/.bin/codex /usr/local/bin/codex && \
    ln -s /tmp/ai-tools/node_modules/.bin/opencode /usr/local/bin/opencode

# yq (pinned binary with integrity check)
ARG YQ_VERSION=4.52.4
ARG YQ_SHA256=0c4d965ea944b64b8fddaf7f27779ee3034e5693263786506ccd1c120f184e8c
RUN curl -fsSL -o /tmp/yq \
    "https://github.com/mikefarah/yq/releases/download/v${YQ_VERSION}/yq_linux_amd64" && \
    echo "${YQ_SHA256}  /tmp/yq" | sha256sum -c - && \
    install -m 755 /tmp/yq /usr/local/bin/yq && \
    rm /tmp/yq


# Tektona login banner (SSH via pam_motd, VNC terminals via /etc/bash.bashrc)
COPY --from=motd-build /out/tektona-motd /usr/local/bin/tektona-motd
COPY tektona-motd/hooks/10-tektona /etc/update-motd.d/10-tektona
COPY tektona-motd/hooks/tektona-motd.sh /etc/tektona-motd.sh
RUN chmod +x /usr/local/bin/tektona-motd /etc/update-motd.d/10-tektona && \
    # Remove Ubuntu's default MOTD scripts so only ours renders.
    find /etc/update-motd.d -mindepth 1 -not -name '10-tektona' -delete && \
    : > /etc/motd && \
    : > /etc/legal && \
    # Source our hook from system-wide bashrc so interactive VNC terminals
    # (non-login shells) also render the banner.
    printf '\n# Tektona login banner\n[ -f /etc/tektona-motd.sh ] && . /etc/tektona-motd.sh\n' \
        >> /etc/bash.bashrc

USER root
WORKDIR /home/tektona

# Default entrypoint — replaced by daemon at runtime
ENTRYPOINT ["sleep", "infinity"]
