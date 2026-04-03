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
  ├── Dockerfile
  ├── package.json     AI CLI tools (claude-code, codex, opencode) — locked deps
  └── Taskfile.yaml
desktop-x11/           X11 desktop layer (extends sandbox-base)
  ├── Dockerfile
  ├── desktop-bg.png
  └── README.md        Architecture docs, package inventory
```

`desktop-x11` builds on top of `sandbox-base` via the `BASE_IMAGE` build arg.

## Updating AI CLI tools

```sh
task update-deps
```

This updates `@anthropic-ai/claude-code`, `@openai/codex`, and `opencode-ai` to latest and regenerates the lockfile.

## Building locally

```sh
# Base image
docker build -t sandbox-base sandbox-base/

# Desktop image (requires base)
docker build --build-arg BASE_IMAGE=sandbox-base -t desktop-x11 desktop-x11/
```

## License

Apache-2.0
