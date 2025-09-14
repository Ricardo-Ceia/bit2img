package img

import (
				"log"
				"image"
				"image/color"
				_ "image/jpeg"
				"os"
)

func ReadImgFile(filePath string) image.Image{
	
	file,err := os.Open(filePath)
	
	defer file.Close()

	if err!=nil{
		log.Fatal(err)
	}

	img,format,err := image.Decode(file)
	log.Printf("Format of the image:%s",format)
	if err!=nil{
		log.Fatal(err)
	}

	return img
}

func ImageToGreyscale(img image.Image)  *image.Gray{
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y;y < bounds.Max.Y;y++{
		for x := bounds.Min.X; x < bounds.Max.X;x++{
			originalColor := img.At(x,y)
			grayColor := color.GrayModel.Convert(originalColor)
			grayImg.Set(x,y,grayColor)
		}
	}
	return grayImg
}
