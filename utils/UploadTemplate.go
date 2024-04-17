package utils

//import (
//	"github.com/gin-gonic/gin"
//	"io"
//	"net/http"
//	"os"
//)
//
//func UploadTemplate(filename string, file) {
//
//// Send multipart request to storage microservice
//url := os.Getenv("STORAGE_MICROSERVICE_HOST") + "/template" // URL of the storage microservice
//req, err := http.NewRequest("POST", url, body)           // Create a new HTTP request
//if err != nil {
//c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create HTTP request"})
//return
//}
//req.Header.Set("Content-Type", writer.FormDataContentType()) // Set the content type
//
//// Send the request
//resp, err := http.DefaultClient.Do(req)
//if err != nil {
//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//return
//}
//
//// Close the response body
//defer func(Body io.ReadCloser) {
//	err := Body.Close()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close response body"})
//		return
//	}
//}(resp.Body)
//
//// Check response status
//if resp.StatusCode != http.StatusOK {
//c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write data packet"})
//return
//}
