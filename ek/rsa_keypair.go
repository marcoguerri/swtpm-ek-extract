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

    d, err := hex.DecodeString("059b78ebae5f808f38e521fd268ef664bcef90fc3e2e3806679118a43c4c6feca0674da6b5673865dbd14f058bbfb19695e8a9ee018f19d14104290d416bf57355635f936485e8723cb0359ce8d1c3d9875a749af3e94b7ce801778cfef79466bd0701ca24f8e8e7d0608e9061633c79aaf6dfe41da8ab22cd8f6bf4a732d127b1c30f3b74f0d2af9a2735c4bc451ea9ddbe6527476bc9597e99cd596b12990d26ef358f734ea077ca5cb432a3a5ef62e8d68e1830e025f5aed531c6de20e6e72aa365a9a31a255e19952f7782f75874b296be3109bb203c9366bf51109b74843bb8493ea52aea4df21f3d0d1c7044d865918a484429ac0a6165807352db5141")
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
	t1, err := hex.DecodeString("e5d12411e2c3a286eedaa478fdb8336ace4c169819eb3694f8f3a90af1d02624e93ba8f3f98a9cee28d36deaa04c628d10d16e952c321f4f13e45a7694b8ef01a430b6202df35dc0e702d6bb1dd6003b89aeb828b8e764cd3760a5907c4b3eaf501bc2a1ca599d1f6db47635e3e3ca4336b2d3dfa897f943959908c9d6abe2df")
	if err != nil {
		panic("t1")
	}
	p1.SetBytes(t1)

	p2 :=big.Int{}
	t2, err := hex.DecodeString("ccad58866d14380c365eaf552e8629683cc85eb56dec73c959ad77d59975878bdcb2076a6fc0c6f91f7257e734c0eea0e41c51befd0fe791926ad6f8fbd99fdea4ef5fd8fa68accf4d7eea8055e2f546d6a0a90bcf5b7bfd349c389ce29658f8903926c968af72cd64ca97ab0024460e0a480f4d571159ff7c8e0493b2c9b34f")
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
