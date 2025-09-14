package fileOperations

import (
	"fmt"
	"log"
	"io"
	"os"
	"sync"
)


// The size of each chunk to be read.
const chunkSize = 1024 * 1024 // 1 MB



func ReadFileOneThread(filePath string) ([]byte, error) {
	fmt.Println("Reading file with a single thread...")
	// os.ReadFile is the idiomatic way to read an entire file into memory.
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file with one thread: %v", err)
	}
	return data, nil
}

func processChunk(file *os.File, offset int64, size int64, wg *sync.WaitGroup, c chan<- []byte) {
	defer wg.Done() // Signal that this goroutine is done when the function exits

	buffer := make([]byte, size)
	n, err := file.ReadAt(buffer, offset)
	if err != nil && err != io.EOF {
		// Log the error but continue processing other chunks
		log.Printf("failed to read chunk at offset %d: %v", offset, err)
		return
	}

	// Send the successfully read bytes to the channel
	c <- buffer[:n]
}



func ReadFileMultiThread(filePath string) ([]byte, error) {
	fmt.Println("Reading file with multiple threads...")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for multi-threaded read: %v", err)
	}
	defer file.Close() // Ensure the file is closed when the function returns

	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}
	fileSize := info.Size()

	var wg sync.WaitGroup
	var offset int64 = 0

	// Create a channel to receive the chunk data
	chunkDataChan := make(chan []byte)

	// Launch a goroutine for each chunk
	for offset < fileSize {
		currentChunkSize := int64(chunkSize)
		// Adjust the last chunk's size if it's smaller than the full chunk size
		if offset+currentChunkSize > fileSize {
			currentChunkSize = fileSize - offset
		}

		wg.Add(1)
		// Pass the channel to the goroutine
		go processChunk(file, offset, currentChunkSize, &wg, chunkDataChan)

		offset += currentChunkSize
	}

	// Launch a separate goroutine to close the channel when all chunks are processed
	go func() {
		wg.Wait()
		close(chunkDataChan)
	}()

	// Collect the data from the channel into a single variable
	var allData []byte
	for chunk := range chunkDataChan {
		allData = append(allData, chunk...)
	}

	fmt.Println("Finished processing all chunks.")

	return allData, nil
}
