#!/bin/sh
# Entry point of the onboarder image. Leaves nothing installed on the host: it supplies the Wi-Fi
# firmware itself, attaches the card, moves it into this container's network namespace, drops
# privileges and runs `shelly-client onboard`.
#
# Run it privileged, with --pid=host (to reach the host network namespace) and its own network
# namespace (--network=none), and a volume from FIRMWARE_HOST_DIR on the host at /firmware.
#
#   FIRMWARE_HOST_DIR  host path of the /firmware volume, e.g. /var/lib/shelly-onboarder/firmware
#   WIFI_PCI_DEV       PCI address of the Wi-Fi card, e.g. 0000:00:14.3
#   WIFI_DRIVER        kernel driver for it (default iwlwifi)
# Arguments are passed on to `shelly-client onboard`.
set -eu

: "${FIRMWARE_HOST_DIR:?}"
: "${WIFI_PCI_DEV:?}"
driver=${WIFI_DRIVER:-iwlwifi}
dev=/sys/bus/pci/devices/$WIFI_PCI_DEV

# 1. Firmware and regulatory database from this image, served to the kernel from the volume.
#    firmware_class.path is a documented runtime setting; it is not persistent and is set each start.
cp /usr/lib/firmware/iwlwifi-* /usr/lib/firmware/regulatory.db* /firmware/
printf '%s' "$FIRMWARE_HOST_DIR" >/sys/module/firmware_class/parameters/path

# 2. Attach the card to its driver if the boot-time probe failed for lack of firmware.
if [ ! -e "$dev/driver" ]; then
	printf '%s' "$WIFI_PCI_DEV" >"/sys/bus/pci/drivers/$driver/bind"
fi

phy=
for _ in $(seq 30); do
	phy=$(ls "$dev/ieee80211" 2>/dev/null | head -n 1)
	[ -n "$phy" ] && break
	sleep 1
done
[ -n "$phy" ] || { echo "no wiphy under $dev" >&2; exit 1; }

# 3. Move the wiphy from the host network namespace into this one. With --pid=host, $$ is this
#    process's host PID and PID 1 is the host's init. The kernel returns the wiphy to the host
#    namespace when this container stops.
if ! iw phy "$phy" info >/dev/null 2>&1; then
	nsenter -t 1 -n iw phy "$phy" set netns $$
fi

# 4. Keep only what onboard needs (nl80211, addresses on its own interface) and run it.
exec setpriv --bounding-set=-all,+net_admin,+net_raw --inh-caps=-all,+net_admin,+net_raw --no-new-privs \
	/usr/local/bin/shelly-client onboard "$@"
