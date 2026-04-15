# Launch Doom in a Tektona sandbox

You are an agent with the `tektona` CLI in PATH (pre-authenticated via config/env —
no login step needed). Your job: boot a Tektona sandbox running the Doom demo
image and get Freedoom Phase 2 playing in a noVNC window.

## Image

`ghcr.io/tektona-ai/doom:latest` — Ubuntu 24.04 + Openbox desktop + chocolate-doom
+ Freedoom WADs. Autostarts Freedoom 2 fullscreen once the desktop session is up.

## Steps

### 1. Create the sandbox

```sh
tektona sandbox create \
  --image ghcr.io/tektona-ai/doom:latest \
  --cpu 2 --memory 2
```

Capture the sandbox ID from the output (format: `01K…`). The sandbox starts in
state `building_image` while `sandbox-snapshot-builder` converts the OCI image
into an ext4 snapshot — first boot of a new image is the slow path (minutes);
subsequent boots are cached.

### 2. Wait until it reaches `running`

Poll with JSON output and key off `.state`. Expected transitions:
`building_image` → `scheduling` → `running`.

```sh
SB=<id>
until [ "$(tektona sandbox info "$SB" -o json | jq -r .state)" = "running" ]; do
  sleep 10
done
```

### 3. Start the desktop session

The VM boots to a TTY. The desktop is started on demand — not automatic.

```sh
tektona sandbox start-desktop "$SB"
```

Openbox launches, and its system-wide autostart (`/etc/xdg/openbox/autostart`)
invokes `/usr/games/chocolate-doom -iwad /usr/share/games/doom/freedoom2.wad
-fullscreen -nograbmouse` about 2 seconds later.

### 4. Open the VNC viewer

```sh
tektona sandbox vnc "$SB" --browser
```

Leaves a local proxy running in the foreground and opens noVNC in the default
browser.

### 5. Verify (optional)

Capture a frame and confirm the Freedoom title screen is visible:

```sh
tektona sandbox screenshot "$SB" -o /tmp/doom.png --no-open
```

## Cleanup

```sh
tektona sandbox delete "$SB"
```

## Troubleshooting

- **Stuck in `building_image` for many minutes on a fresh image tag** — normal.
  The snapshot builder is pulling ~2 GiB and converting to ext4.
- **Desktop shows wallpaper but no Doom window** — the autostart hook only fires
  when openbox first starts. If you quit the game, relaunch from a terminal
  inside the sandbox: `tektona sandbox ssh "$SB" -- 'DISPLAY=:0 /usr/games/chocolate-doom -iwad /usr/share/games/doom/freedoom2.wad -fullscreen -nograbmouse &'`.
- **Silent** — audio is not piped through the current VNC stack. Expected until
  audio transport is wired.
