package img

import (
				"log"
				"img2bit/fileOperations"
)

func ReadImgFile(filePath string) []byte{
	
	imgBuffer,err := fileOperations.ReadFileMultiThread(filePath)
	
	if err != nil{
		log.Fatal(err)
	}

 return imgBuffer
}
