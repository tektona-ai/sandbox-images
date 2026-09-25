# syntax=docker/dockerfile:1.7

# ----- tektona-motd build stage -----
FROM golang:1.27-bookworm@sha256:69a7b9788769bec032d238959b61854e9ae87f57be9029ec04e9885fabf99195 AS motd-build
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
FROM ubuntu:26.04@sha256:da6fc2be547864451aa253836dd926da33623312df4a9a243e35dc877c378a78

LABEL org.opencontainers.image.source="https://github.com/tektona-ai/sandbox-images"
LABEL org.opencontainers.image.description="Tektona sandbox base image — Ubuntu 26.04 with Node.js and AI coding CLIs"
LABEL org.opencontainers.image.licenses="MIT"

ENV DEBIAN_FRONTEND=noninteractive
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Restore man pages and docs that the Docker image strips, to match a real
# Ubuntu Server install (the sandbox runs as a VM). Must run before the
# package installs so those packages keep their docs. `yes` dies of SIGPIPE
# once unminimize exits, which pipefail would turn into a failed build.
RUN apt-get update && apt-get install -y unminimize && (yes || true) | unminimize \
    && rm -rf /var/lib/apt/lists/*

# Core system packages. Rarely changes; add new packages to the extras layer
# below so this layer's cache survives. dist-upgrade first: the FROM digest is
# pinned, so base packages only pick up Ubuntu security updates when this layer
# rebuilds.
RUN apt-get update && apt-get dist-upgrade -y && apt-get install -y \
    7zip \
    bash-completion \
    bc \
    bind9-dnsutils \
    build-essential \
    bzip2 \
    ca-certificates \
    cmake \
    cron \
    curl \
    e2fsprogs \
    file \
    git \
    git-lfs \
    gnupg2 \
    htop \
    iotop \
    iproute2 \
    iputils-ping \
    jq \
    kmod \
    less \
    locales \
    logrotate \
    lsb-release \
    lsof \
    mtr-tiny \
    nano \
    net-tools \
    netcat-openbsd \
    openssh-client \
    patch \
    pciutils \
    pkg-config \
    procps \
    psmisc \
    python3 \
    python3-pip \
    rsync \
    shellcheck \
    software-properties-common \
    strace \
    sudo \
    sysstat \
    systemd \
    systemd-sysv \
    dbus \
    libpam-systemd \
    tcpdump \
    telnet \
    time \
    traceroute \
    tree \
    tzdata \
    unzip \
    usbutils \
    vim \
    wget \
    xz-utils \
    zip \
    zstd \
    && locale-gen en_US.UTF-8 de_DE.UTF-8 es_ES.UTF-8 fr_FR.UTF-8 pt_PT.UTF-8 pt_BR.UTF-8 \
    && rm -rf /var/lib/apt/lists/*

ENV LANG=en_US.UTF-8

# Point the Ubuntu apt sources at https://. apt checks package signatures either
# way, but plain HTTP exposes what a sandbox installs, and apt fails when
# archive.ubuntu.com stops answering on port 80 but still serves port 443.
# The stock ubuntu image has no CA bundle, so this has to follow the layer that
# installs ca-certificates.
RUN sed -i 's|http://|https://|' /etc/apt/sources.list.d/ubuntu.sources \
    && ! grep -q 'http://' /etc/apt/sources.list.d/ubuntu.sources

# Extra packages. Add new packages here, not to the core layer above, so only
# this layer rebuilds.
RUN apt-get update && apt-get install -y \
    ffmpeg \
    man-db \
    manpages \
    pipx \
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
# - sysstat*: sar history collection. iostat and pidstat do not need it.
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
        dpkg-db-backup.timer \
        sysstat.service \
        sysstat-collect.timer \
        sysstat-summary.timer; do \
        ln -sf /dev/null "/etc/systemd/system/$unit"; \
    done

# Session user. Every way into the sandbox — an SSH session, a remote command,
# `tektona sandbox cp`, the desktop — runs as the image USER, and $HOME follows
# that user's passwd entry. Declaring `USER tektona` (at the end of this file)
# against a matching WORKDIR is what keeps the start directory and $HOME the
# same directory. When they differ, a tool that creates a path with `~` and then
# writes to a path relative to the session uses two different places; editors
# that upload a server component do exactly that and fail with
# "No such file or directory".
#
# Ubuntu ships a stock `ubuntu` account already holding uid 1000, so it goes
# first — otherwise useradd lands tektona on 1001 and `--user 1000` names
# somebody else.
#
# NOPASSWD sudo is load-bearing rather than a convenience now that sessions are
# not root: apt, systemctl and `npm install -g` all need it. 0440 is the mode
# sudo expects of a sudoers drop-in.
RUN userdel --remove ubuntu 2>/dev/null || true; \
    groupdel ubuntu 2>/dev/null || true; \
    useradd --create-home --shell /bin/bash --uid 1000 --user-group tektona \
    && echo "tektona ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/tektona \
    && chmod 0440 /etc/sudoers.d/tektona \
    && install -d -o tektona -g tektona /home/tektona/.local /home/tektona/.local/bin

# Node.js 22 LTS (pinned binary with integrity check)
ARG NODE_VERSION=22.23.3
ARG NODE_SHA256=df450af89261115ef9f9e3830c3eeb2cc9213b63c720b1af623cb5dcbe2e02de
RUN curl -fsSL -o /tmp/node.tar.xz \
    "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.xz" && \
    echo "${NODE_SHA256}  /tmp/node.tar.xz" | sha256sum -c - && \
    tar -xJf /tmp/node.tar.xz -C /usr/local --strip-components=1 && \
    rm /tmp/node.tar.xz

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
ARG YQ_VERSION=4.53.6
ARG YQ_SHA256=c5f056448f973ae7d39b5401949648a78f2dc1947d6a8eb65be60d5c504b9385
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

# The ubuntu image ships /usr/sbin/policy-rc.d exiting 101, which blocks
# `apt install <daemon>` from starting the service (container convention).
# Our images boot as real systemd VMs, so drop it to restore normal
# install-starts-the-service behavior. Kept until here so this file's own
# apt steps still build daemon-free; derived images (desktop-x11) build
# without it, which is safe — maintainer scripts skip service starts when
# systemd isn't booted (no /run/systemd/system in a build container).
RUN rm -f /usr/sbin/policy-rc.d

WORKDIR /home/tektona
USER tektona

# ~/.local/bin is where `pip install --user` and pipx put CLIs. Ubuntu's
# ~/.profile adds it only for login shells, so `tektona ssh <id> -- aws` would
# miss it. sandbox-init hands the image PATH to every session, so set it here.
ENV PATH=/home/tektona/.local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

# PID 1 is systemd. sandbox-init detects /sbin/init in the rootfs and execs it
# in a fresh PID/mount namespace (see internal/services/sandbox-init's
# selectEntrypoint precedence rules). The detection only wins when the OCI
# spec has no CMD — but ubuntu:26.04 inherits CMD ["bash"] from its parent,
# which would otherwise mask systemd. Set TEKTONA_PREFER_SYSTEMD=true so the
# preference branch wins over the hasCmd branch regardless of inheritance.
#
# The USER above does not reach this process. An init system is root's job: as
# a non-root PID 1 it cannot mount cgroups and the boot dies. sandbox-init
# applies the image USER only to an entrypoint the image itself declared, not
# to an init it chose on the image's behalf — so systemd runs as root while
# sessions and the desktop run as tektona. That guard is a platform-side
# requirement of this image: on a sandbox-init without it, this image does
# not boot.
ENV TEKTONA_PREFER_SYSTEMD=true
