// Package lzwcodec 用定长码做 LZW 编解码，只用标准库。
package lzwcodec

import (
	"encoding/binary"
	"errors"
)

// FirstCode 是第一个可分配的字典编号。
const FirstCode = 256

var (
	// ErrOddLength 表示数据长度不是 2 的倍数。
	ErrOddLength = errors.New("lzwcodec: 数据长度不是 2 的倍数")
	// ErrCode 表示码不在字典里。
	ErrCode = errors.New("lzwcodec: 未知的码")
)

// Result 是一次编码的结果。
type Result struct {
	Data  []byte
	Codes int
	Dict  int
}

// Encode 编码。缺陷：每个码固定写 2 字节，字典不设上限，末尾不写结束码。
func Encode(data []byte) Result {
	table := map[string]int{}
	for value := 0; value < 256; value++ {
		table[string([]byte{byte(value)})] = value
	}
	next := FirstCode
	out := make([]byte, 0, len(data)*2)
	codes := 0
	emit := func(code int) {
		var buffer [2]byte
		binary.BigEndian.PutUint16(buffer[:], uint16(code))
		out = append(out, buffer[:]...)
		codes++
	}
	current := ""
	for _, value := range data {
		piece := string([]byte{value})
		if _, ok := table[current+piece]; ok {
			current += piece
			continue
		}
		emit(table[current])
		table[current+piece] = next
		next++
		current = piece
	}
	if current != "" {
		emit(table[current])
	}
	return Result{Data: out, Codes: codes, Dict: len(table)}
}

// Decode 解码。缺陷：按固定 2 字节读码，不认 ClearCode / EOICode，字典也没有上限。
func Decode(data []byte) ([]byte, error) {
	if len(data)%2 != 0 {
		return nil, ErrOddLength
	}
	table := map[int]string{}
	for value := 0; value < 256; value++ {
		table[value] = string([]byte{byte(value)})
	}
	next := FirstCode
	out := make([]byte, 0, len(data))
	previous := ""
	for at := 0; at < len(data); at += 2 {
		code := int(binary.BigEndian.Uint16(data[at : at+2]))
		entry, ok := table[code]
		if !ok {
			if code == next && previous != "" {
				entry = previous + previous[:1]
			} else {
				return nil, ErrCode
			}
		}
		out = append(out, entry...)
		if previous != "" {
			table[next] = previous + entry[:1]
			next++
		}
		previous = entry
	}
	return out, nil
}
