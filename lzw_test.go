package lzwcodec

import (
	"bytes"
	"testing"
)

func TestEncodeEmpty(t *testing.T) {
	result := Encode(nil)
	if len(result.Data) != 0 || result.Codes != 0 {
		t.Fatalf("空输入应为空：%+v", result)
	}
}

func TestEncodeCountsCodes(t *testing.T) {
	result := Encode([]byte("AB"))
	if result.Codes != 2 {
		t.Fatalf("两个不同字符应给 2 个码：%d", result.Codes)
	}
	if len(result.Data) == 0 {
		t.Fatal("编码结果不该为空")
	}
}

func TestRoundTrip(t *testing.T) {
	for _, text := range []string{"A", "AB", "banana", "TOBEORNOTTOBEORTOBEORNOT"} {
		result := Encode([]byte(text))
		decoded, err := Decode(result.Data)
		if err != nil {
			t.Fatalf("Decode(%q) 出错：%v", text, err)
		}
		if string(decoded) != text {
			t.Fatalf("往返不一致：%q -> %q", text, decoded)
		}
	}
}

func TestDecodeRoundTripsBytes(t *testing.T) {
	raw := []byte{0, 1, 2, 255, 254, 200, 7}
	result := Encode(raw)
	decoded, err := Decode(result.Data)
	if err != nil {
		t.Fatalf("Decode 出错：%v", err)
	}
	if !bytes.Equal(decoded, raw) {
		t.Fatalf("字节往返不一致：%v -> %v", raw, decoded)
	}
}

func TestDecodeRejectsOddLength(t *testing.T) {
	if _, err := Decode([]byte{1, 2, 3}); err == nil {
		t.Fatal("长度不是 2 的倍数应该报错")
	}
}

func TestDecodeRejectsUnknownCode(t *testing.T) {
	if _, err := Decode([]byte{0x01, 0x2c}); err == nil {
		t.Fatal("未知的码应该报错")
	}
}
