package controllers

import (
	"bytes"
	"database/sql"
	"dpacks-go-services-template/models"
	"dpacks-go-services-template/validators"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func GetVisitorInfo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		// get page id parameter
		page := c.Param("page")

		// get count parameter
		count := c.Param("count")

		// Convert page and count to integers
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		countInt, err := strconv.Atoi(count)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		// Calculate offset
		offset := (pageInt - 1) * countInt

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT * FROM visitor_info WHERE webid = $3 ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset, id) // Added 'id' to the args slice for SQL query

		if val != "" && key != "" {
			switch key {
			case "id":
				// Already filtering by ID, so no need to change the query
			case "device":
				query = "SELECT * FROM visitor_info WHERE webid = $3 AND device LIKE $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal) // Adjust the args slice accordingly
			case "country":
				query = "SELECT * FROM visitor_info WHERE webid = $3 AND country LIKE $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			}
		}

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with bound parameters
		rows, err := stmt.Query(args...)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into WebpageModel structs
		var visitorsInfo []models.VisitorInfo

		for rows.Next() {
			var visitorInfo models.VisitorInfo
			if err := rows.Scan(&visitorInfo.ID, &visitorInfo.IpAddres, &visitorInfo.Device, &visitorInfo.Country, &visitorInfo.Source, &visitorInfo.VisitedDate, &visitorInfo.WebId); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			visitorsInfo = append(visitorsInfo, visitorInfo)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, visitorsInfo)
	}
}

func GetVisitorInfoCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the visitor ID from the URL parameter
		id := c.Param("id")

		var count int

		// Get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Start building the query
		query := "SELECT COUNT(*) FROM visitor_info WHERE webid = $1"
		args = append(args, id)

		if val != "" && key != "" {
			// Depending on the key, modify the query to include an additional filter
			switch key {
			case "id":
				// No additional filter needed if key is 'id' since we're already filtering by id
			case "device":
				query += " AND device LIKE $2"
				args = append(args, escapedVal)
			case "country":
				query += " AND country LIKE $2"
				args = append(args, escapedVal)
			}
		}

		// Prepare the SQL statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing query"})
			return
		}
		defer stmt.Close()

		// Execute the query
		err = stmt.QueryRow(args...).Scan(&count)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing query"})
			return
		}

		//backend validation
		validators.NumberValidation(count)

		// Return the count as JSON
		c.JSON(http.StatusOK, count)
	}
}

func GetVisitorInfoByDatetime(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get page id parameter
		page := c.Param("page")

		// get count parameter
		count := c.Param("count")

		// Convert page and count to integers
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		countInt, err := strconv.Atoi(count)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		// Calculate offset
		offset := (pageInt - 1) * countInt

		// get query parameters
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT * FROM visitor_info ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT * FROM visitor_info WHERE visited_time BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT * FROM visitor_info WHERE id = $5 AND visited_time BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "name":
				query = "SELECT * FROM visitor_info WHERE device LIKE $5 AND visited_time BETWEEN $3 AND $4 ORDER BY CASE WHEN device = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "path":
				query = "SELECT * FROM visitor_info WHERE country LIKE $5 AND visited_time BETWEEN $3 AND $4 ORDER BY CASE WHEN country = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			}
		}

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		rows, err := stmt.Query(args...)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into WebpageModel structs
		var visitorsInfo []models.VisitorInfo

		for rows.Next() {
			var visitorInfo models.VisitorInfo
			if err := rows.Scan(&visitorInfo.ID, &visitorInfo.IpAddres, &visitorInfo.Device, &visitorInfo.Country, &visitorInfo.Source, &visitorInfo.VisitedDate); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			visitorsInfo = append(visitorsInfo, visitorInfo)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, visitorsInfo)

	}

}

func GetVisitorByDatetimeCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get query parameters
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM visitor_info"

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT COUNT(*) FROM visitor_info WHERE visited_time BETWEEN $1 AND $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM visitor_info WHERE id = $3 AND visited_time BETWEEN $1 AND $2"
				args = append(args, val)
			case "name":
				query = "SELECT COUNT(*) FROM visitor_info WHERE device LIKE $3 AND visited_time BETWEEN $1 AND $2"
				args = append(args, escapedVal)
			case "path":
				query = "SELECT COUNT(*) FROM visitor_info WHERE country LIKE $3 AND visited_time BETWEEN $1 AND $2"
				args = append(args, escapedVal)
			}
		}

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		var count int
		err = stmt.QueryRow(args...).Scan(&count)

		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Close the statement
		defer stmt.Close()

		// Return all webpages as JSON
		c.JSON(http.StatusOK, count)

	}

}

func GetVisitorInfoById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT * FROM visitor_info WHERE id = $1", id)

		// Create a WebpageModel to hold the data
		var visitorInfo models.VisitorInfo

		// Scan the row data into the WebpageModel
		err := row.Scan(&visitorInfo.ID, &visitorInfo.IpAddres, &visitorInfo.Device, &visitorInfo.Country, &visitorInfo.Source, &visitorInfo.VisitedDate)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the webpage as JSON
		c.JSON(http.StatusOK, visitorInfo)

	}

}
func GetSource(db *sql.DB) gin.HandlerFunc {
	// Return a handler function
	return func(c *gin.Context) {
		// Get the website ID from the URL params
		websiteID := c.Param("id")

		// Execute the SQL query
		rows, err := db.Query(`
            SELECT
                src.type AS user_source,
                COUNT(*) AS user_count
            FROM
                public.sessions AS s
            JOIN
                public.source AS src ON s.source_id = src.id
            WHERE
                s.web_id = $1
            GROUP BY
                s.web_id,
                src.type
            ORDER BY
                s.web_id,
                user_count DESC;
        `, websiteID)
		if err != nil {
			// Handle any errors
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}
		defer rows.Close()

		// Iterate over the result set
		var results []gin.H
		for rows.Next() {
			var userSource string
			var userCount int
			if err := rows.Scan(&userSource, &userCount); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning database rows"})
				return
			}
			// Append the result to the results slice
			results = append(results, gin.H{"user_source": userSource, "user_count": userCount})
		}

		// Return the results as JSON
		c.JSON(http.StatusOK, results)
	}
}

func GetSessions(db *sql.DB) gin.HandlerFunc {
	// Return a handler function
	return func(c *gin.Context) {
		// Get the website ID from the URL params
		websiteID := c.Param("id")

		// Execute the SQL query with the website ID as a parameter
		rows, err := db.Query(`
            SELECT
                EXTRACT(DOW FROM sessionstart) AS day_of_week,
                COUNT(*) AS session_count
            FROM
                public.sessions
            WHERE
                web_id = $1
            GROUP BY
                day_of_week
            ORDER BY
                day_of_week;
        `, websiteID)
		if err != nil {
			// Handle any errors
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}
		defer rows.Close()

		// Iterate over the result set
		var results []gin.H
		for rows.Next() {
			var dayOfWeek int
			var sessionCount int
			if err := rows.Scan(&dayOfWeek, &sessionCount); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning database rows"})
				return
			}
			// Append the result to the results slice
			results = append(results, gin.H{"day_of_week": dayOfWeek, "session_count": sessionCount})
		}

		// Return the results as JSON
		c.JSON(http.StatusOK, results)
	}
}
func GetDevices(db *sql.DB) gin.HandlerFunc {
	// Return a handler function
	return func(c *gin.Context) {
		// Get the website ID from the URL params
		websiteID := c.Param("id")

		// Execute the SQL query with the website ID as a parameter
		rows, err := db.Query(`
            SELECT
                d.devicename,
                COUNT(*) AS device_count
            FROM
                public.sessions AS s
            JOIN
                public.devices AS d ON s.deviceid = d.deviceid
            WHERE
                s.web_id IN ($1)
            GROUP BY
                s.web_id,
                d.devicename
            ORDER BY
                s.web_id,
                d.devicename;
        `, websiteID)
		if err != nil {
			// Handle any errors
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}
		defer rows.Close()

		// Iterate over the result set
		var results []gin.H
		for rows.Next() {
			var deviceName string
			var deviceCount int
			if err := rows.Scan(&deviceName, &deviceCount); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning database rows"})
				return
			}
			// Append the result to the results slice
			results = append(results, gin.H{"device_name": deviceName, "device_count": deviceCount})
		}

		// Return the results as JSON
		c.JSON(http.StatusOK, results)
	}
}

