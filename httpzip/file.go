package httpzip

import (
	"compress/flate"
	"encoding/binary"
	"fmt"
	"io"
)

type ZipCloser struct {
	io.Reader
	rc   io.Closer
	body io.Closer
}

func (zc *ZipCloser) Close() error {
	if err := zc.rc.Close(); err != nil {
		_ = zc.body.Close()
		return err
	}
	return zc.body.Close()
}

// File 表示 ZIP 文件中的一个文件条目（Central Directory 中的记录）
type File struct {
	// Name 文件名（相对路径），来源于 Central Directory/File Header
	Name string

	// CompressedSize 压缩后的大小（字节数）
	// 对应 Central Directory 中的 compressed size 字段
	CompressedSize uint64

	// UncompressedSize 解压后的大小（字节数）
	// 对应 Central Directory 中的 uncompressed size 字段
	UncompressedSize uint64

	// Compression 压缩算法标识
	// 0 = Store（无压缩），8 = Deflate，其它见 ZIP 规范
	Compression uint16

	// Offset 文件数据在 ZIP 中的偏移量（相对于整个 ZIP 文件的开头）
	// 一般指向对应 Local File Header 的起始位置
	Offset uint64

	// reader 指向 ZIP Reader，用于按需加载该文件的数据
	// 通常实现为一个 ReaderAt + 解压逻辑
	reader *Reader
}

// Open 打开文件内容，返回一个 io.ReadCloser
func (f *File) Open() (io.ReadCloser, error) {
	// 读取 Local Header
	localHeader, err := f.reader.httpRange(int64(f.Offset), int64(f.Offset)+30+256)
	if err != nil {
		return nil, err
	}

	nameLen := binary.LittleEndian.Uint16(localHeader[26:])
	extraLen := binary.LittleEndian.Uint16(localHeader[28:])

	// 文件数据起始位置 = 偏移量 + header大小 + 文件名 + extra
	dataOffset := f.Offset + 30 + uint64(nameLen) + uint64(extraLen)

	// 读取文件数据
	compData, err := f.reader.httpRangeRaw(int64(dataOffset), int64(dataOffset)+int64(f.CompressedSize)-1)
	if err != nil {
		return nil, err
	}

	switch f.Compression {
	case 0: // Store（无压缩）
		return compData.Body, nil
	case 8: // Deflate
		rc := flate.NewReader(compData.Body)
		return &ZipCloser{
			Reader: rc,
			rc:     rc,
			body:   compData.Body,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported compression: %d", f.Compression)
	}
}
