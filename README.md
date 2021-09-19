# swtpm Docker container

Docker container which allows to run `swtpm`, exporting server and control sockets externally.
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
TPM2_PT_HR_NV_INDEX: 0x0
TPM2_PT_HR_LOADED: 0x0
TPM2_PT_HR_LOADED_AVAIL: 0x3
TPM2_PT_HR_ACTIVE: 0x0
TPM2_PT_HR_ACTIVE_AVAIL: 0x40
TPM2_PT_HR_TRANSIENT_AVAIL: 0x3
TPM2_PT_HR_PERSISTENT: 0x0
TPM2_PT_HR_PERSISTENT_AVAIL: 0x41
TPM2_PT_NV_COUNTERS: 0x0
TPM2_PT_NV_COUNTERS_AVAIL: 0x19
TPM2_PT_ALGORITHM_SET: 0x0
TPM2_PT_LOADED_CURVES: 0x8
TPM2_PT_LOCKOUT_COUNTER: 0x0
TPM2_PT_MAX_AUTH_FAIL: 0x3
TPM2_PT_LOCKOUT_INTERVAL: 0x3E8
TPM2_PT_LOCKOUT_RECOVERY: 0x3E8
TPM2_PT_NV_WRITE_RECOVERY: 0x0
TPM2_PT_AUDIT_COUNTER_0: 0x0
TPM2_PT_AUDIT_COUNTER_1: 0x0
```
