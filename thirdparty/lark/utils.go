package lark

import (
	"fmt"
	"io"
	"strings"
)

func getAudioLength(reader io.Reader, bitRateKbps float64) (int, error) {
	// 将数据从 reader 复制到 io.Discard，计算出总字节数
	bytesCopied, err := io.Copy(io.Discard, reader)
	if err != nil {
		return 0, fmt.Errorf("failed to get audio length: %w", err)
	}

	// 将字节数转换为千字节
	sizeInKB := float64(bytesCopied) / 1024.0

	// 计算音频长度，比特率单位为千比特每秒（kbps）
	lengthInMilliSeconds := sizeInKB / bitRateKbps * 1000

	return int(lengthInMilliSeconds), nil
}

func escapePayload(payload string) string {
	// 转义 payload 中的特殊字符
	payload = strings.ReplaceAll(payload, `"`, `\"`)
	return payload
}
