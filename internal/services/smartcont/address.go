package smartcont

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/joaojeronimo/go-crc16"
)

func toHumanRepresentationAddr(wc int8, addr string) (string, error) {
	i := new(big.Int)
	i.SetString(addr, 10)

	hexAddr := fmt.Sprintf("%x", i)
	addr, err := packSmcAddr(int8(wc), hexAddr, true, true)
	if err != nil {
		return "", err
	}

	return addr, nil
}

func packSmcAddr(wc int8, hexAddr string, bounceble bool, testnet bool) (string, error) {
	_hexAddr := hexAddr
	if len(hexAddr) < 64 {
		_hexAddr = fmt.Sprintf("%064s", hexAddr)
	}

	tag := 0x11 // for "bounceable" addresses
	if !bounceble {
		tag = 0x51 // for "non-bounceable"
	}
	if testnet {
		tag += 0x80 // if the address should not be accepted by software running in the production network
	}

	var x []byte
	_tag, err := hex.DecodeString(fmt.Sprintf("%x", tag)) // one tag byte
	if err != nil {
		panic(err)
	}
	x = append(x, _tag...)

	_wc := []byte{0x00} // for the basic workchain
	if wc != 0 {
		// one byte containing a signed 8-bit integer with the workchain_id
		tmp, err := hex.DecodeString(strconv.FormatUint(uint64(wc), 16))
		if err != nil {
			panic(err)
		}
		_wc = []byte{tmp[0]}
	}
	x = append(x, _wc[0])

	_addr, err := hex.DecodeString(_hexAddr) // 32 bytes containing 256 bits of the smart-contract address inside the workchain (big-endian)
	if err != nil {
		panic(err)
	}
	x = append(x, _addr...)

	crc := crc16.Crc16(x) // 2 bytes containing CRC16-CCITT of the previous 34 bytes

	crcFix := fmt.Sprintf("%x", crc)
	_crc, err := hex.DecodeString(fmt.Sprintf("%04s", crcFix))
	if err != nil {
		panic(err)
	}

	x = append(x, _crc...)

	encoded := base64.URLEncoding.EncodeToString(x)
	return encoded, nil
}
