package image

import (
	"github.com/disintegration/imaging"
	"image/jpeg"
	"os"
	"testing"
)

func TestResize(t *testing.T) {
	suffix := "gif"
	path := "C:\\Users\\admin\\Desktop\\1558257408278$177$3480148d7f904c88a3890d762a7dde3." + suffix
	img, err := Resize(path, 40, 40, imaging.Lanczos)
	if err != nil{
		t.Log(err)
		return
	}

	f, _ := os.Create("C:\\Users\\admin\\Desktop\\1558257408278$177$3480148d7f904c88a3890d762a7dde3." + suffix)     //创建文件
	defer f.Close()                   //关闭文件
	jpeg.Encode(f, img, nil)
}

func TestResizeAuto(t *testing.T) {
	suffix := "png"
	path := "C:\\Users\\admin\\Desktop\\aboutUs@2x." + suffix
	img, err := ResizeAuto(path, 80, imaging.Lanczos)
	if err != nil{
		t.Log(err)
		return
	}

	f, _ := os.Create("C:\\Users\\admin\\Desktop\\aboutUs@2x_80." + suffix)     //创建文件
	defer f.Close()                   //关闭文件
	jpeg.Encode(f, img, nil)
}