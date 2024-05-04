package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAllDpacksSites function
func GetAllDpacksSites(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		query := `SELECT id,name,domain,category from sites`

		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		defer stmt.Close()

		rows, err := stmt.Query()
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		defer rows.Close()

		var sites []models.Site
		for rows.Next() {
			var site models.Site
			err := rows.Scan(&site.ID, &site.Name, &site.Domain, &site.Category)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			sites = append(sites, site)
		}
		// Return a message to display this endpoint is working
		//c.JSON(200, gin.H{"message": "GetAllWebContents endpoint is working!!!!!!!!!!"})
		c.JSON(http.StatusOK, sites)

	}
}

// GetUpdatedWebContents function
func GetUpdatedWebContents(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Return a message to display this endpoint is working
		c.JSON(200, gin.H{"message": "GetUpdatedWebContents endpoint is working!!!!!!!!!!"})

	}
}
