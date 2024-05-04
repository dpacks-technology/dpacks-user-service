package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func UploadTemplate(filename string, file *multipart.FileHeader) error {
	// Open the uploaded file
	uploadedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	// Create a buffer to store the file content
	fileContent := bytes.Buffer{}
	_, err = io.Copy(&fileContent, uploadedFile)
	if err != nil {
		return err
	}

	// Send multipart request to storage microservice
	url := os.Getenv("STORAGE_MICROSERVICE_HOST") + "/template" // URL of the storage microservice

	// Create a new form data buffer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file to the form data
	fileWriter, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(fileWriter, &fileContent)
	if err != nil {
		return err
	}

	// Add filename as a form field
	err = writer.WriteField("filename", filename)
	if err != nil {
		return err
	}

	// Close the form data writer
	writer.Close()

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType()) // Set the content type

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// print the response
	fmt.Println(resp)

	// Check response status
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("%s\n", err)
		return err
	}

	return nil
}
