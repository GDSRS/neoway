package main

import (
	"errors"
	"fmt"
	"io"
	"neoway/config"
	"neoway/database"
	"neoway/utils"
	"os"
	"strings"
)

const (
	filePath  = "./base_teste.txt"
	chunkSize = 1024 * 1024 // Max chunk size to read from file, in bytes.

)

func findLatestChar(arr []byte, sentinel rune) (int, error) {
	for i := len(arr) - 1; i > -1; i -= 1 {
		if arr[i] == byte(sentinel) {
			return i, nil
		}
	}
	return -1, errors.New("Sentinel not found :-(")
}

func skipFirstLine(file *os.File) {
	numberBytesRead := 0
	newLineIndex := -1
	bytesRead := make([]byte, 1024) // 1KB
	for newLineIndex == -1 {
		_, err := file.Read(bytesRead)
		if err != nil {
			panic(err)
		}

		for i, char := range bytesRead {
			if string(char) == "\n" {
				newLineIndex = i + numberBytesRead
				break
			}
		}
		numberBytesRead += len(bytesRead) - 1
	}
	file.Seek(int64(newLineIndex-1), 0)
}

func main() {
	// Config
	config.LoadConfig("./shared/config-sample.yml")

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Initialize db
	database.InitDatabase()
	defer database.Pool.Close()

	skipFirstLine(file)

	fileChunk := make([]byte, chunkSize)
	for {
		bytesRead, err := file.Read(fileChunk)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if bytesRead == 0 || err == io.EOF {
			fmt.Printf("\nNo more file content\n")
			break
		}

		fileChunk = fileChunk[:bytesRead]
		// if the last byte from fileChunk isn't a new line character,
		// remove the content from fileChunk until previous new line
		if fileChunk[bytesRead-1] != byte('\n') && bytesRead == chunkSize {
			newLineIndex, err := findLatestChar(fileChunk, '\n')
			if err != nil {
				panic(err)
			}
			// go back to previous new line character in file seeker
			_, err = file.Seek(int64(-1*(len(fileChunk)-newLineIndex-1)), 1)
			if err != nil {
				panic(err)
			}
			fileChunk = fileChunk[:newLineIndex]
		}

		fileChunkStr := string(fileChunk)
		lines := strings.Split(fileChunkStr, "\n")

		var insertStatment strings.Builder
		insertStatment.Grow(len(fileChunkStr))
		insertStatment.WriteString("INSERT INTO public.file_data VALUES\n")

		batchSize := 0
		for _, line := range lines {
			// It's said, in instructions, that validation should be done only after persistence but
			// to conform with insert operation at least number of columns must be verified
			columnsData := strings.Fields(line)
			if len(columnsData) < 8 {
				// fmt.Printf("Skiping line: %s\n", line)
				continue
			}
			if batchSize > 0 {
				insertStatment.WriteRune(',')
			}
			insertStatment.WriteString(utils.GetInputLine(columnsData))

			batchSize++
		}
		insertStatment.WriteRune(';')
		// fmt.Println(insertStatment.String())

		database.Pool.MustExec(insertStatment.String())
	}
	database.CleanDataScripts()
}
