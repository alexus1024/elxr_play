# Elxr Documentation Review — Findings from Building a Pi 5 Image

Based on hands-on experience building an Elxr image for Raspberry Pi 5 using Rucksack/Pasha/elxr-config.

---

## Critical Issues

### 1. No module options documentation (Pasha)
Each module is listed by name only. There is no documentation of what `options:` each module accepts. Users must read Python source code to discover valid options.

For example:
- What options does `bootstrap.mmdebstrap` accept? (suite, architecture, packages, keyring, mirrors, setup_hooks, customize_hooks, etc.)
- What does the `keyring` option do? (see #3 below)
- What does `image.raw` accept? (name, size)
- What does `utils.run_shell` accept? (commands, from_root_dir)

**Impact**: This is the biggest barrier to adoption. We had to read `pasha/modules/bootstrap.py` and `pasha/setup.py` to figure out module names and options.

### 2. No guide for adding new board support
There is no documentation for how to:
- Create a new board directory under `boards/arm64/`
- Write a board-specific image manifest
- Handle non-UEFI bootloaders (Pi, other SBCs)
- Select appropriate kernel/firmware packages
- Create board-specific overlays

**Impact**: We had to reverse-engineer the IMX config and adapt it through trial and error.

### 3. `keyring` option replaces default keyrings silently
The `keyring` field in pasha passes `--keyring` to mmdebstrap, which **replaces** (not supplements) the default trusted keyrings. If you specify a custom keyring for an external repo, you lose the Elxr/Debian keyrings.

**Our experience**: We needed the Raspberry Pi repo key alongside Elxr's key. The interaction between `keyring`, `signed-by` in mirror lines, and the host's `trusted.gpg.d` is undocumented. We solved it by copying the RPi keyring to `/usr/share/keyrings/` on the host and referencing it in the `keyring` field.

### 4. Architecture-specific packages in shared manifests
`manifests/packages-elxr.yaml` contains x86-specific packages:
- `systemd-boot` — UEFI bootloader
- `dracut` — initramfs for UEFI
- `uefi-ext4` — UEFI filesystem driver
- `linux-image-amd64` — wrong architecture

Anyone copying an existing config for arm64 and including `${include:manifests/packages-elxr.yaml}` will get build failures. There is no `packages-elxr-arm64.yaml` equivalent, and no warning about this.

**Suggestion**: Split into `packages-elxr-common.yaml` (arch-independent) and `packages-elxr-amd64.yaml` / `packages-elxr-arm64.yaml`.

### 5. Bootloader module is UEFI-only (undocumented limitation)
The `installer.bootloader` module only supports:
- `systemd-boot` (default)
- `grub-efi`

There is no support for non-UEFI boot (Raspberry Pi, U-Boot, etc.). This is never stated explicitly. We had to skip the bootloader module entirely and handle Pi boot setup via `customize_hooks` and `utils.run_shell`.

**Suggestion**: Document this limitation. Consider adding a generic "copy files to boot partition" module, or a `bootloader.raw` option for boards with their own boot firmware.

---

## Moderate Issues

### 6. Three repos, no relationship map
Nowhere is the relationship between rucksack, pasha, and elxr-config explained:
- **rucksack** = CLI tool (`rs` command), orchestrator
- **pasha** = module library (bootstrap, disk, bootloader, etc.)
- **elxr-config** = YAML manifests + overlays

We had to discover this by exploring the repos. A simple diagram would save significant time.

### 7. mmdebstrap `sync-in` relative path issues
Setup hooks like `sync-in overlay/rpi5/ /` fail with `tar: This does not look like a tar archive`. Absolute paths work: `sync-in /usr/src/elxr-config/overlay/rpi5/ /`.

The existing IMX configs use relative paths (`sync-in overlay/elxr/ /`) and presumably work. The difference may be related to which directory `rs` is invoked from. This behavior is undocumented.

**Note**: This might be a working directory issue specific to our container mount setup, but either way it should be documented.

### 8. Container setup not fully documented
The rucksack README doesn't clearly document:
- The container **must** run with `--privileged` (for loopback devices, mounts)
- The `-v /dev:/dev` mount is required
- All 3 repos should be mounted into `/usr/src/`
- `pip install ./pasha ./rucksack` must be run inside the container
- Disk space requirements (several GB needed)

### 9. Module names in YAML don't match intuitive naming
Module names are registered as stevedore entry_points in `pasha/setup.py`. The naming isn't intuitive:
- `utils.run_shell` (not `cmd.shell` or `shell`)
- `utils.run_command` (not `cmd.run`)
- `image.raw` (not `disk.raw` or `create.image`)

Without reading `setup.py`, you can't know the correct names. We hit `No 'pasha.modules' driver found, looking for 'cmd.shell'` because the naming wasn't documented.

### 10. `${include:...}` syntax undocumented
The YAML include syntax (e.g., `${include:manifests/packages-ostree.yaml}`) is used in configs but never explained:
- Is it OmegaConf? Custom resolver?
- Are paths relative to the config file or working directory?
- What happens if the included file is missing?
- Are nested includes supported?

---

## Minor Issues

### 11. Download page has template URL
The download page shows `wget https://downloads.elxr.dev/<your_image_name>` with no actual filenames, no link to browse available images, and no checksums.

### 12. `dd` command has no safety warnings
The quick start shows `dd if=... of=/dev/sdb bs=1M` with no warning about data loss, no `sudo`, and no instruction to verify the correct device with `lsblk`.

### 13. No mention of Podman as alternative
"Currently no support for Podman" is stated with no explanation of why or timeline.

### 14. Build warnings not documented
Common build warnings like "No zstd in path, using gzip" and "Couldn't identify type of root file system for fsck hook" are not explained. Users can't tell if these are safe to ignore.

---

## Positive Notes

Things that worked well:
- The YAML manifest format is clean and readable once you understand it
- The overlay system is a good pattern for board-specific customization
- `mmdebstrap` with QEMU cross-compilation works reliably for arm64
- The `elxr-growfs` service is a nice touch for auto-expanding root partitions
- Partition creation and formatting stages are solid
- Image compression (zstd) is handled automatically

---

*Generated from hands-on experience building Elxr for Raspberry Pi 5, February 2026.*
