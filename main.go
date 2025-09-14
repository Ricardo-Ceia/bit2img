package main

import (
	"fmt"
	"log"
	"img2bit/img"
)


func main() {
	var filePath string
	fmt.Print("Please enter the path to your file relative to the current folder: ")
	_, err := fmt.Scanln(&filePath)

	if err != nil {
		log.Fatal(err)
	}

	imgBuffer := img.ReadImgFile(filePath)
	imgGreyscale := img.ImageToGreyscale(imgBuffer)
	log.Println("IMAGE DATA:%s",imgGreyscale)
}
