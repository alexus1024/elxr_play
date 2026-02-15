# Elxr on Raspberry Pi 5 — Build Log

## Goal
Build a custom Elxr (Debian-based) image for Raspberry Pi 5 using Rucksack, for headless k3s deployment.

## Setup

### Prerequisites
- **WSL 2.6.3.0** with Ubuntu
- **Docker Desktop** (Windows, with WSL integration enabled)
- **Raspberry Pi Imager** (Windows)

### Workspace
- WSL native filesystem: `~/elxr/` (faster than `/mnt/d/` for builds)
- Three repos cloned from GitLab:
  ```
  git clone https://gitlab.com/elxr/tools/rucksack.git
  git clone https://gitlab.com/elxr/tools/pasha.git
  git clone https://gitlab.com/elxr/tools/elxr-config.git
  ```

### Dev Container
Build:
```bash
cd ~/elxr && docker build -t elxr-sdk -f rucksack/docker/Dockerfile rucksack/docker/
```

Run (mount all 3 repos):
```bash
docker run -it --privileged \
  -v $(pwd)/rucksack:/usr/src/rucksack \
  -v $(pwd)/pasha:/usr/src/pasha \
  -v $(pwd)/elxr-config:/usr/src/elxr-config \
  -v /dev:/dev \
  elxr-sdk
```

Install tools inside container:
```bash
pip install ./pasha ./rucksack
```

## Key Findings

### Elxr Architecture
- **Rucksack** (`rs` command) — the image builder, lives in pasha actually (modules)
- **Pasha** — the module library (bootstrap, disk, bootloader, etc.)
- **elxr-config** — YAML manifests + overlays defining what goes into each image
- **Release codenames**: `aria` (current), `bianca` (newer/older)
- **Image types**: `vm/` (virtual machines, x86 UEFI), `boards/` (physical ARM boards)
- **Existing boards**: IMX (NXP), Orin (NVIDIA Jetson) — NO Pi 5 yet

### Pi 5 Specifics
- **SoC**: BCM2712
- **Kernel package**: `linux-image-rpi-2712` (meta-package, currently pulls 6.12.62)
- **Firmware**: `raspi-firmware` (boot firmware, config.txt, start.elf)
- **WiFi/BT**: `firmware-brcm80211`
- **Boot process**: Pi's own bootloader reads `config.txt` from FAT32 partition — NOT UEFI, NOT GRUB
- **Partition layout**: 2 partitions (FAT32 boot + ext4 root), unlike UEFI boards (EFI + BOOT + ROOT)
- **Kernel source**: Raspberry Pi's own apt repo `http://archive.raspberrypi.com/debian bookworm main`

### Packages NOT to include for Pi 5
These are x86/UEFI specific from `packages-elxr.yaml`:
- `systemd-boot` — UEFI bootloader, Pi doesn't use it
- `dracut` — initramfs generator for UEFI boot
- `uefi-ext4` — UEFI filesystem driver
- `linux-image-amd64` — wrong architecture
- `xterm` — not needed for headless

## Troubles & Solutions

### 1. sync-in relative path fails
**Error**: `tar: This does not look like a tar archive`
**Cause**: mmdebstrap's `sync-in` hook resolves paths relative to cwd. Rucksack creates a temp workspace and may change context. Relative paths like `overlay/rpi5/` don't resolve correctly.
**Fix**: Use absolute paths in setup_hooks:
```yaml
setup_hooks:
  - 'sync-in /usr/src/elxr-config/overlay/rpi5/ /'
  - 'sync-in /usr/src/elxr-config/overlay/systemd-networkd/ /'
```
**TODO**: Figure out proper relative path resolution — the IMX configs use relative paths and work. Might be a working directory issue specific to our setup.

