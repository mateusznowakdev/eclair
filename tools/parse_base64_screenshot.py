#!/usr/bin/env python3

import base64
import sys

if len(sys.argv) != 2:
    print("usage: encoded.txt", file=sys.stderr)
    sys.exit(1)

# load and decode

with open(sys.argv[1], "rt") as encoded_file:
    encoded = encoded_file.read()
    decoded = base64.b64decode(encoded)

# convert into PGM file

print("P2")
print("128 32")
print("255")

for i in range(0, len(decoded), 128):
    chunk = decoded[i : i + 128]

    for s in range(8):
        for b in chunk:
            print("255" if b & (1 << s) else "0")
