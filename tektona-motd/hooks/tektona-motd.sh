# Tektona MOTD — shown on interactive non-SSH shells (VNC terminals).
# SSH sessions are handled by pam_motd via /etc/update-motd.d/10-tektona.

case $- in
    *i*) ;;
    *) return ;;
esac

[ -n "$SSH_CONNECTION" ] && return 0
[ -x /usr/local/bin/tektona-motd ] || return 0

/usr/local/bin/tektona-motd 2>/dev/null || true
