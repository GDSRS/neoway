package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
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
			slog.Error("Error skiping first line: %s\n", err)
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

	// slog.SetLogLoggerLevel(slog.LevelDebug)
	// Config
	slog.Info("Loading configuration file")
	config.LoadConfig("./shared/config-sample.yml")

	// Open file
	slog.Info("Opening data file")
	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("Error opening base_teste.txt file", err)
		panic(err)
	}
	defer file.Close()

	// Initialize db
	slog.Info("Initializing database")
	database.InitDatabase()
	defer database.Pool.Close()

	slog.Info("Skiping header line")
	skipFirstLine(file)

	fileChunk := make([]byte, chunkSize)

	slog.Info("Reading file")
	for {
		bytesRead, err := file.Read(fileChunk)
		if err != nil && err != io.EOF {
			slog.Error("Error reading fileChunk: %s\n", err)
			panic(err)
		}
		if bytesRead == 0 || err == io.EOF {
			slog.Debug("No more file content")
			break
		}

		fileChunk = fileChunk[:bytesRead]
		// if the last byte from fileChunk isn't a new line character,
		// remove the content from fileChunk until previous new line
		if fileChunk[bytesRead-1] != byte('\n') && bytesRead == chunkSize {
			newLineIndex, err := findLatestChar(fileChunk, '\n')
			if err != nil {
				slog.Error("Error finding previous character: %s", err)
				panic(err)
			}
			// go back to previous new line character in file seeker
			_, err = file.Seek(int64(-1*(len(fileChunk)-newLineIndex-1)), 1)
			if err != nil {
				slog.Error("Error updating file seeker: %s", err)
				panic(err)
			}
			fileChunk = fileChunk[:newLineIndex]
		}

		fileChunkStr := string(fileChunk)
		lines := strings.Split(fileChunkStr, "\n")

		var insertStatment strings.Builder
		insertStatment.Grow(len(fileChunkStr))
		insertStatment.WriteString("INSERT INTO public.file_data VALUES\n")

		inserted := false
		for _, line := range lines {
			// It's said, in instructions, that validation should be done only after persistence but
			// to conform with insert operation at least number of columns must be verified
			columnsData := strings.Fields(line)
			if len(columnsData) < 8 {
				continue
			}

			inputLine, err := utils.GetInputLine(columnsData)
			if err != nil {
				slog.Info(fmt.Sprintf("Invalid input: %v skipping...", err))
				continue
			}

			if inserted {
				insertStatment.WriteRune(',')
			}

			insertStatment.WriteString(inputLine)

			inserted = true
		}
		if !inserted {
			continue
		}

		insertStatment.WriteRune(';')

		database.Pool.MustExec(insertStatment.String())
	}
	slog.Info("Running clean scripts in database")
	database.CleanDataScripts()
	slog.Info("Finishing execution")
}
