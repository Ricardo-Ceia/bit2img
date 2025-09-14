package main

import (
	"fmt"
	"log"
	"img2bit/img"
	"image"
	"image/color"
	"image/jpeg"
	tea "github.com/charmbracelet/bubbletea"
)


type Model struct{
	Image *image.Gray
	TerminalWidth int 
	TerminalHeight int
	CursorRow int
	IsQuitting 	bool
	Error error
}

func NewModel(img *image.Gray,width,height) Model{
	return Model{
		Image: img,
		TerminalWidth: width,
		TerminalHeight:	height,
		CursorRow:	0,
		IsQuitting: false,
		Error: nil,
	}
}

func main() {
	if len(os.Args < 2){
		log.Fatal("Usage: go run main. go <iamge_path>")
	}	

	filePath := os.Args[1]

	imgBuffer := img.ReadImgFile(filePath)

	grayImg := img.ImageToGrey(imgBuffer)

	width,height,err := term.GetSize(int(os.Strdin.Fd()))

	if err != nil{
		log.Fatal(err)
	}

	p := tea.NewProgram(NewModel(grayImg,width,height))

	if _,err := p.Run(); err != nil{
		log.Fatal("Error starting program: %v",err)
	}
}
