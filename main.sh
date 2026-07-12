#!/bin/bash

# Target drive definition
DRIVE="/dev/sda"

echo "=== ALERT: Formatting $DRIVE as a single FAT32 partition ==="
echo "Make sure your USB drive is connected as $DRIVE before proceeding."
echo "--------------------------------------------------------"

# 1. Unmount existing partitions to prevent busy errors
echo "[1/3] Unmounting existing partitions..."
sudo umount ${DRIVE}* 2>/dev/null

# 2. Wipe and create a fresh single MBR FAT32 partition layout via fdisk
echo "[2/3] Creating fresh partition table and layout..."
sudo fdisk $DRIVE <<EOF
o
n
p
1


t
b
w
EOF

# 3. Apply the final FAT32 filesystem format
echo "[3/3] Executing final FAT32 format on ${DRIVE}1..."
sudo mkfs.vfat -F 32 "${DRIVE}1"

echo "--------------------------------------------------------"
echo "Success! Your USB drive is now formatted into a single partition."
