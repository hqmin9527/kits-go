package crypto

import (
	"fmt"
	"testing"
)

func TestStorageEncrypt(t *testing.T) {
	srcStr := ""
	crypto := 1
	dstStr := Encrypt([]byte(srcStr), crypto, []byte("xxxxxxxxxxxxxxxx"))
	fmt.Println(string(dstStr))
}
