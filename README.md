This repository contains a Docker configuration file to deploy `swtpm` in container and 
a C application to extract private portion of TPM EK from `swtpm` state file. See the 
corresponding [README for more information](https://github.com/marcoguerri/swtpm-ek-extract/blob/main/ek/README.md).

Why would one want to extract the private counterpart of TPM Endorsement Key? Few possible reasons are the following:
* Purely for educational purposes, e.g. implementing crypto operations natively and comparing the results with swtpm
* To debug unexpected results obtained from swtpm

# swtpm Docker container

The Docker container allows to run `swtpm`, exporting server and control sockets to the host
Execute `/home/swtpm/init.sh` to start `swtpm` instance with socket interface bound to localhost
and create a `socat` tunnel which can be reachable from outside after exporting server and control ports.


# Usage

Start container and redirect server and control ports:

```
sudo docker run -p 2322:2322 -p 2323:2323 -ti swtpm /home/swtpm/init.sh
```

Configure `tpm2-tools` to use `swtpm` TCTI interface:
```
export TPM2TOOLS_TCTI=swtpm:port=2322
```

Startup TPM and run commands against it:
```
$ export TPM2TOOLS_TCTI=swtpm:port=2322
$ tpm2_startup -c                      
$ tpm2_getcap properties-variable
TPM2_PT_PERSISTENT:
  ownerAuthSet:              0
  endorsementAuthSet:        0
  lockoutAuthSet:            0
  reserved1:                 0
  disableClear:              0
  inLockout:                 0
  tpmGeneratedEPS:           1
  reserved2:                 0
TPM2_PT_STARTUP_CLEAR:
  phEnable:                  1
  shEnable:                  1
  ehEnable:                  1
  phEnableNV:                1
  reserved1:                 0
  orderly:                   1
[...]
```

swtpm will show corresponding incoming commands and responses:
```
$ sudo docker run -p 2322:2322 -p 2323:2323 -ti swtpm
[sudo] password for marcoguerri: 
 Ctrl Cmd: length 5
 00 00 00 05 00 
 Ctrl Rsp: length 4
 00 00 00 00 
 SWTPM_IO_Read: length 12
 80 01 00 00 00 0C 00 00 01 44 00 00 
 SWTPM_IO_Write: length 10
 80 01 00 00 00 0A 00 00 00 00 
 Ctrl Cmd: length 5
 00 00 00 05 00 
 Ctrl Rsp: length 4
 00 00 00 00 
 SWTPM_IO_Read: length 22
 80 01 00 00 00 16 00 00 01 7A 00 00 00 06 00 00 
 02 00 00 00 00 7F 
 SWTPM_IO_Write: length 187
 80 01 00 00 00 BB 00 00 00 00 00 00 00 00 06 00 
 00 00 15 00 00 02 00 00 00 04 00 00 00 02 01 80 
 00 00 0F 00 00 02 02 00 00 00 00 00 00 02 03 00 
 00 00 00 00 00 02 04 00 00 00 03 00 00 02 05 00 
 00 00 00 00 00 02 06 00 00 00 40 00 00 02 07 00 
 00 00 03 00 00 02 08 00 00 00 00 00 00 02 09 00 
 00 00 41 00 00 02 0A 00 00 00 00 00 00 02 0B 00 
 00 00 19 00 00 02 0C 00 00 00 00 00 00 02 0D 00 
 00 00 08 00 00 02 0E 00 00 00 00 00 00 02 0F 00 
 00 00 03 00 00 02 10 00 00 03 E8 00 00 02 11 00 
 00 03 E8 00 00 02 12 00 00 00 00 00 00 02 13 00 
 00 00 00 00 00 02 14 00 00 00 00
```

swtpm uses two ports for communication, which are assumed to be at consecutive numbers.


We can then create an Endorsement Key:
```
tpm2_createek -G rsa -c 0x81010002
```

We can then read the public area at the handle specified:
```
tpm2_readpublic -c 0x81010002 -o ek.pub -f pem
```

# EK extraction
See corresponding README in the [ek directory](https://github.com/marcoguerri/swtpm-ek-extract/blob/main/ek/README.md).
