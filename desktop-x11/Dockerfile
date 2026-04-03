ARG BASE_IMAGE
FROM ${BASE_IMAGE}

USER root
ENV DEBIAN_FRONTEND=noninteractive
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Desktop environment (renders to virtio-gpu via Xorg modesetting/KMS)
# - openbox: lightweight window manager
# - plank: macOS-style dock
# No Xvfb/x11vnc/websockify — QEMU's built-in VNC handles display from outside the VM
RUN apt-get update && apt-get install -y \
    dbus-x11 \
    openbox \
    plank \
    xserver-xorg-core \
    xserver-xorg-input-libinput \
    xinit \
    # Wallpaper + cursor
    hsetroot \
    adwaita-icon-theme \
    dmz-cursor-theme \
    dconf-cli \
    # X11 tools (used by sandbox-init desktop API)
    x11-utils \
    x11-xserver-utils \
    wmctrl \
    xdotool \
    imagemagick \
    # X11 libraries
    libx11-6 \
    libxext6 \
    libxfixes3 \
    libxi6 \
    libxrandr2 \
    libxrender1 \
    libxss1 \
    libxtst6 \
    # Desktop apps for AI computer-use agents
    xfce4-terminal \
    pcmanfm \
    mousepad \
    galculator \
    gedit \
    xpdf \
    xpaint \
    libreoffice \
    # Screen recording (agent session replay)
    ffmpeg \
    && rm -rf /var/lib/apt/lists/*

# Google Chrome (Ubuntu 24.04's chromium-browser is a snap wrapper)
RUN curl -fsSL https://dl.google.com/linux/linux_signing_key.pub \
    | gpg --dearmor -o /usr/share/keyrings/google-chrome.gpg \
    && echo "deb [arch=amd64 signed-by=/usr/share/keyrings/google-chrome.gpg] https://dl.google.com/linux/chrome/deb/ stable main" \
    > /etc/apt/sources.list.d/google-chrome.list \
    && apt-get update && apt-get install -y \
    google-chrome-stable \
    fonts-liberation \
    fonts-noto-cjk \
    fonts-noto-core \
    fonts-noto-color-emoji \
    fonts-noto-mono \
    && rm -rf /var/lib/apt/lists/*

# Desktop wallpaper
COPY desktop-bg.png /usr/share/backgrounds/tektona-desktop.png

# Chrome desktop entry with flags to skip first-run setup
RUN sed -i 's|Exec=/usr/bin/google-chrome-stable|Exec=/usr/bin/google-chrome-stable --no-sandbox --no-first-run --no-default-browser-check --disable-sync --disable-gpu-sandbox --disable-features=TranslateUI|' \
    /usr/share/applications/google-chrome.desktop

# Configure DMZ-Black cursor (high visibility for AI computer-use)
RUN mkdir -p /usr/share/icons/default && \
    printf '[Icon Theme]\nInherits=DMZ-Black\n' > /usr/share/icons/default/index.theme

# Plank dock configuration (macOS-style centered bottom dock)
# Stored in /etc/skel/ because /home/tektona is initialized dynamically at runtime.
# The daemon copies these to $HOME on first start.
RUN mkdir -p /etc/skel/.config/plank/dock1/launchers && \
    printf '[PlankDockPreferences]\nHideMode=0\nTheme=Transparent\nPosition=3\nIconSize=48\nDockItems=google-chrome.dockitem;xfce4-terminal.dockitem;pcmanfm.dockitem\n' \
    > /etc/skel/.config/plank/dock1/settings && \
    printf '[PlankDockItemPreferences]\nLauncher=file:///usr/share/applications/google-chrome.desktop\n' \
    > /etc/skel/.config/plank/dock1/launchers/google-chrome.dockitem && \
    printf '[PlankDockItemPreferences]\nLauncher=file:///usr/share/applications/xfce4-terminal.desktop\n' \
    > /etc/skel/.config/plank/dock1/launchers/xfce4-terminal.dockitem && \
    printf '[PlankDockItemPreferences]\nLauncher=file:///usr/share/applications/pcmanfm.desktop\n' \
    > /etc/skel/.config/plank/dock1/launchers/pcmanfm.dockitem

# Plank dconf system defaults (bounce animation on launch, dock always visible)
RUN mkdir -p /etc/dconf/profile /etc/dconf/db/local.d && \
    printf 'user-db:user\nsystem-db:local\n' > /etc/dconf/profile/user && \
    printf '[net/launchpad/plank/docks/dock1]\nanimate-on-launch=true\nhide-mode='"'"'none'"'"'\n' \
    > /etc/dconf/db/local.d/01-plank && \
    dconf update 2>/dev/null || true

# PCManFM desktop wallpaper configuration
RUN mkdir -p /etc/skel/.config/pcmanfm/default && \
    printf '[*]\nwallpaper_mode=stretch\nwallpaper=/usr/share/backgrounds/tektona-desktop.png\ndesktop_bg=#1a1a2e\nshow_documents=0\nshow_trash=0\nshow_mounts=0\n' \
    > /etc/skel/.config/pcmanfm/default/desktop-items-0.conf

USER root

ENV DISPLAY=:0
ENV XCURSOR_THEME=DMZ-Black
ENV XCURSOR_SIZE=32
ENV TEKTONA_IMAGE_ID=desktop-x11
