package main

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"
	"time"

	"img2bit/img"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nfnt/resize"
	"golang.org/x/term"
)

type tickMsg time.Time

type Particle struct {
	X, Y   int
	Color  string
	Symbol string
	Dx, Dy int
}

type Model struct {
	Particles      []Particle
	TerminalWidth  int
	TerminalHeight int
}

func NewModel(imgBuffer image.Image, width, height int) Model {
	downscaled := resizeImageToTerminal(imgBuffer, width, height)
	bounds := downscaled.Bounds()

	const asciiRamp = "@$B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/|()1{}[]?-_+~<>i!lI;:,\"^`'. "
	var particles []Particle

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			r, g, b, _ := downscaled.At(x, y).RGBA()
			gray := uint8((0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)))
			index := int((float64(gray) / 255.0) * float64(len(asciiRamp)-1))
			char := string(asciiRamp[index])
			color := fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r>>8, g>>8, b>>8)
			particles = append(particles, Particle{
				X:      x,
				Y:      y,
				Color:  color,
				Symbol: char,
			})
		}
	}

	return Model{Particles: particles, TerminalWidth: width, TerminalHeight: height}
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.TerminalWidth = msg.Width
		m.TerminalHeight = msg.Height
	case tickMsg:
		for i := range m.Particles {
			if rand.Float64() < 0.05 {
				m.Particles[i].Dx = rand.Intn(3) - 1
				m.Particles[i].Dy = rand.Intn(3) - 1
			} else {
				m.Particles[i].Dx, m.Particles[i].Dy = 0, 0
			}
		}
		return m, tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
	}
	return m, nil
}

func (m Model) View() string {
	canvas := make([][]string, m.TerminalHeight)
	for i := range canvas {
		canvas[i] = make([]string, m.TerminalWidth)
		for j := range canvas[i] {
			canvas[i][j] = " "
		}
	}

	for _, p := range m.Particles {
		x := p.X + p.Dx
		y := p.Y + p.Dy
		if x >= 0 && x < m.TerminalWidth && y >= 0 && y < m.TerminalHeight {
			canvas[y][x] = p.Color + p.Symbol + "\x1b[0m"
		}
	}

	var output string
	for _, row := range canvas {
		for _, cell := range row {
			output += cell
		}
		output += "\n"
	}
	return output
}

func resizeImageToTerminal(imgBuffer image.Image, termWidth int, termHeight int) image.Image {
	if imgBuffer == nil {
		return nil
	}
	aspectRatio := 0.55
	targetHeight := uint(float64(termHeight) * aspectRatio)
	return resize.Resize(uint(termWidth), targetHeight, imgBuffer, resize.Lanczos3)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Error: Please provide a path to an image file.\nUsage: go run . <image_path>")
	}
	filePath := os.Args[1]
	imgBuffer := img.ReadImgFile(filePath)
	if imgBuffer == nil {
		log.Fatal("Error: Could not read or decode image.")
	}
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Error getting terminal size: %v", err)
	}
	p := tea.NewProgram(NewModel(imgBuffer, width, height), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting program: %v", err)
	}
}

