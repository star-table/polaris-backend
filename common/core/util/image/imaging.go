package image

import (
	"fmt"
	"github.com/star-table/common/core/logger"
	"github.com/disintegration/imaging"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

var log = logger.GetDefaultLogger()

func Resize(sourcePath string, width, height int, filter imaging.ResampleFilter) (*image.NRGBA, error){
	file, _ := os.Open(sourcePath)
	defer func() {
		if err := file.Close(); err != nil{
			fmt.Println(err)
		}
	}()

	img, _, err := image.Decode(file)
	if err != nil{
		fmt.Println(err)
		return nil, err
	}
	return imaging.Resize(img, width, height, filter), nil
}

//ratio: 比例 = 长:宽
func ResizeAuto(sourcePath string, height int, filter imaging.ResampleFilter) (*image.NRGBA, error){
	file, _ := os.Open(sourcePath)
	defer func() {
		if err := file.Close(); err != nil{
			fmt.Println(err)
		}
	}()

	imgConfig, _, err := image.DecodeConfig(file)
	if err != nil{
		return nil, err
	}

	file, _ = os.Open(sourcePath)
	img, _, err := image.Decode(file)
	if err != nil{
		fmt.Println(err)
		return nil, err
	}

	width := imgConfig.Width / (imgConfig.Height / height)

	return imaging.Resize(img, width, height, filter), nil
}