func GetCountry(db *sql.DB) gin.HandlerFunc {
	// Return a handler function
	return func(c *gin.Context) {
		// Get the website ID from the URL params
		websiteID := c.Param("id")

		// Execute the SQL query with the website ID as a parameter
		rows, err := db.Query(`
            SELECT
                c.countrycode,
                COUNT(DISTINCT s.ipaddress) AS user_count
            FROM
                public.sessions AS s
            JOIN
                public.countries AS c ON s.countrycode = c.countrycode
            WHERE
                s.web_id = $1
            GROUP BY
                s.web_id,
                c.countrycode
            ORDER BY
                s.web_id,
                user_count DESC;
        `, websiteID)
		if err != nil {
			// Handle any errors
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}
		defer rows.Close()

		// Iterate over the result set
		var results []gin.H
		for rows.Next() {
			var countryCode string
			var userCount int
			if err := rows.Scan(&countryCode, &userCount); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning database rows"})
				return
			}
			// Append the result to the results slice
			results = append(results, gin.H{"country_code": countryCode, "user_count": userCount})
		}

		// Return the results as JSON
		c.JSON(http.StatusOK, results)
	}
}

func CreateNewAlert(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		// get the JSON data
		var alert models.CreateNewUserAlert
		if err := c.ShouldBindJSON(&alert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := "INSERT INTO useralerts (alert_threshold, alert_subject,alert_content,when_alert_required,website_id) VALUES ($1, $2, $3, $4,$5)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		fmt.Printf("test3")

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(alert.AlertThreshold, alert.AlertSubject, alert.AlertContent, alert.WhenAlertRequired, alert.WebsiteeId)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		fmt.Printf("test4")

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"message": "Alert set Succesfully"})

	}
}

func GetAllAlert(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		//get id parameter
		id := c.Param("id")

		// get page id parameter
		page := c.Param("page")

		// get count parameter
		count := c.Param("count")

		// Convert page and count to integers
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		countInt, err := strconv.Atoi(count)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		// Calculate offset
		offset := (pageInt - 1) * countInt

		//get query parameters
		key := c.Query("key")
		val := c.Query("val")
		//escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT id,alert_threshold,alert_subject,alert_content,when_alert_required,status,website_id FROM useralerts WHERE website_id=$3 ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset, id)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT * FROM useralerts WHERE id = $3  ORDER BY id LIMIT $1 OFFSET $2"
				fmt.Printf("%s", args)
				args = append(args, val)

			}
		}

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {

			fmt.Printf("%s\n", err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with bound parameters
		rows, err := stmt.Query(args...)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		defer rows.Close()

		//close the rows when the surrounding function returns(handler function)

		// Iterate over the rows and scan them into WebpageModel structs
		var alerts []models.UserAlertsShow
		for rows.Next() {
			var alert models.UserAlertsShow
			if err := rows.Scan(&alert.AlertID, &alert.AlertThreshold, &alert.AlertSubject, &alert.AlertContent, &alert.WhenAlertRequired, &alert.Status, &alert.WebsiteeId); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			alerts = append(alerts, alert)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, alerts)

	}

}

func GetAlertsCount(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		var count int
		id := c.Param("id")

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM useralerts WHERE website_id = $1"
		args = append(args, id)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM useralerts WHERE id = $1 AND website_id = $2"
				args = append(args, val, id)
			case "alertthreshold":
				query = "SELECT COUNT(*) FROM useralerts WHERE alertthreshold LIKE $1 AND website_id = $2"
				args = append(args, escapedVal, id)
			}
		}

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		err = stmt.QueryRow(args...).Scan(&count)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Close the statement
		defer stmt.Close()

		//backend validation
		validators.NumberValidation(count)

		// Return all webpages as JSON
		c.JSON(http.StatusOK, count)

	}

}

