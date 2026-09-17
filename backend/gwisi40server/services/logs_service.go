package services

import (
	"archive/zip"
	"gwisi40server/models"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func LoadLogFiles() []models.Logfile {

	pathSeparator := string(os.PathSeparator)

	path := ".." + pathSeparator + ".." + pathSeparator + "gateway"

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var result []models.Logfile

	logId := uint(0)

	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".log.gz") || e.Name() == "gateway.log" {
			info, _ := e.Info()
			abs, _ := filepath.Abs(path + pathSeparator + e.Name())

			log := models.Logfile{
				ID:       logId,
				FullName: abs,
				Name:     e.Name(),
				Date:     info.ModTime().Format("02/01/2006"),
				Time:     info.ModTime().Format("15:04:05"),
			}
			result = append(result, log)
			logId++
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date > result[j].Date
	})

	return result
}

func ZipLogFiles() string {
	// Create a new ZIP file
	zipFileName := "gateway.zip"
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer zipFile.Close()

	// Create a new ZIP writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// List of files to add to the ZIP
	entries := LoadLogFiles()

	for _, e := range entries {
		file, err := os.Open(e.Name)

		if err != nil {
			log.Println(err)
			continue
		}

		defer file.Close()

		// Create a new file entry in the ZIP
		fileInfo, err := file.Stat()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("File name: %s, size %d\n", fileInfo.Name(), fileInfo.Size())

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			log.Println(err)
			continue
		}

		header.Name = e.Name

		// Add the file to the ZIP
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			log.Println(err)
			continue
		}

		// Copy the file data into the ZIP
		_, err = io.Copy(writer, file)

		if err != nil {
			log.Println(err)
		}
	}

	log.Println("Logs ZIP file created successfully:", zipFileName)

	return zipFileName
}
