# Tektona Sandbox Images

Container images for [Tektona](https://github.com/tektona-ai/tektona) sandbox VMs. These are the rootfs images that run inside QEMU/KVM sandboxes.

## Images

| Image | Description | Pull |
|-------|-------------|------|
| **sandbox-base** | Ubuntu 24.04 base with dev tools, Node.js, code-server, AI coding CLIs | `ghcr.io/tektona-ai/sandbox-base:latest` |
| **desktop-x11** | X11 desktop environment with Chrome, LibreOffice, and desktop apps for AI computer-use | `ghcr.io/tektona-ai/desktop-x11:latest` |

## Structure

```
sandbox-base/          Ubuntu 24.04 base image (all sandboxes inherit from this)
  └── Dockerfile
desktop-x11/           X11 desktop layer (extends sandbox-base)
  ├── Dockerfile
  ├── desktop-bg.png
  └── README.md        Architecture docs, package inventory
```

`desktop-x11` builds on top of `sandbox-base` via the `BASE_IMAGE` build arg.

The AI CLI tools (`@anthropic-ai/claude-code`, `@openai/codex`, `opencode-ai`) are installed globally and unpinned — each build picks up their latest stable release.

## The session user

Sessions run as **`tektona`** (uid 1000), not root, with `HOME=/home/tektona`.

Every way into the sandbox uses that identity: an SSH session, a remote
command, `tektona sandbox cp`, and the desktop session. The images declare
`USER tektona` against a matching `WORKDIR`, which is what keeps the directory
a session starts in and its `$HOME` the same directory — when those differ, a
tool that creates a path with `~` and then writes relative to the session uses
two different places, and remote editors fail to upload their server component.

Because sessions are not root, anything system-wide needs `sudo` — which is
configured NOPASSWD:

```sh
sudo apt-get install -y <pkg>
sudo systemctl status <unit>
sudo npm install -g <pkg>
```

For a root session instead, name it explicitly:

```sh
tektona ssh <sandbox-id> --user root
```

systemd still runs as root as PID 1. The image `USER` applies to the sandbox's
sessions, not to the init system Tektona starts on the image's behalf.

That last part is a property of the platform, not of this image: an init system
cannot run as a non-root process, so on a Tektona whose `sandbox-init` predates
[tektona-ai/tektona#1534](https://github.com/tektona-ai/tektona/pull/1534) this
image does not boot at all. Use a `0.4.x` tag against an older platform.

## Building locally

Both images build from the **repository root** — `sandbox-base` copies in
`tektona-motd/`, so its build context is not its own directory:

```sh
# Base image
docker build -f sandbox-base/Dockerfile -t sandbox-base .

# Desktop image (requires base)
docker build --build-arg BASE_IMAGE=sandbox-base -t desktop-x11 desktop-x11/
```

## License

Apache-2.0