func GetAlertbyId(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT id,alert_threshold,alert_subject,alert_content,when_alert_required,status,website_id FROM useralerts WHERE id = $1", id)

		// Create a WebpageModel to hold the data
		var alert models.UserAlertsShow

		// Scan the row data into the WebpageModel
		err := row.Scan(&alert.AlertID, &alert.AlertThreshold, &alert.AlertSubject, &alert.AlertContent, &alert.WhenAlertRequired, &alert.Status, &alert.WebsiteeId)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the webpage as JSON
		c.JSON(http.StatusOK, alert)

	}

}

func GetAlertsByStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get status parameter (array)
		statuses := c.Query("status")

		// get page id parameter
		page := c.Param("page")

		// get count parameter
		count := c.Param("count")

		// Convert page and count to integers
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		countInt, err := strconv.Atoi(count)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		// Calculate offset
		offset := (pageInt - 1) * countInt

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT * FROM useralerts ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		switch statuses {
		case "1":
			query = "SELECT * FROM useralerts WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 1)
		case "0":
			query = "SELECT * FROM useralerts WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT * FROM useralerts WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				query = "SELECT * FROM useralerts WHERE id = $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "alertthreshold":
				query = "SELECT * FROM useralerts WHERE alertthreshold LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "whenalertrequired":
				query = "SELECT * FROM useralerts WHERE whenalertrequired LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			}
		}

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		rows, err := stmt.Query(args...)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into WebpageModel structs
		var alerts []models.UserAlertsShow

		for rows.Next() {
			var alert models.UserAlertsShow
			if err := rows.Scan(&alert.AlertID, &alert.AlertThreshold, &alert.AlertSubject, &alert.AlertContent, &alert.WhenAlertRequired, &alert.Status, &alert.WebsiteeId); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			alerts = append(alerts, alert)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, alerts)

	}

}
func GetAlertsByStatusCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get status parameter (array)
		statuses := c.Query("status")

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM useralerts"

		switch statuses {
		case "1":
			query = "SELECT COUNT(*) FROM useralerts WHERE status IN ($1)"
			args = append(args, 1)
		case "0":
			query = "SELECT COUNT(*) FROM useralerts WHERE status IN ($1)"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM useralerts WHERE id = $2 AND status IN ($1)"
				args = append(args, val)
			case "alertthreshold":
				query = "SELECT COUNT(*) FROM useralerts WHERE alertthreshold LIKE $2 AND status IN ($1)"
				args = append(args, escapedVal)
			case "whenalertrequired":
				query = "SELECT COUNT(*) FROM useralerts WHERE whenalertrequired LIKE $2 AND status IN ($1)"
				args = append(args, escapedVal)
			}
		}

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		var count int
		err = stmt.QueryRow(args...).Scan(&count)

		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Close the statement
		defer stmt.Close()

		//backend validation
		validators.NumberValidation(count)

		// Return all webpages as JSON
		c.JSON(http.StatusOK, count)

	}

}

//func EditAlert(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		// get id parameter
//		id := c.Param("id")
//
//		// get the JSON data - only the name
//		var alert models.UserAlertsModel
//		if err := c.ShouldBindJSON(&alert); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//
//		//// Validate the webpage data
//		//if err := validators.ValidateName(alert, false); err != nil {
//		//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		//	return
//		//}
//
//		// Update the webpage in the database
//		_, err := db.Exec("UPDATE useralerts SET name = $1 WHERE id = $2", alert.Name, id)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		// Return a success message
//		c.JSON(http.StatusOK, gin.H{"message": "Webpage updated successfully"})
//
//	}
//
//}

func DeleteAlertByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// query to delete the webpage
		query := "DELETE FROM useralerts WHERE id = $1"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Alert deleted successfully"})

	}

}

func DeleteAlertByIDBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get ids array as a parameter as integer
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// Delete the webpage from the database
		for _, id := range ids {
			// query to delete the webpage
			query := "DELETE FROM useralerts WHERE id = $1"

			// Prepare the statement
			stmt, err := db.Prepare(query)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			// Execute the prepared statement with bound parameters
			_, err = stmt.Exec(id)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Alert bulk deleted successfully"})

	}

}

func UpdateAlertStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the status
		var alert models.UserAlertStatus
		if err := c.ShouldBindJSON(&alert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to update the webpage status
		query := "UPDATE useralerts SET status = $1 WHERE id = $2"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(alert.Status, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Webpage status updated successfully"})

	}

}

func UpdateAlertStatusBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// get the JSON data - only the status
		var alert models.UserAlertStatus
		if err := c.ShouldBindJSON(&alert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the webpage status in the database
		for _, id := range ids {

			query := "UPDATE useralerts SET status = $1 WHERE id = $2"

			// Prepare the statement
			stmt, err := db.Prepare(query)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			// Execute the prepared statement with bound parameters
			_, err = stmt.Exec(alert.Status, id)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Webpage status updated successfully"})

	}
}

func EditAlert(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the name
		var updateAlert models.UpdateUserAlert
		if err := c.ShouldBindJSON(&updateAlert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("test %s", updateAlert)

		// Update the webpage in the database
		_, err := db.Exec("UPDATE useralerts SET alert_threshold=$1,alert_subject=$2,alert_content=$3,when_alert_required=$4 where id=$5", updateAlert.AlertThreshold, updateAlert.AlertSubject, updateAlert.AlertContent, updateAlert.WhenAlertRequired, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Webpage updated successfully"})

	}
}

func SessionRecord(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var session models.SessionRecord
		if err := c.ShouldBindJSON(&session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//backend validation
		validators.StringValidation(session.CountryCode)
		validators.NumberValidation(session.DeviceId)
		validators.NumberValidation(session.SourceId)

		query := "INSERT INTO sessions (sessionid,ipaddress,countrycode,deviceid,source_id,landingpage,web_id) VALUES ($1, $2, $3, $4,$5,$6,$7)"
		stmt, err := db.Prepare(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(session.SessionID, session.IpAddress, session.CountryCode, session.DeviceId, session.SourceId, session.LandingPage, session.WebId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		countQuery := "SELECT count FROM session_count WHERE website_id = $1"
		var count int
		err = db.QueryRow(countQuery, session.WebId).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		alertQuery := "SELECT id, alert_threshold FROM useralerts WHERE status = 1 and website_id = $1"
		rows, err := db.Query(alertQuery, session.WebId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var alertsToUpdate []struct {
			ID             int
			AlertThreshold int
		}

		for rows.Next() {
			var alert struct {
				ID             int
				AlertThreshold int
			}
			err := rows.Scan(&alert.ID, &alert.AlertThreshold)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			alertsToUpdate = append(alertsToUpdate, alert)
		}

		fmt.Printf("Count: %d\n", count)

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		quotaExceeded := false
		for _, alert := range alertsToUpdate {
			if count >= alert.AlertThreshold {
				quotaExceeded = true
				// Update the status of the user alert record to 0
				updateQuery := "UPDATE useralerts SET status = 0 WHERE id = $1"
				_, err := db.Exec(updateQuery, alert.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}
		fmt.Printf("Quota Exceeded: %t\n", quotaExceeded)

		//get the user email
		if quotaExceeded {

			fmt.Printf("hellow test 1")
			// Retrieve the email of the users whose quota has been exceeded
			emailQuery := "SELECT email FROM users WHERE id IN (SELECT user_id FROM user_site WHERE site_id = $1)"
			rows, err := db.Query(emailQuery, session.WebId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			defer rows.Close()

			// Send email to each user
			for rows.Next() {
				var email string
				err := rows.Scan(&email)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				// Send email using the provided endpoint
				sendEmail(email)
			}
			fmt.Printf("hellow test 3")

		}

		c.JSON(http.StatusCreated, gin.H{"message": "Added session record successfully"})
	}
}

func sendEmail(email string) {
	payload := map[string]interface{}{

		//get from env file api key from
		"api_key": "ctn47m4o8mSwqo8Fgr89gwrQorhoyDrhp9qHgtp9tSgotiSmDyxepGulHoaUPalzxu0zw0peqwdk9tuc",
		"to":      email,
		"subject": "User Count Limit Exceeded",
		"message": "User Count is Reached the Limit to your website",
		"size":    "sm",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf(string(jsonData))

	resp, err := http.Post("http://34.47.130.27:4005/api/email/send", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to send email. Status code: %d\n", resp.StatusCode)
	}
}
