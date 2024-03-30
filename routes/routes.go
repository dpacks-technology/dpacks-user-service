package routes

import (
	"database/sql"
	"dpacks-go-services-template/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutesFunc(r *gin.Engine, db *sql.DB) {
	api := r.Group("/api")
	{
		exampleRoutes := api.Group("/example") // example api group
		{
			exampleRoutes.GET("/", controllers.GetExample(db))               // get all examples
			exampleRoutes.GET("/:id", controllers.GetExampleByID(db))        // get example by ID
			exampleRoutes.POST("/", controllers.AddExample(db))              // add example
			exampleRoutes.PUT("/:id", controllers.UpdateExample(db))         // update example by id
			exampleRoutes.PUT("/bulk", controllers.UpdateExampleBulk(db))    // update examples (bulk) by id
			exampleRoutes.DELETE("/:id", controllers.DeleteExample(db))      // update examples (bulk) by id
			exampleRoutes.DELETE("/bulk", controllers.DeleteExampleBulk(db)) // update examples (bulk) by id
		}

		webRoutes := api.Group("/web") // web api group
		{
			webRoutes.POST("/webpage", controllers.AddWebPage(db)) // add webpage

			webRoutes.GET("/webpages/:count/:page", controllers.GetWebPages(db))                    // get all webpages
			webRoutes.GET("/webpage/:id", controllers.GetWebPageById(db))                           // get a webpage by id
			webRoutes.GET("/webpages/status/:count/:page", controllers.GetWebPagesByStatus(db))     // get all webpages by status
			webRoutes.GET("/webpages/status/count", controllers.GetWebPagesByStatusCount(db))       // get all webpages by status
			webRoutes.GET("/webpages/datetime/:count/:page", controllers.GetWebPagesByDatetime(db)) // get all webpages by datetime
			webRoutes.GET("/webpages/datetime/count", controllers.GetWebPagesByDatetimeCount(db))   // get all webpages by datetime
			webRoutes.GET("/webpages/count", controllers.GetWebPagesCount(db))                      // get all webpages count

			webRoutes.PUT("/webpages/status/:id", controllers.UpdateWebPageStatus(db))          // update webpage status by id
			webRoutes.PUT("/webpages/:id", controllers.EditWebPage(db))                         // edit webpage by id
			webRoutes.PUT("/webpages/status/bulk/:id", controllers.UpdateWebPageStatusBulk(db)) // update webpage status by id (bulk)

			webRoutes.DELETE("/webpages/:id", controllers.DeleteWebPageByID(db))          // delete webpage by ID
			webRoutes.DELETE("/webpages/bulk/:id", controllers.DeleteWebPageByIDBulk(db)) // delete webpage by ID (bulk)
		}

		adminUserRoutes := api.Group("/admin_user") // admin user api group
		{
			adminUserRoutes.GET("/", controllers.GetAdminUsers(db)) // get all admin users
		}

		autoRespondRoutes := api.Group("/auto_respond") // auto respond api group
		{
			autoRespondRoutes.GET("/", controllers.GetAutoResponds(db)) // get all auto responds
		}

		analyticalAlertsRoutes := api.Group("/analytical_alerts") // analytical alerts api group
		{
			analyticalAlertsRoutes.GET("/", controllers.GetAnalyticalAlerts(db)) // get all analytical alerts
		}

		keyPairsRoutes := api.Group("/keypairs") // keypairs api group
		{
			keyPairsRoutes.GET("/", controllers.GetKeyPairs(db)) // get all keypairs
		}

		subscriptionPlansRoutes := api.Group("/subscription_plans") // subscription plans api group
		{
			subscriptionPlansRoutes.GET("/", controllers.GetSubscriptionPlans(db)) // get all subscription plans
		}

		templateRoutes := api.Group("/marketplace") // marketplace api group
		{
			templateRoutes.GET("/templates/:count/:page", controllers.GetTemplates(db)) // get all
			//templateRoutes.GET("/webpage/:id", controllers.GetWebPageById(db))                           // get by id
			templateRoutes.GET("/templates/status/:count/:page", controllers.GetTemplatesByStatus(db))     // get all by status
			templateRoutes.GET("/templates/status/count", controllers.GetTemplatesByStatusCount(db))       // get all by status
			templateRoutes.GET("/templates/datetime/:count/:page", controllers.GetTemplatesByDatetime(db)) // get all by datetime
			templateRoutes.GET("/templates/datetime/count", controllers.GetTemplatesByDatetimeCount(db))   // get all by datetime
			templateRoutes.GET("/templates/count", controllers.GetTemplatesCount(db))                      // get all count

			templateRoutes.PUT("/templates/status/:id", controllers.UpdateTemplatesStatus(db)) // update status by id
			//templateRoutes.PUT("/templates/:id", controllers.EditWebPage(db))                         // edit by id
			templateRoutes.PUT("/templates/status/bulk/:id", controllers.UpdateTemplatesStatusBulk(db)) // update status by id (bulk)

			templateRoutes.DELETE("/templates/:id", controllers.DeleteTemplateByID(db))          // delete by ID
			templateRoutes.DELETE("/templates/bulk/:id", controllers.DeleteTemplateByIDBulk(db)) // delete by ID (bulk)

		}

		visitorUserRoutes := api.Group("/visitor_user") // visitor user api group
		{
			visitorUserRoutes.GET("/", controllers.GetVisitorUsers(db)) // get all visitor users
		}
	}
}
