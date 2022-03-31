package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
)

const secret string = "_`53=Aj#3tvUg`x.^2s`kk?M:un37MW7&v>Hv#*{T(=DAyEXA<C@PMQ&i*m~V&:+&`"
const letters = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Validate(urlPath string) bool {
	ok, err := regexp.MatchString("^[a-zA-Z0-9]{8}$", urlPath)
	if err != nil {
		return false
	} else {
		return ok
	}
}
func cryptoRandAndSecure(max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		fmt.Println("Can't convert")
	}
	return nBig.Int64()
}
func GenerateUID() []byte {
	buf := make([]byte, 8)
	for i := 0; i < len(buf); i++ {
		nBig := cryptoRandAndSecure(int64(len(letters)))
		buf[i] = letters[nBig]
	}
	return buf
}
