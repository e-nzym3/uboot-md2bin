# uboot-md2bin
Convert standard output of the U-Boot "md" command to a binary file.

# Install
```
go install github.com/e-nzym3/uboot-md2bin@latest
```

# Usage
```
Usage: uboot-md2bin [options] <fileName>
Options:
  --fix    Fix non-sequential memory addresses
  -h, --help    Display this help message
```

# Features
- Converts U-Boot memory dump output to binary files
- Validates input format to ensure correct memory dump structure
- Option to fix non-sequential memory addresses with `--fix` flag

# Input Format
The tool expects each line to follow this exact format:
```
00000000: ff ff ff ff ff ff ff ff ff ff ff ff ff ff ff ff    ................
```
Where:
- First 8 characters: Memory address in hexadecimal
- Followed by colon and space
- 16 space-separated hex bytes
- 4 spaces
- 16 ASCII characters representing the hex bytes
- Line ending

<b>Note: This is the standard format that the u-boot command `md.b` spits out. </b><br><br>

The tool will validate this format before processing and exit if any line doesn't match.

# Memory Address Fixing
When using the `--fix` flag, the tool will:
1. Create a new file with "_fixed" appended to the original filename
2. Fix any non-sequential memory addresses to be sequential
3. Process the fixed version to create the binary output

This is useful when dealing with memory dumps that have gaps or non-sequential addresses.

# Example
```
# uboot-md2bin --fix rootfs2.dump 

╻━━    ┳┳  ┳┓            ┓┏┓┓ •      ━━╻
┃━━━━  ┃┃━━┣┫┏┓┏┓╋   ┏┳┓┏┫┏┛┣┓┓┏┓  ━━━━┃
╹━━    ┗┛  ┻┛┗┛┗┛┗   ┛┗┗┗┻┗━┗┛┗┛┗    ━━╹
   - By enzym3 (github.com/e-nzym3) -

[*] Reading rootfs2.dump file for processing...
[*] Validating contents...
[*] Fixing non-sequential memory addresses...
[*] Created fixed version: rootfs2_fixed.dump
[*] Processing complete! Resulting file is: rootfs2_output.bin
```

# Credits
https://github.com/nmatt0/firmwaretools/tree/master
https://github.com/gmbnomis/uboot-mdb-dump
