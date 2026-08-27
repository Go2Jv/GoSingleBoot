package Test // 或与包同名，但测试文件通常放在同一包下

import (
	"GoSingleBoot/internal/bizErr"
	"errors"
	"testing"
)

func TestBizErr(t *testing.T) {
	err := bizErr.Warp(200, "hello test")
	bz := bizErr.BizErr{
		Code: 200,
		Msg:  "hello world",
	}

	if errors.As(err, &bz) {
		t.Log(bz.Error())
	}
}
