ARG BASE_IMAGE
FROM ${BASE_IMAGE}

LABEL org.opencontainers.image.source="https://github.com/tektona-ai/sandbox-images"
LABEL org.opencontainers.image.description="Tektona desktop sandbox — X11 desktop with Openbox, Chrome, LibreOffice, and Plank dock for AI computer-use agents"
LABEL org.opencontainers.image.licenses="MIT"

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
    libreoffice-writer \
    libreoffice-calc \
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
    fonts-noto-core \
    fonts-noto-color-emoji \
    fonts-noto-mono \
    && rm -rf /var/lib/apt/lists/*

# Desktop wallpaper
COPY desktop-bg.png /usr/share/backgrounds/tektona-desktop.png

# Chrome wrapper: adds --no-sandbox when running as root, plus VM-friendly flags.
# --test-type suppresses the "unsupported command-line flag: --no-sandbox" infobar.
RUN printf '#!/bin/bash\nARGS=(--no-first-run --no-default-browser-check --disable-sync --password-store=basic --disable-features=TranslateUI --disable-gpu --disable-dev-shm-usage)\nif [ "$(id -u)" -eq 0 ]; then\n  ARGS+=(--no-sandbox --test-type)\nfi\nexec /opt/google/chrome/chrome "${ARGS[@]}" "$@"\n' \
    > /usr/local/bin/google-chrome-tektona && chmod +x /usr/local/bin/google-chrome-tektona

# Desktop configuration (single layer for all config files)
RUN set -e && \
    # Chrome: point .desktop at our wrapper
    sed -i 's|Exec=/usr/bin/google-chrome-stable|Exec=/usr/local/bin/google-chrome-tektona|g' \
    /usr/share/applications/google-chrome.desktop && \
    # DMZ-Black cursor (high visibility for AI computer-use)
    mkdir -p /usr/share/icons/default && \
    printf '[Icon Theme]\nInherits=DMZ-Black\n' > /usr/share/icons/default/index.theme && \
    # Plank dock (macOS-style centered bottom dock)
    # Stored in /etc/skel/ because /home/tektona is initialized dynamically at runtime.
    mkdir -p /etc/skel/.config/plank/dock1/launchers && \
    printf '[PlankDockPreferences]\nHideMode=0\nTheme=Transparent\nPosition=3\nIconSize=48\nDockItems=google-chrome.dockitem;xfce4-terminal.dockitem;pcmanfm.dockitem;mousepad.dockitem;gedit.dockitem;galculator.dockitem;libreoffice-writer.dockitem;libreoffice-calc.dockitem\n' \
    > /etc/skel/.config/plank/dock1/settings && \
    for app in google-chrome:google-chrome \
               xfce4-terminal:xfce4-terminal \
               pcmanfm:pcmanfm \
               mousepad:org.xfce.mousepad \
               gedit:org.gnome.gedit \
               galculator:galculator \
               libreoffice-writer:libreoffice-writer \
               libreoffice-calc:libreoffice-calc; do \
        name="${app%%:*}"; desktop="${app##*:}"; \
        printf '[PlankDockItemPreferences]\nLauncher=file:///usr/share/applications/%s.desktop\n' "$desktop" \
        > "/etc/skel/.config/plank/dock1/launchers/${name}.dockitem"; \
    done && \
    # Plank dconf defaults (bounce animation, dock always visible)
    mkdir -p /etc/dconf/profile /etc/dconf/db/local.d && \
    printf 'user-db:user\nsystem-db:local\n' > /etc/dconf/profile/user && \
    printf '[net/launchpad/plank/docks/dock1]\nanimate-on-launch=true\nhide-mode='"'"'none'"'"'\n' \
    > /etc/dconf/db/local.d/01-plank && \
    dconf update 2>/dev/null || true && \
    # Openbox menu (no "Restart" entry)
    mkdir -p /etc/skel/.config/openbox && \
    printf '<?xml version="1.0" encoding="UTF-8"?>\n<openbox_menu xmlns="http://openbox.org/3.4/menu">\n  <menu id="root-menu" label="Desktop">\n    <item label="Terminal"><action name="Execute"><execute>xfce4-terminal</execute></action></item>\n    <item label="File Manager"><action name="Execute"><execute>pcmanfm</execute></action></item>\n    <item label="Chrome"><action name="Execute"><execute>google-chrome-tektona</execute></action></item>\n    <separator />\n    <item label="Reconfigure"><action name="Reconfigure" /></item>\n  </menu>\n</openbox_menu>\n' \
    > /etc/skel/.config/openbox/menu.xml && \
    # Openbox Nightmare theme
    mkdir -p /usr/share/themes/Nightmare/openbox-3 && \
    printf '%s\n' \
    'window.active.border.color: #333333' \
    'window.inactive.border.color: #222222' \
    'window.active.title.bg: flat solid' \
    'window.active.title.bg.color: #1a1a1a' \
    'window.active.label.bg: parentrelative' \
    'window.active.label.text.color: #e0e0e0' \
    'window.active.label.text.font: shadow=y:shadowoffset=1:shadowtint=70' \
    'window.active.handle.bg: flat solid' \
    'window.active.handle.bg.color: #1a1a1a' \
    'window.active.grip.bg: flat solid' \
    'window.active.grip.bg.color: #333333' \
    'window.active.button.unpressed.bg: flat solid' \
    'window.active.button.unpressed.bg.color: #1a1a1a' \
    'window.active.button.unpressed.image.color: #aaaaaa' \
    'window.active.button.pressed.bg: flat solid' \
    'window.active.button.pressed.bg.color: #333333' \
    'window.active.button.pressed.image.color: #ffffff' \
    'window.active.button.hover.bg: flat solid' \
    'window.active.button.hover.bg.color: #333333' \
    'window.active.button.hover.image.color: #ffffff' \
    'window.active.button.disabled.bg: flat solid' \
    'window.active.button.disabled.bg.color: #1a1a1a' \
    'window.active.button.disabled.image.color: #555555' \
    'window.inactive.title.bg: flat solid' \
    'window.inactive.title.bg.color: #0f0f0f' \
    'window.inactive.label.bg: parentrelative' \
    'window.inactive.label.text.color: #666666' \
    'window.inactive.handle.bg: flat solid' \
    'window.inactive.handle.bg.color: #0f0f0f' \
    'window.inactive.grip.bg: flat solid' \
    'window.inactive.grip.bg.color: #222222' \
    'window.inactive.button.unpressed.bg: flat solid' \
    'window.inactive.button.unpressed.bg.color: #0f0f0f' \
    'window.inactive.button.unpressed.image.color: #555555' \
    'window.inactive.button.pressed.bg: flat solid' \
    'window.inactive.button.pressed.bg.color: #222222' \
    'window.inactive.button.pressed.image.color: #aaaaaa' \
    'window.inactive.button.hover.bg: flat solid' \
    'window.inactive.button.hover.bg.color: #222222' \
    'window.inactive.button.hover.image.color: #aaaaaa' \
    'menu.title.bg: flat solid' \
    'menu.title.bg.color: #1a1a1a' \
    'menu.title.text.color: #e0e0e0' \
    'menu.items.bg: flat solid' \
    'menu.items.bg.color: #111111' \
    'menu.items.text.color: #cccccc' \
    'menu.items.active.bg: flat solid' \
    'menu.items.active.bg.color: #333333' \
    'menu.items.active.text.color: #ffffff' \
    'menu.items.disabled.text.color: #555555' \
    'menu.separator.color: #333333' \
    'border.width: 1' \
    'padding.width: 4' \
    'padding.height: 4' \
    'window.handle.width: 2' \
    'menu.separator.width: 1' \
    'menu.separator.padding.width: 6' \
    'menu.separator.padding.height: 3' \
    > /usr/share/themes/Nightmare/openbox-3/themerc && \
    # Openbox rc.xml: use Nightmare theme
    printf '<?xml version="1.0" encoding="UTF-8"?>\n<openbox_config xmlns="http://openbox.org/3.4/rc">\n  <theme>\n    <name>Nightmare</name>\n    <titleLayout>NLIMC</titleLayout>\n    <keepBorder>yes</keepBorder>\n    <animateIconify>yes</animateIconify>\n    <font place="ActiveWindow"><name>sans</name><size>10</size><weight>Bold</weight></font>\n    <font place="InactiveWindow"><name>sans</name><size>10</size><weight>Bold</weight></font>\n    <font place="MenuHeader"><name>sans</name><size>10</size><weight>Bold</weight></font>\n    <font place="MenuItem"><name>sans</name><size>10</size><weight>Normal</weight></font>\n  </theme>\n  <desktops><number>1</number></desktops>\n  <resize><drawContents>yes</drawContents></resize>\n  <focus><followMouse>no</followMouse></focus>\n  <placement><policy>Smart</policy></placement>\n  <menu><file>menu.xml</file></menu>\n</openbox_config>\n' \
    > /etc/skel/.config/openbox/rc.xml && \
    # PCManFM desktop wallpaper
    mkdir -p /etc/skel/.config/pcmanfm/default && \
    printf '[*]\nwallpaper_mode=stretch\nwallpaper=/usr/share/backgrounds/tektona-desktop.png\ndesktop_bg=#1a1a2e\nshow_documents=1\nshow_trash=0\nshow_mounts=0\n' \
    > /etc/skel/.config/pcmanfm/default/desktop-items-0.conf && \
    # xfce4-terminal: use Noto Sans Mono (from fonts-noto-mono) for full coverage
    # of the Unicode glyphs (block-drawing, braille) rendered by tektona-motd.
    mkdir -p /etc/skel/.config/xfce4/terminal && \
    printf '[Configuration]\nFontName=Noto Sans Mono 11\nScrollingUnlimited=FALSE\nScrollingLines=10000\nMiscCursorBlinks=FALSE\nMiscCursorShape=TERMINAL_CURSOR_SHAPE_BLOCK\nMiscMenubarDefault=FALSE\nMiscToolbarDefault=FALSE\nMiscShowUnsafePasteDialog=FALSE\n' \
    > /etc/skel/.config/xfce4/terminal/terminalrc

USER root

ENV DISPLAY=:0
ENV XCURSOR_THEME=DMZ-Black
ENV XCURSOR_SIZE=32
ENV TEKTONA_IMAGE_ID=desktop-x11
