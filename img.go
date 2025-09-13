package img

import  "os"
				"io"
				"log"

func readImgFile(filePath) byte[]{
	file,err := os.Open(filePath)
	
	if err != nil{
		log.Fatal(err)
	}
	
	imgBuffer,err := io.Copy(file)
	
	if err != nil{
		log.Fatal(err)
	}
 return imgBuffer
}
