package main
import (
        "bytes"
         "crypto/rsa"
        "fmt"
        "io/ioutil"
        "encoding/binary"
		"encoding/hex"
        "math/big"
)

func main() {
    
    ek, err := ioutil.ReadFile("ek.pub.tss")
    if err != nil {
        panic("ek reading")
    }

    var size uint16

    ekBuff := bytes.NewBuffer(ek)
    err = binary.Read(ekBuff, binary.BigEndian, &size)
    if err != nil {
        panic("public area size")
    }
    fmt.Printf("Size of public EK area: %d\n", size)

    tpmt_public_buff := make([]byte, size)

    err = binary.Read(ekBuff, binary.BigEndian, tpmt_public_buff)
    if err != nil {
        panic("tpmt_public")
    }
    
    tpmt_public := bytes.NewBuffer(tpmt_public_buff)
    var _type uint16
    err = binary.Read(tpmt_public, binary.BigEndian, &_type)
    if err != nil {
        panic("type")
    }
    var alg uint16
    err = binary.Read(tpmt_public, binary.BigEndian, &alg)
    if err != nil {
        panic("algo")
    }

    var attr uint32
    err = binary.Read(tpmt_public, binary.BigEndian, &attr)
    if err != nil {
        panic("attr")
    }

    var digestSize uint16
    err = binary.Read(tpmt_public, binary.BigEndian, &digestSize)
    if err != nil {
        panic("digest size")
    }
    digest := make([]byte, digestSize)

    err = binary.Read(tpmt_public, binary.BigEndian, digest)
    if err != nil {
        panic("digest")
    }

    // ------------ TPMU_PUBLIC_PARMS -------------
    var algSymObject uint16

    // ------------ TPMT_SYM_DEF_OBJECT ------------
    err = binary.Read(tpmt_public, binary.BigEndian, &algSymObject)
    if err != nil {
        panic("algSymObject")
    }
    var rsaKeyBits, symMode uint16
    err = binary.Read(tpmt_public, binary.BigEndian, &rsaKeyBits)
    if err != nil {
        panic("rsa key bits")
    }
    err = binary.Read(tpmt_public, binary.BigEndian, &symMode)
    if err != nil {
        panic("sym mode")
    }
    // ------------ TPMT_SYM_DEF_OBJECT ------------


    var algRsaScheme uint16
    err = binary.Read(tpmt_public, binary.BigEndian, &algRsaScheme)
    if err != nil {
        panic("algRsaScheme")
    }

    err = binary.Read(tpmt_public, binary.BigEndian, &rsaKeyBits)
    if err != nil {
        panic("rsa key bits")
    }

    var exponent uint32

    err = binary.Read(tpmt_public, binary.BigEndian, &exponent)
    if err != nil {
        panic("exponent")
    }

    // ----- TPM2B_PUBLIC_KEY_RSA ------
    err = binary.Read(tpmt_public, binary.BigEndian, &size)
    if err != nil {
        panic("size of TPM2B_PUBLIC_KEY_RSA")
    }

    modulus := make([]byte, size)
    err = binary.Read(tpmt_public, binary.LittleEndian, &modulus)

    if err != nil {
        panic("modulus")
    }

     n := big.NewInt(0)
     n.SetBytes(modulus)

    rsaKey := rsa.PublicKey{N: n, E: int(0x00010001)}
    fmt.Printf("Rsa key: %v\n", rsaKey)

    d, err := hex.DecodeString("46713ed5059824a95e46e6fa1f81fd06586a38a89ce71d8ed4a2990041697428d05f6b73579416359c90d8ff5a61baf7e2ec77ef0aaf6758e42668edfd30133bd136974152e6e0ba7f4f66b393f0e5f8f131622eb6d7caa7ff75d6b7893241a230c32d06a5c1a38d5b1913a418f145b88c910c05f294f9adb775bc9f2be88d4f87d0425aa41fb84457f998a28b8ae450b2462ecc2d63965b4d2fadf12c05be940e20dec9015cc36feb6589937b4a188eea7bcffa5d27b85d1a6f508d2b3356d6ae7f87371395404fbfc0c152e0155fdd1e548a57b461dd48cc79bb1686886fc5e6a01e4aabd08d0298dfacfced7412824836959745e7f4ce475db57de724da91")
	if err != nil {
		panic("d decode")
	}

    dBi := big.Int{}
	dBi.SetBytes(d)

    // Private key structure

	// type PrivateKey struct {
	// PublicKey            // public part.
	// D         *big.Int   // private exponent
	// Primes    []*big.Int // prime factors of N, has >= 2 elements.

	// Precomputed contains precomputed values that speed up private
	// operations, if available.
	// Precomputed PrecomputedValues
	// }

	p1 := big.Int{}
	t1, err := hex.DecodeString("b558750725176f2d2ce3465c9d51b2708e70809d6e96cf2066cffddc43c0f25bcef5c9d10e3b6a2939e932f5eda5eea91e28a6c7882ce92a51bb976ec090c95d3733383b7162bc9c2fc5736226a125dfbdf4ad00773816d29bd922c4f44995707eee1573bd2c728fe3d7d291c2fee3c2b5fc35ed4c792f42a7ded715841b8079")
	if err != nil {
		panic("t1")
	}
	p1.SetBytes(t1)

	p2 :=big.Int{}
	t2, err := hex.DecodeString("fa120dee7e4184cc77fd14d3a834108ea106b82f4fd552027cd84aa4d69866e9caab67a9ab45f334154850d45dedce4c18eefc2ac755d119da6cc79b5986f38b647c091f13ee6555fa0a91c34f4fb68a11d13669479515b3c6707c7e5548bb5dbc89064f1172b588f7f2c623e29f4ae8283a28226f97f9c539b4dec5478f85c7")
	if err != nil {
		panic("t2")
	} 
	p2.SetBytes(t2)

	pk := rsa.PrivateKey{
		PublicKey: rsaKey,
		D: &dBi,
		Primes: []*big.Int{&p1, &p2},
	}

	if err := pk.Validate(); err != nil {
		panic("key is invalid")
	}
}
