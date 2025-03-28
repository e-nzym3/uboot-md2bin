# uboot-md2bin
Convert standard output of the U-Boot "md" command to a binary file.

# Install
```
go install github.com/e-nzym3/uboot-md2bin@latest
```

# Usage
```
Usage: uboot-md2bin <fileName>
-- Converts standard U-Boot md command output to a binary file
```
Ensure to only have the memory dump in the file you are trying to process. Strip it from any u-boot command output. See example below for appropriate format. <br><br>
Example:
```
❯ head firmware.bin
80000000: 00 00 00 03 00 00 00 01 ff ff 66 77 00 00 00 00    ..........fw....
80000010: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00    ................
80000020: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00    ................
80000030: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00    ................
80000040: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00    ................
80000050: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00    ................
80000060: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00    ................
80000070: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00    ................
80000080: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00    ................
80000090: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00    ................
❯ uboot-md2bin firmware.bin

╻━━    ┳┳  ┳┓            ┓┏┓┓ •      ━━╻
┃━━━━  ┃┃━━┣┫┏┓┏┓╋   ┏┳┓┏┫┏┛┣┓┓┏┓  ━━━━┃
╹━━    ┗┛  ┻┛┗┛┗┛┗   ┛┗┗┗┻┗━┗┛┗┛┗    ━━╹
   - By enzym3 (github.com/e-nzym3) -

[*] Reading firmware.cur file for processing...
[*] Processing complete! Resulting file is: firmware_output.bin

❯ ls
firmware.dump  firmware_output.bin

```
# Credits
Matt Brown's firmwaretools repo: https://github.com/nmatt0/firmwaretools/tree/master
