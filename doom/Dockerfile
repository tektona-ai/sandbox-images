ARG BASE_IMAGE=ghcr.io/tektona-ai/desktop-x11:0.3.1
FROM ${BASE_IMAGE}

LABEL org.opencontainers.image.source="https://github.com/tektona-ai/sandbox-images"
LABEL org.opencontainers.image.description="Tektona Doom sandbox — chocolate-doom + Freedoom, autostarts fullscreen in the desktop session"
LABEL org.opencontainers.image.licenses="MIT"

USER root
ENV DEBIAN_FRONTEND=noninteractive
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# chocolate-doom: faithful native Doom engine (renders via SDL to the X display).
# freedoom: FOSS WADs (freedoom1.wad, freedoom2.wad) — no commercial Doom licensing.
RUN apt-get update && apt-get install -y \
    chocolate-doom \
    freedoom \
    && rm -rf /var/lib/apt/lists/*

# Autostart Doom when the openbox session comes up. `/etc/skel/.config/openbox/`
# is the per-user config skeleton copied into /home/tektona at runtime; openbox
# executes `autostart` after the window manager initialises. --nograbmouse keeps
# the pointer usable over VNC (cursor not locked to the window).
RUN mkdir -p /etc/skel/.config/openbox && \
    printf '#!/bin/sh\n(sleep 2 && chocolate-doom -iwad /usr/share/games/doom/freedoom2.wad -fullscreen -nograbmouse) &\n' \
    > /etc/skel/.config/openbox/autostart && \
    chmod +x /etc/skel/.config/openbox/autostart

# Desktop launcher for manual re-launch from the openbox menu / terminal.
RUN cat > /usr/share/applications/freedoom2.desktop <<'EOF'
[Desktop Entry]
Type=Application
Name=Freedoom: Phase 2
Comment=Doom engine with FOSS assets
Exec=chocolate-doom -iwad /usr/share/games/doom/freedoom2.wad -fullscreen -nograbmouse
Icon=chocolate-doom
Terminal=false
Categories=Game;
EOF

# Keep the "desktop" prefix so sandbox-init enables the VNC/desktop code path,
# but distinguish this image from plain desktop-x11.
ENV TEKTONA_IMAGE_ID=desktop-doom
