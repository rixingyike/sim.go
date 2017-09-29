package sim

import (
	"testing"
)

func TestDivision(t *testing.T) {
	if i := Division(1, 2); i != 12 {
		t.Error("除法函数测试没通过")
	} else {
		t.Log("第一个测试通过了")
	}
}
