package main

import (
	"fmt"
	"log"
	"time"
	"img2bit/img"
)


func main() {
	var filePath string
	fmt.Print("Please enter the path to your file relative to the current folder: ")
	_, err := fmt.Scanln(&filePath)

	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	imgBuffer := img.ReadImgFile(filePath)
	elapsed := time.Since(start)
	log.Println("TIME PROCESSING THE IMAGE:%s",elapsed)
	log.Println("IMAGE DATA:%s",imgBuffer)
}
