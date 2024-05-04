package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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

		limit := c.Param("limit")
		//this is dummy site id used to demonstrate our system is can have 1m data at time
		query := `SELECT dp.site, dp.page, dp.element, st.domain
					from data_packets dp
					INNER JOIN sites st ON dp.site = st.id::text
					WHERE site NOT IN ('6fa3aa36-04b5-4a1d-b426-d3c76d87ff12')
					GROUP BY dp.site, dp.page, dp.element, st.domain, dp.last_updated
					ORDER BY dp.last_updated DESC LIMIT $1`

		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		defer stmt.Close()

		rows, err := stmt.Query(limit)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		defer rows.Close()

		var updatedData []json.RawMessage
		for rows.Next() {
			var SiteID string
			var Page string
			var Element string
			var Domain string
			err := rows.Scan(&SiteID, &Page, &Element, &Domain)
			if err != nil {
				fmt.Printf("%s\n", " happend here")
				return
			}
			//dataPackets = append(dataPackets, dataPacket)

			url := fmt.Sprintf("https://web.dpacks.net/api/v1/data-packets/fetch/%s/%s/%s", SiteID, Page, Element)

			response, err := doGetRequest(url)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			var data json.RawMessage
			err = json.Unmarshal(response, &data)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			// Create a map to hold domain name and data
			domainData := map[string]interface{}{
				"domain": Domain,
				"data":   data,
			}

			// Marshal the map into JSON
			jsonData, err := json.Marshal(domainData)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			updatedData = append(updatedData, jsonData)
		}

		c.JSON(200, updatedData)
		//c.JSON(200, dataPackets)

	}
}

func doGetRequest(url string) ([]byte, error) {

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
