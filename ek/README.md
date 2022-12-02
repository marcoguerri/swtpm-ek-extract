# EK extraction

`ek_extract.cpp` rebuilds swtpm endorsement key from its state file (by default called `tpm2-00.permall`). The endorsement key
is derived from a seed calculated based on the TPM endorsement hierarchy primary seed, which swtpm stores in the state file.
`ek_extract.cpp` unmarshals the content of `tpm2-00.permall` and builds `TPMT_PUBLIC` object based on TCG default template 
to calculate the name of the key and feeds the result into the key derivation process. Within `TPMT_PUBLIC`, the initial value 
of `unique` (the public key bytes) is zero.

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

Copy `tpm2-00.permall` from swtpm instance, remove header which we are not interested in:
```
$ dd if=tpm2-00.permall.header skip=16 bs=1 of=tpm2-00.permall
```

Execute `ek_extract` by pointing to custom libtpms library:
```
$ LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:/tmp/libtpms/src/.libs" ./ek_extract
Read bytes from tpm2-00.permall: 1299
Private exponent: 
059b78ebae5f808f38e521fd268ef664bcef90fc3e2e3806679118a43c4c6feca0674da6b5673865dbd14f058bbfb19695e8a9ee018f19d14104290d416bf57355635f936485e8723cb0359ce8d1c3d9875a749af3e94b7ce801778cfef79466bd0701ca24f8e8e7d0608e9061633c79aaf6dfe41da8ab22cd8f6bf4a732d127b1c30f3b74f0d2af9a2735c4bc451ea9ddbe6527476bc9597e99cd596b12990d26ef358f734ea077ca5cb432a3a5ef62e8d68e1830e025f5aed531c6de20e6e72aa365a9a31a255e19952f7782f75874b296be3109bb203c9366bf51109b74843bb8493ea52aea4df21f3d0d1c7044d865918a484429ac0a6165807352db5141
Prime 1
e5d12411e2c3a286eedaa478fdb8336ace4c169819eb3694f8f3a90af1d02624e93ba8f3f98a9cee28d36deaa04c628d10d16e952c321f4f13e45a7694b8ef01a430b6202df35dc0e702d6bb1dd6003b89aeb828b8e764cd3760a5907c4b3eaf501bc2a1ca599d1f6db47635e3e3ca4336b2d3dfa897f943959908c9d6abe2df
Prime 2
ccad58866d14380c365eaf552e8629683cc85eb56dec73c959ad77d59975878bdcb2076a6fc0c6f91f7257e734c0eea0e41c51befd0fe791926ad6f8fbd99fdea4ef5fd8fa68accf4d7eea8055e2f546d6a0a90bcf5b7bfd349c389ce29658f8903926c968af72cd64ca97ab0024460e0a480f4d571159ff7c8e0493b2c9b34f
```

## Verify the private/public key pair match

After deriving the keypair directly from the endorsement seed, we can verify that the key created by the TPM
matches our result. Create endorsement key on the TPM:

```
tpm2_createek -u ek.pub.tss -c context.ctx
```

Use `rsa_keypair.go` to verify that the public key created by the TPM matches the key derived from the EPS.
Replace exponent, prime1 and prime2 in `rsa_keypair.go`.

```
 $ go run rsa_keypair.go
Size of public EK area: 314
Rsa key: {23195438192713682485470870009973843005146007915315544307342170175411393239335737355502810102836498661185637594617032550391123583111475287786722668730804512249405050432890919983758683122295351180758485414512554352022266229028903635616082167605265399887474315291602195622990372925293352885824353567480687186121373138270794279301133654386998703015833006341259854319153504531707146516204612711156721354874415734485549214424768509487237106390961710033723013129923910496948465329629663502756282394363092307901413919170848502572422703312781706972736934107454944951785947476881491803298599654402217150657643429509259478691793 65537}
```
