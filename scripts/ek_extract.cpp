#include <iostream>
#include <fstream>
#include <cstring>
#include <iomanip>

extern "C" {
    #include "NVMarshal.h"
}

#define BUF_SIZE 4096


void BIGNUM_print(const char *label, const BIGNUM *a, BOOL eol);
using namespace std;

int main() {
    char buffer[BUF_SIZE];

    // Reading tpm2-00.permall to acquire the EPS
    ifstream permstateStream;
    permstateStream.open("tpm2-00.permall", ios::binary);
    if (permstateStream.fail()) {
        return -1;
    }
    permstateStream.seekg (0, permstateStream.end);
    int length = permstateStream.tellg();
    permstateStream.seekg (0, permstateStream.beg);
    if (length == 0 || length > BUF_SIZE ) {
        std::cout << "Buffer length (" << BUF_SIZE << ") not sufficient for state file length (" << length << ")" << std::endl;
        return -1;
    }
    permstateStream.read(buffer, length);
    std::cout << "Read bytes from tpm2-00.permall: " << permstateStream.gcount() << std::endl;

    // Unmarshalling tpm2-00.permall
    INT32 size = length;
    BYTE *bBuffer = (BYTE*)buffer;
    PERSISTENT_DATA pd;
    PERSISTENT_ALL_Unmarshal(&bBuffer, &size, &pd);

    // Read the TPMT_PUBLIC structure representing the endorsement key
    TPM2B_PUBLIC publicArea;
    char publicAreaBuffer[BUF_SIZE];
    BYTE *ptrPa = (BYTE*)publicAreaBuffer;

    ifstream publicAreaStream;
    publicAreaStream.open("ek.pub.tss");
    if (publicAreaStream.fail()) {
        std::cout << "Could not read ek.pub.tss" << std::endl;
        return -1;
    }

    publicAreaStream.seekg (0, publicAreaStream.end);
    int lengthPublicArea = publicAreaStream.tellg();
    publicAreaStream.seekg (0, publicAreaStream.beg);
    if (lengthPublicArea == 0 || lengthPublicArea > BUF_SIZE) {
        std::cout << "Buffer length (" << BUF_SIZE << ") not sufficient for public key area length (" << length << ")" << std::endl;
        return -1;
    }

    publicAreaStream.read(publicAreaBuffer, lengthPublicArea);
    int sizeUnmarshal = lengthPublicArea;
    TPM_RC rc = TPM2B_PUBLIC_Unmarshal(&publicArea, &ptrPa , &sizeUnmarshal, FALSE);
    if (rc != TPM_RC_SUCCESS) {
        std::cout << "TPM2B_PUBLIC_Unmarshal filed" << std::endl;
        return -1;
    }

    // Define the default TPM2B_PUBLIC content as specified by TCG specs
    // EK Profile template: https://github.com/tpm2-software/tpm2-tools/issues/2710
    TPM2B_PUBLIC referencePublic;
    memset(&referencePublic, 0, sizeof(referencePublic));
    referencePublic.publicArea.nameAlg = TPM_ALG_SHA256;
    referencePublic.publicArea.type = TPM_ALG_RSA;
    referencePublic.publicArea.objectAttributes = 0x300b2;
    referencePublic.publicArea.parameters.rsaDetail.keyBits = 2048;
    referencePublic.publicArea.parameters.rsaDetail.symmetric.algorithm = TPM_ALG_AES;
    referencePublic.publicArea.parameters.rsaDetail.symmetric.keyBits.aes = 128;
    referencePublic.publicArea.parameters.rsaDetail.symmetric.mode.aes = TPM_ALG_CFB;
    referencePublic.publicArea.parameters.rsaDetail.scheme.scheme = TPM_ALG_NULL;
    referencePublic.publicArea.unique.rsa.b.size = 256;

    referencePublic.publicArea.authPolicy.b.size = 32;

    BYTE referencePolicy[] = {
        0x83, 0x71, 0x97, 0x67, 0x44, 0x84, 0xb3, 0xf8,
        0x1a, 0x90, 0xcc, 0x8d, 0x46, 0xa5, 0xd7, 0x24,
        0xfd, 0x52, 0xd7, 0x6e, 0x06, 0x52, 0x0b, 0x64,
        0xf2, 0xa1, 0xda, 0x1b, 0x33, 0x14, 0x69, 0xaa
    };

    std::memcpy(referencePublic.publicArea.authPolicy.b.buffer, referencePolicy, 32);

    //memset(&publicArea.publicArea.parameters, 0, sizeof(publicArea.publicArea.parameters));

    DRBG_STATE rand;
    TPM2B_NAME name;

    rc = DRBG_InstantiateSeeded(
        &rand,
        &pd.EPSeed.b,
        PRIMARY_OBJECT_CREATION,
        (TPM2B *)PublicMarshalAndComputeName(&referencePublic.publicArea, &name),
        // in-sensitive, additional data, set all to all zeros by tpm2-tools
        NULL,
        pd.EPSeedCompatLevel
    );

    if (rc != TPM_RC_SUCCESS) {
        std::cout << "DRBG_InstantiateSeeded returned with error" << std::endl;
        return -1;
    }
    OBJECT  object;
    object.publicArea.parameters.rsaDetail.keyBits = 2048;
    // RSA_DEFAULT_PUBLIC_EXPONENT is set in CryptRsaGenerateKey if e is zero
    object.publicArea.parameters.rsaDetail.exponent = RSA_DEFAULT_PUBLIC_EXPONENT;
    rc = CryptRsaGenerateKey(&object, (RAND_STATE*)&rand);

    if (rc != TPM_RC_SUCCESS) {
        std::cout << "CryptRsaGenerateKey failed" << std::endl;
        return -1;
    }


    BN_RSA(bnE);
    BnSetWord(bnE, object.publicArea.parameters.rsaDetail.exponent);
    // First prime factor
    BN_PRIME(bnP);
    BnFrom2B(bnP, &object.sensitive.sensitive.rsa.b);

    // Public modulus
    BN_RSA(bnN);
    BnFrom2B(bnN, &object.publicArea.unique.rsa.b);

    BYTE bModulus[BUF_SIZE];
    NUMBYTES mSize = 0;
    BnToBytes(bnN, bModulus, &mSize);

    BN_PRIME(bnQr);
    BN_PRIME(bnQ);

    // Identify the second prime by division
    BnDiv(bnQ, bnQr, bnN, bnP);
    if(!BnEqualZero(bnQr)) {
        std::cout << "BnDiv did not result in zero remainder" << std::endl;
        return -1;
    }

    NUMBYTES firstPrimeSize = 0;
    NUMBYTES secondPrimeSize = 0;
    BYTE bFirstPrime[BUF_SIZE];
    firstPrimeSize = 0;
    BnToBytes(bnP, bFirstPrime, &firstPrimeSize);

    BYTE bSecondPrime[BUF_SIZE];
    secondPrimeSize = 0;
    BnToBytes(bnQ, bSecondPrime, &secondPrimeSize);

    // N = P * Q. Double check that it matches the modulus obtained from the public key structure
    BN_RSA(bnNVerify);
    BnMult(bnNVerify, bnP, bnQ);

    BOOL pOK;

    // We can compute the private exponent with ComputePrivateExponent,
    // but this will be in format Q, dP, dP, qInv (CRT format)
    pOK = FALSE;
    BN_MAX(bnD);
    BN_RSA(bnPhi);

    // Get compute Phi = (p - 1)(q - 1) = pq - p - q + 1 = n - p - q + 1
    pOK = BnCopy(bnPhi, bnN);
    pOK = pOK && BnSub(bnPhi, bnPhi, bnP);
    pOK = pOK && BnSub(bnPhi, bnPhi, bnQ);
    pOK = pOK && BnAddWord(bnPhi, bnPhi, 1);
    // Compute the multiplicative inverse d = 1/e mod Phi
    pOK = pOK && BnModInverse(bnD, bnE, bnPhi);

    unsigned requiredSize = BnSizeInBits(bnD);
    BYTE bD[BUF_SIZE];
    NUMBYTES dSize = 0;
    BnToBytes(bnD, bD, &dSize);

    if (pOK) {
       std::cout << "Private exponent: " << std::endl;
       for(int i=0; i < dSize; i++) {
           std::cout <<  std::hex << std::setfill('0') << std::setw(2) << static_cast<int>(bD[i]);
       }

       std::cout << std::endl;

       std::cout << "Prime 1" << std::endl;
       for(int i=0; i < firstPrimeSize; i++) {
           std::cout <<  std::hex << std::setfill('0') << std::setw(2) << static_cast<int>(bFirstPrime[i]);
       }

       std::cout << std::endl;

       std::cout << "Prime 2" << std::endl;
       for(int i=0; i < secondPrimeSize; i++) {
           std::cout <<  std::hex << std::setfill('0') << std::setw(2) << static_cast<int>(bSecondPrime[i]);
       }
    }

}