### 2. Raspberry Pi repo GPG key not trusted
**Error**: `NO_PUBKEY 82B129927FA3303E`
**Cause**: Elxr's chroot doesn't have the Raspberry Pi archive signing key.
**Fix**:
1. Download the keyring: `curl -sL http://archive.raspberrypi.com/debian/pool/main/r/raspberrypi-archive-keyring/raspberrypi-archive-keyring_2021.1.1+rpt1_all.deb -o /tmp/rpi-keyring.deb`
2. Extract: `dpkg-deb -x /tmp/rpi-keyring.deb /tmp/rpi-keyring`
3. Copy to system keyrings: `cp /tmp/rpi-keyring/usr/share/keyrings/raspberrypi-archive-keyring.gpg /usr/share/keyrings/`
4. Reference in image.yaml:
```yaml
keyring:
  - /usr/share/keyrings/raspberrypi-archive-keyring.gpg
```
5. Use `signed-by` in mirror line:
```yaml
mirrors:
  - deb http://mirror.elxr.dev/elxr aria main non-free-firmware contrib
  - deb [signed-by=/usr/share/keyrings/raspberrypi-archive-keyring.gpg] http://archive.raspberrypi.com/debian bookworm main
```
**Note**: The `keyring` field in pasha passes `--keyring` to mmdebstrap. The key file must exist at the specified path on the host (inside the container).

### 3. Overlay directory permissions
**Error**: `EACCES: permission denied` in VSCode
**Cause**: Directory created as root inside container, WSL user can't write.
**Fix**: `sudo chown -R alexus:alexus ~/elxr/elxr-config/elxr/aria/boards/arm64/rpi5/`

### 4. Overlay file path typo
`overlay/rpi5/lib/systemd/system/` should be `overlay/rpi5/usr/lib/systemd/system/` (matching the original elxr overlay structure).

## Files Created

### elxr-config changes (branch: `feature/rpi5-support`)
```
elxr/aria/boards/arm64/rpi5/minimal/image.yaml   — Pi 5 build manifest
overlay/rpi5/etc/fstab                            — 2-partition fstab (BOOT + ROOT)
overlay/rpi5/etc/hostname                         — "elxr-rpi5"
overlay/rpi5/etc/elxr/variant                     — "edge"
overlay/rpi5/usr/bin/elxr-growfs                  — auto-expand root partition (copied from elxr overlay)
overlay/rpi5/usr/lib/systemd/system/elxr-growfs.service — systemd service for growfs
overlay/rpi5/usr/share/keyrings/raspberrypi-archive-keyring.gpg — RPi repo signing key
```

## TODOs

- [ ] **Build completes** — currently running, first full build attempt
- [ ] **Bootloader setup** — we skipped `installer.bootloader` module. Pi firmware boot handled via customize_hooks (config.txt + cmdline.txt). May need adjustment.
- [ ] **Kernel + DTB placement** — verify raspi-firmware puts kernel/DTBs in the right place on the FAT32 boot partition
- [ ] **fstab mount point** — boot partition mounts at `/boot/firmware` (Raspberry Pi convention). Verify raspi-firmware expects this.
- [ ] **SSH keys** — headless setup needs SSH enabled on first boot (systemctl enable ssh in customize_hooks)
- [ ] **Fix relative paths** — understand why `overlay/rpi5/` doesn't work as relative path in setup_hooks
- [ ] **First boot test** — flash to SD card, boot Pi 5, SSH in
- [ ] **Install k3s** — after successful boot
- [ ] **Reduce image size** — 16GB SD card is tight, optimize packages
- [ ] **Network config** — verify systemd-networkd DHCP works on Pi 5's ethernet interface

## Build Command
```bash
cd /usr/src/elxr-config && rs --workspace /usr/src build --debug --config elxr/aria/boards/arm64/rpi5/minimal/image.yaml
```

## Useful Commands
```bash
# Check what packages are in Elxr's arm64 repo
curl -sL https://mirror.elxr.dev/elxr/dists/aria/main/binary-arm64/Packages.gz | zcat | grep "Package:" | head -20

# Check Raspberry Pi repo packages
curl -sL http://archive.raspberrypi.com/debian/dists/bookworm/main/binary-arm64/Packages.gz | zcat | grep "Package:" | head -20

# Pi 5 kernel versions available
curl -sL http://archive.raspberrypi.com/debian/dists/bookworm/main/binary-arm64/Packages.gz | zcat | grep "Package: linux-image" | grep 2712 | grep -v dbg
```
