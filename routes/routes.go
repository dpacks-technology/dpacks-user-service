package routes

import (
	"database/sql"
	"dpacks-go-services-template/controllers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(r *gin.Engine, db *sql.DB) {
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
			webRoutes.GET("/pages/", controllers.GetWebPages(db)) // get all webpages
		}

		adminUserRoutes := api.Group("/admin_user") // admin user api group
		{
			adminUserRoutes.GET("/", controllers.GetAdminUsers(db)) // get all admin users
		}

		//autoRespondRoutes := api.Group("/auto_respond") // auto respond api group
		//{
		//	autoRespondRoutes.GET("/", controllers.GetAutoResponds(db)) // get all auto responds
		//}
		//
		//analyticalAlertsRoutes := api.Group("/analytical_alerts") // analytical alerts api group
		//{
		//	analyticalAlertsRoutes.GET("/", controllers.GetAnalyticalAlerts(db)) // get all analytical alerts
		//}
		//
		//keyPairsRoutes := api.Group("/keypairs") // keypairs api group
		//{
		//	keyPairsRoutes.GET("/", controllers.GetKeyPairs(db)) // get all keypairs
		//}
		//
		//subscriptionPlansRoutes := api.Group("/subscription_plans") // subscription plans api group
		//{
		//	subscriptionPlansRoutes.GET("/", controllers.GetSubscriptionPlans(db)) // get all subscription plans
		//}
		//
		//templateRoutes := api.Group("/template") // template api group
		//{
		//	templateRoutes.GET("/", controllers.GetTemplates(db)) // get all templates
		//}
		//
		//visitorUserRoutes := api.Group("/visitor_user") // visitor user api group
		//{
		//	visitorUserRoutes.GET("/", controllers.GetVisitorUsers(db)) // get all visitor users
		//}
	}
}
