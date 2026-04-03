# desktop-x11

X11 desktop environment for Tektona sandbox VMs. Provides a graphical desktop for AI computer-use agents.

## Display stack

- **Xorg** with modesetting driver renders to the VM's virtio-gpu
- **openbox** as window manager (lightweight, no compositing)
- **plank** dock with launchers for Chrome, terminal, and file manager
- VNC is handled outside the VM by QEMU — no VNC server inside the image
- Desktop is started on demand, not at boot — headless sandboxes stay lightweight

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
docker build --build-arg SANDBOXD_BASE=<base-image> -t desktop-x11 desktop-x11/
```
