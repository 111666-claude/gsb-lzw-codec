# lzw-codec

LZW（Lempel-Ziv-Welch）编解码，只用 Go 标准库。

```
go test ./...
go run ./cmd/lzw-codec --text TOBEORNOTTOBEORTOBEORNOT --roundtrip
go run ./cmd/lzw-codec --pseudo 10000 --roundtrip
```

## 口径（README 为准）

- **码宽要自适应**：字典编号从 9 位码开始，每写满 2^k 个条目就加一位，上限 16 位。
  `--text TOBEORNOTTOBEORTOBEORNOT` 现在 `codes=16 bytes=32`（每个码固定 2 字节），按口径应该 ≤ 20 字节。
- **要写 ClearCode(256) 与 EOICode(257)**：字典到达 4096 条时发 ClearCode 把字典重置回 258 条，
  编码末尾必须发 EOICode；解码遇到 EOI 就结束，缺失 EOI、截断、码超界、码宽切换错位都要报错并给出字节偏移。
  `--pseudo 10000` 现在 `dict=9620`（字典无上限），按口径 `dict` 不能超过 4096；
  `--pseudo 100000` 现在编号溢出 `uint16`，`roundtrip=false`，按口径必须 `true`。
- 编码/解码按字节进行（不是按 rune），任意二进制输入都要能往返。
- 不变量：`Decode(Encode(data)) == data`（含空输入、全同字节、随机字节）；同一输入重复编码结果逐字节相同；
  `codes` 只数数据码，不含 ClearCode 与 EOICode。
- 规模：10 MB 输入单遍处理，内存与字典条目数成正比（≤ 4096 条），不许整份复制输入。

## 现在的行为

- 每个码固定 `binary.BigEndian.PutUint16` 写 2 字节：`--text TOBEORNOTTOBEORTOBEORNOT` 给 `bytes=32`。
- 没有 ClearCode / EOICode：字典编号一直涨，`--pseudo 10000` 给 `dict=9620`；
  `--pseudo 100000` 的编号超过 65535 被 `uint16` 截断，`roundtrip=false`。
- 解码按固定 2 字节步长读，长度不是 2 的倍数才报错；`--decode-hex 00410042` 照样解出 `text=AB`（缺 EOI 也不报）。

## 计划新增的能力

```
--text 文本       要编码的文本
--repeat N        生成 N 个 'A' 作为输入
--pseudo N        生成 N 个伪随机字节作为输入
--decode-hex HEX  解码十六进制
--roundtrip       编码后解码并比对
```

## 输出

```
codes=16 bytes=18 dict=271 ratio=0.018
roundtrip=true
```

## 目录

```
lzw.go              编解码
cmd/lzw-codec       命令行入口
lzw_test.go         go test 用例
```
