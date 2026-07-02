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
LABEL org.opencontainers.image.description="Tektona sandbox base image — Ubuntu 24.04 with Node.js, code-server, and AI coding CLIs"
LABEL org.opencontainers.image.licenses="MIT"

ENV DEBIAN_FRONTEND=noninteractive
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# System packages. dist-upgrade first: the FROM digest is pinned, so base
# packages only pick up Ubuntu security updates when this layer rebuilds.
RUN apt-get update && apt-get dist-upgrade -y && apt-get install -y \
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
    systemd \
    systemd-sysv \
    dbus \
    libpam-systemd \
    tree \
    unzip \
    wget \
    xz-utils \
    zstd \
    && rm -rf /var/lib/apt/lists/*

# Mask systemd units that would fight sandbox-init's pre-boot setup.
# Done as symlinks (not `systemctl mask`) because there's no running systemd
# inside the build context.
# - systemd-resolved/-networkd: sandbox-init owns /etc/resolv.conf and the NIC
#   (configured from kernel cmdline before switchRoot).
# - systemd-modules-load: the custom QEMU kernel ships drivers built-in;
#   there's no matching /lib/modules tree to load against.
# - systemd-firstboot: blocks boot prompting for locale/timezone on a "fresh"
#   rootfs, which every snapshot clone looks like once machine-id is reset.
RUN for unit in \
        systemd-resolved.service \
        systemd-networkd.service \
        systemd-networkd.socket \
        systemd-modules-load.service \
        systemd-firstboot.service; do \
        ln -sf /dev/null "/etc/systemd/system/$unit"; \
    done

# Mask units that are pure overhead in an ephemeral, single-use sandbox VM.
# These delay boot and/or burn cycles on a disposable rootfs that is discarded
# at shutdown. Masking (symlink to /dev/null) keeps systemd as init but turns
# each unit into a no-op. None of these are needed for SSH/exec, the desktop,
# or the agent-runtime.
# - systemd-timesyncd: the host provides time to the guest; NTP is redundant.
# - e2scrub_*: online fsck of an ephemeral, snapshot-backed ext4 rootfs.
# - motd-news + apt-daily*: package/news upkeep on a throwaway image.
# - fstrim / dpkg-db-backup: disk maintenance pointless on a disposable VM.
RUN for unit in \
        systemd-timesyncd.service \
        e2scrub_reap.service \
        e2scrub_all.timer \
        motd-news.service \
        motd-news.timer \
        apt-daily.service \
        apt-daily.timer \
        apt-daily-upgrade.service \
        apt-daily-upgrade.timer \
        fstrim.timer \
        dpkg-db-backup.timer; do \
        ln -sf /dev/null "/etc/systemd/system/$unit"; \
    done

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

# The Node tarball bundles npm 10, whose vendored sigstore 3.x carries
# CVE-2026-48815; npm 11 bundles sigstore 4.x.
RUN npm install -g npm@11

# AI coding CLIs (claude-code, codex, opencode) — installed globally so npm
# owns the /usr/local/bin symlinks. Unpinned: each build picks up the latest
# stable release, which is what users expect.
RUN npm install -g \
    @anthropic-ai/claude-code \
    @openai/codex \
    opencode-ai

# yq (pinned binary with integrity check)
ARG YQ_VERSION=4.53.2
ARG YQ_SHA256=d56bf5c6819e8e696340c312bd70f849dc1678a7cda9c2ad63eebd906371d56b
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

# Refresh /etc/ld.so.cache at build time so ldconfig.service short-circuits at
# boot: it only runs when the cache is stale relative to the library dirs
# (ConditionNeedsUpdate), so writing it here as the last /etc change makes the
# boot-time unit a no-op. Must stay near the end — any later step that writes
# under /etc/ can re-stale the condition.
RUN ldconfig

# ubuntu:24.04 ships /usr/sbin/policy-rc.d exiting 101, which blocks
# `apt install <daemon>` from starting the service (container convention).
# Our images boot as real systemd VMs, so allow all actions (exit 0) to
# restore normal install-starts-the-service behavior. Overwritten rather
# than removed: a cross-layer rm needs an OCI whiteout, which not every
# rootfs converter honors, while an overwrite is plain layer content.
# Kept until here so this file's own apt steps still build daemon-free;
# derived images (desktop-x11) build with exit 0, which is safe —
# maintainer scripts skip service starts when systemd isn't booted
# (no /run/systemd/system in a build container).
RUN printf '#!/bin/sh\nexit 0\n' > /usr/sbin/policy-rc.d && chmod 755 /usr/sbin/policy-rc.d

USER root
WORKDIR /home/tektona

# PID 1 is systemd. sandbox-init detects /sbin/init in the rootfs and execs it
# in a fresh PID/mount namespace (see internal/services/sandbox-init's
# selectEntrypoint precedence rules). The detection only wins when the OCI
# spec has no CMD — but ubuntu:24.04 inherits CMD ["bash"] from its parent,
# which would otherwise mask systemd. Set TEKTONA_PREFER_SYSTEMD=true so the
# preference branch wins over the hasCmd branch regardless of inheritance.
ENV TEKTONA_PREFER_SYSTEMD=true
