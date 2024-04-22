package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"dpacks-go-services-template/validators"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AddSite handles POST /api/web/site - CREATE
func AddSite(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data
		// get the JSON data
		var site models.Site
		if err := c.ShouldBindJSON(&site); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the webpage data
		if err := validators.ValidateSite(site, true); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// insert data
		query := `INSERT INTO sites (name, description, category, domain, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// execute the statement
		err = stmt.QueryRow(site.Name, site.Description, site.Category, site.Domain, 1).Scan(&site.ID)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		query = `INSERT INTO user_site (user_id, site_id) VALUES ($1, $2)`

		// Prepare the statement
		stmt, err = db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// get the authenticated user data
		userid, _ := c.Get("auth_userId")

		// execute the statement
		_, err = stmt.Exec(userid, site.ID)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// return the response
		c.JSON(http.StatusCreated, gin.H{"message": "Webpage added successfully"})

	}
}

// ReadSites handles GET /api/web/sites - READ
func ReadSites(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get the authenticated user data
		userid, _ := c.Get("auth_userId")
		userKey, _ := c.Get("auth_userKey")
		username, _ := c.Get("auth_username")
		status, _ := c.Get("auth_status")
		roles, _ := c.Get("auth_roles")

		// print the user data
		fmt.Printf("User data: %v, %v, %v, %v, %v\n", userid, userKey, username, status, roles)

		// get data
		query := `SELECT id, name, description, category, domain, status, last_updated FROM sites, user_site WHERE user_site.site_id = sites.id AND user_site.user_id = $1 ORDER BY seq_id`

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// execute the statement
		rows, err := stmt.Query(userid)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// iterate over the rows
		var sites []models.Site
		for rows.Next() {
			var site models.Site
			err := rows.Scan(&site.ID, &site.Name, &site.Description, &site.Category, &site.Domain, &site.Status, &site.LastUpdated)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			sites = append(sites, site)
		}

		// return the response
		c.JSON(http.StatusOK, sites)

	}
}

// GetSiteById handles GET /api/web/site/:id - READ
func GetSiteById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get the site id
		siteId := c.Param("id")

		// get data
		query := `SELECT name, domain, description, category, status, last_updated FROM sites WHERE id = $1`

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// execute the statement
		var site models.Site
		err = stmt.QueryRow(siteId).Scan(&site.Name, &site.Domain, &site.Description, &site.Category, &site.Status, &site.LastUpdated)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// return the response
		c.JSON(http.StatusOK, site)

	}
}

// EditSite handles PUT /api/web/site/:id - UPDATE
func EditSite(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get the site id
		siteId := c.Param("id")

		// get data
		// get the JSON data
		var site models.Site
		if err := c.ShouldBindJSON(&site); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the webpage data
		if err := validators.ValidateSite(site, false); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// update data
		query := `UPDATE sites SET name = $1, description = $2, category = $3 WHERE id = $4`

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// execute the statement
		_, err = stmt.Exec(site.Name, site.Description, site.Category, siteId)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// return the response
		c.JSON(http.StatusOK, gin.H{"message": "Webpage updated successfully"})

	}
}

// DeleteSite handles DELETE /api/web/site/:id - DELETE
func DeleteSite(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get the site id
		siteId := c.Param("id")

		// delete data
		query := `DELETE FROM sites WHERE id = $1`

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// execute the statement
		_, err = stmt.Exec(siteId)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// return the response
		c.JSON(http.StatusOK, gin.H{"message": "Webpage deleted successfully"})

	}
}
