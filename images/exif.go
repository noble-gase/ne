package images

import (
	"os"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/shopspring/decimal"
)

// Rect 定义一个矩形框
type Rect struct {
	X int `json:"x"` // 左上角X坐标
	Y int `json:"y"` // 左上角Y坐标
	W int `json:"w"` // 宽度
	H int `json:"h"` // 高度
}

// Orientation 图片的旋转方向
type Orientation int

func (o Orientation) String() string {
	s := ""
	switch o {
	case TopLeft:
		s = "Top-Left"
	case TopRight:
		s = "Top-Right"
	case BottomRight:
		s = "Bottom-Right"
	case BottomLeft:
		s = "Bottom-Left"
	case LeftTop:
		s = "Left-Top"
	case RightTop:
		s = "Right-Top"
	case RightBottom:
		s = "Right-Bottom"
	case LeftBottom:
		s = "Left-Bottom"
	}
	return s
}

const (
	TopLeft     Orientation = 1
	TopRight    Orientation = 2
	BottomRight Orientation = 3
	BottomLeft  Orientation = 4
	LeftTop     Orientation = 5
	RightTop    Orientation = 6
	RightBottom Orientation = 7
	LeftBottom  Orientation = 8
)

// EXIF 定义图片EXIF
type EXIF struct {
	Size        int64
	Format      string
	Width       int
	Height      int
	Orientation string
	Longitude   decimal.Decimal
	Latitude    decimal.Decimal
}

// ParseEXIF 解析图片EXIF
func ParseEXIF(filename string) (*EXIF, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	data := &EXIF{
		Size: stat.Size(),
	}

	// 图片格式
	format, _ := imaging.FormatFromFilename(filename)
	if format < 0 {
		return data, nil
	}
	data.Format = format.String()

	// 解析EXIF
	if x, _ := exif.Decode(f); x != nil {
		// 经纬度
		lat, lng, _ := x.LatLong()
		data.Longitude = decimal.NewFromFloat(lng)
		data.Latitude = decimal.NewFromFloat(lat)
		// 宽
		if tag, _ := x.Get(exif.PixelXDimension); tag != nil {
			v, _ := tag.Int(0)
			data.Width = v
		}
		// 高
		if tag, _ := x.Get(exif.PixelYDimension); tag != nil {
			v, _ := tag.Int(0)
			data.Height = v
		}
		// 转向
		if tag, _ := x.Get(exif.Orientation); tag != nil {
			v, _ := tag.Int(0)
			data.Orientation = Orientation(v).String()
		}
	}

	// 如果宽度或高度为0，则从图片中获取
	if data.Width == 0 || data.Height == 0 {
		if img, _ := imaging.Open(filename); img != nil {
			rect := img.Bounds()
			data.Width = rect.Dx()
			data.Height = rect.Dy()
		}
	}

	return data, nil
}
