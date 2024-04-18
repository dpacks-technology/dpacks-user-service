package middleware

import (
	"dpacks-go-services-template/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the JWT string from the header
		tokenString := c.GetHeader("Authorization")

		// api call to check if the token is valid
		url := os.Getenv("AUTH_API_HOST") + "/api/auth/check-auth"

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		req.Header.Add("Authorization", tokenString)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// return response data
		if resp.StatusCode == http.StatusOK {
			// Read the body content
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading body:", err)
				c.Abort()
				return
			}

			// Unmarshal the JSON response into a ResponseBody struct
			var responseBody models.AuthResponseBody
			err = json.Unmarshal(bodyBytes, &responseBody)
			if err != nil {
				fmt.Println("Error unmarshalling response:", err)
				c.Abort()
				return
			}

			// Set the values in the context
			c.Set("auth_userId", responseBody.ID)
			c.Set("auth_userKey", responseBody.UserKey)
			c.Set("auth_username", responseBody.Username)
			c.Set("auth_status", responseBody.Status)
			c.Set("auth_roles", responseBody.Role)
		} else if resp.StatusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
			c.Abort()
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error making the request"})
			c.Abort()
			return
		}

		// Proceed with the request
		c.Next()

	}
}
