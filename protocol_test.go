// Create by Yale 2018/1/11 17:49
package gosocket

import (
	"fmt"
	"testing"
)

func TestProtocol_Decode(t *testing.T) {

	p := &Protocol{}
	e := p.Encode([]byte("aaa"))
	d := p.Decode(e)
	if d && p.Success {
		fmt.Println(string(p.Data))
	}

}
