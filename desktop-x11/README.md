# desktop-x11

X11 desktop environment for Tektona sandbox VMs. Provides a graphical desktop for AI computer-use agents.

## Display stack

- **Xorg** with modesetting driver renders to the VM's virtio-gpu
- **openbox** as window manager (lightweight, no compositing)
- **plank** dock with launchers for Chrome, terminal, and file manager
- VNC is handled outside the VM by QEMU — no VNC server inside the image
- Desktop is started on demand, not at boot — headless sandboxes stay lightweight

## Desktop session

The desktop is split between the platform and this image.

**The platform owns the display.** It settles udev, hands over the virtio
device nodes, writes `/etc/X11/xorg.conf` for `/dev/dri/card0` and
`/dev/input/event1|2`, starts Xorg on `:0`, waits for the socket, and applies
the resolution with `xrandr`.

**This image owns the session.** It declares one by shipping an executable at
`/etc/tektona/desktop-session` — the presence of that file is the capability
check. The platform runs it as the session user with `DISPLAY=:0` and `HOME`
already exported, once the X server accepts connections. The file sits under
`/etc/tektona` because it is a declaration the image makes, not a command
anyone runs; it needs no place on `PATH`, and it leaves the bin directories
free for the tools of images layered on top of this one.

The session starts openbox, paints the wallpaper, launches a session bus and
plank, and disables screen blanking. It stays in the foreground: the desktop
is stopped by killing the session and the X server.

It also seeds the desktop configuration from `/etc/skel/.config` into the
session's `$XDG_CONFIG_HOME`, per top-level entry and only where nothing exists
yet. `/etc/skel` is copied into a home directory only when that directory is
created, and neither home here gets it: `/root` predates the image, and
`/home/tektona` is created by `useradd -m` in sandbox-base, layers before these
config files are written. Existing configuration is never overwritten.

## Installed apps

| App | Purpose |
|-----|---------|
| Google Chrome | Web browsing (primary tool for AI agents) |
| LibreOffice | Office suite (docs, spreadsheets, presentations) |
| xfce4-terminal | Terminal emulator |
| PCManFM | File manager |
| Mousepad | Text editor |
| gedit | Text editor (GUI) |
| Galculator | Calculator |
| xpdf | PDF viewer |
| xpaint | Image editor |
| ffmpeg | Screen recording (agent session replay) |

## X11 tooling

| Tool | Package | Used for |
|------|---------|----------|
| `xdotool` | xdotool | Mouse/keyboard input injection |
| `import` | imagemagick | Screenshot capture |
| `xrandr` | x11-xserver-utils | Resolution management |
| `xdpyinfo` | x11-utils | Display info queries |
| `wmctrl` | wmctrl | Window listing |
| `hsetroot` | hsetroot | Wallpaper |

## Building

```sh
docker build --build-arg BASE_IMAGE=sandbox-base -t desktop-x11 desktop-x11/
```
