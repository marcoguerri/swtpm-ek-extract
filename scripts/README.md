# EK extraction

`ek_extract.cpp` rebuilds swtpm endorsement key from its state file (by default called `tpm2-00.permall`). The endorsement key
is derived from a seed calculated based on the TPM endorsement hierarchy primary seed, which swtpm stores in the state file.
`ek_extract.cpp` unmarshals the content of `tpm2-00.permall` and `TPMT_PUBLIC` structure corresponding to the Endorsement Key. 
`TPMT_PUBLIC` contribution is necessary because TPM key derivation depends also on the name of the key, which is a hash over 
the serialization of the `TPMT_PUBLIC` structure. Within `TPMT_PUBLIC`, the initial value of `unique` (the public key bytes)
is zero.

The code needs to be compiled against `libtpms` with two amendments:
* Some symbols need to be exported to be linked against. This can be done through `src/libtpms.syms`, by setting the following symbols as non-local:
```
PERSISTENT_ALL_Unmarshal;
CryptRsaGenerateKey;
TPM_PrintAll;
TPMLIB_LogPrintf;
BnSetWord;
BnInit;
ComputePrivateExponent;
BnFrom2B;
BnDiv;
BnMult;
BnCopy;
BnSub;
BnAddWord;
BnModInverse;
BnToBytes;
BnSizeInBits;
TPM2B_PUBLIC_Unmarshal;
PublicMarshalAndComputeName;
DRBG_InstantiateSeeded;
PRIMARY_OBJECT_CREATION;
```

* `PERSISTENT_ALL_Unmarshal` needs to be patched to set `PERSISTENT_DATA` structure, which contains also the `EPSeed`, the primary seed
of the endorsement hierarchy.

The file `patch` contains a patch that can be applied 

# Example

## Setup swtpm
```
$ git clone https://github.com/stefanberger/libtpms.git
$ ./bootstrap.sh
$ ./configure

[export symbols from `src/libtpms.syms`]
```

Apply `0001-Modify-libtpms-to-export-some-symbols.patch`.

## Compile
```
$ g++ ek_extract.cpp -o ek_extract  \
    -I /tmp/libtpms/src/tpm2 \
    -I /tmp/libtpms/src/tpm2/crypto/openssl \
    -I /tmp/libtpms/src/tpm2/crypto \
    -I /tmp/libtpms/include/libtpms \
    -DTPM_POSIX \
    -ltpms \
    -g \
    -L /tmp/libtpms/src/.libs
```

Execute by pointing to custom libtpms library:
```
$ LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:/tmp/libtpms/src/.libs" ./ek_extract 
Read bytes from tpm2-00.permall: 1299
Private exponent: 
46713ed5059824a95e46e6fa1f81fd06586a38a89ce71d8ed4a2990041697428d05f6b73579416359c90d8ff5a61baf7e2ec77ef0aaf6758e42668edfd30133bd136974152e6e0ba7f4f66b393f0e5f8f131622eb6d7caa7ff75d6b7893241a230c32d06a5c1a38d5b1913a418f145b88c910c05f294f9adb775bc9f2be88d4f87d0425aa41fb84457f998a28b8ae450b2462ecc2d63965b4d2fadf12c05be940e20dec9015cc36feb6589937b4a188eea7bcffa5d27b85d1a6f508d2b3356d6ae7f87371395404fbfc0c152e0155fdd1e548a57b461dd48cc79bb1686886fc5e6a01e4aabd08d0298dfacfced7412824836959745e7f4ce475db57de724da91
Prime 1
b558750725176f2d2ce3465c9d51b2708e70809d6e96cf2066cffddc43c0f25bcef5c9d10e3b6a2939e932f5eda5eea91e28a6c7882ce92a51bb976ec090c95d3733383b7162bc9c2fc5736226a125dfbdf4ad00773816d29bd922c4f44995707eee1573bd2c728fe3d7d291c2fee3c2b5fc35ed4c792f42a7ded715841b8079
Prime 2
fa120dee7e4184cc77fd14d3a834108ea106b82f4fd552027cd84aa4d69866e9caab67a9ab45f334154850d45dedce4c18eefc2ac755d119da6cc79b5986f38b647c091f13ee6555fa0a91c34f4fb68a11d13669479515b3c6707c7e5548bb5dbc89064f1172b588f7f2c623e29f4ae8283a28226f97f9c539b4dec5478f85c7
```

## Verify the private/public key pair match

After deriving the keypair directly from the endorsement seed, we can verify that the key created by the TPM
matches our result. Create endorsement key on the TPM:

```
tpm2_createek -u ek.pub.tss -c context.ctx
```

Use `rsa_keypair.go` to verify that the public key created by the TPM matches the key derived from the EPS.
Replace exponent, prime1 and prime2 in `rsa_keypair.go`.
