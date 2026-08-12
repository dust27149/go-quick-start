package compress

import (
	"compress/gzip"
	"io"
	"os"
)

// CompressGzipFile 压缩指定的源文件到目标文件，使用 gzip 格式。
func CompressGzipFile(sourcePath string) (err error) {
	targetPath := sourcePath + ".gz"
	// 打开源文件
	src, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer src.Close()

	// 创建目标文件
	target, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer target.Close()

	// 创建 gzip 写入器
	gz := gzip.NewWriter(target)
	defer gz.Close()
	// 将源文件内容复制到 gzip 写入器中，实现压缩
	if _, err := io.Copy(gz, src); err != nil {
		return err
	}
	// 关闭 gzip 写入器，确保所有数据都被写入目标文件
	return gz.Close()
}
