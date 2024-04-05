package routes

import (
	"database/sql"
	"dpacks-go-services-template/controllers"
	"dpacks-go-services-template/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutesFunc(r *gin.Engine, db *sql.DB) {
	// Create a new rate limiter middleware
	rateLimiter, err := middleware.NewRateLimit(db)
	if err != nil {
		// Handle error
		panic(err)
	}

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

		apiSubscribersRoutes := api.Group("/api_subscribers") // admin api subscriber  api group
		{
			apiSubscribersRoutes.POST("/subscriber", controllers.AddSubscribers(db))

			apiSubscribersRoutes.GET("/subscribers/:count/:page", controllers.GetApiSubscribers(db))
			apiSubscribersRoutes.GET("/subscriber/:id", controllers.GetApiSubscriberById(db))
			apiSubscribersRoutes.GET("/subscribers/datetime/:count/:page", controllers.GetApiSubscribersByDatetime(db))
			apiSubscribersRoutes.GET("/subscribers/datetime/count", controllers.GetApiSubscribersByDatetimeCount(db))
			apiSubscribersRoutes.GET("/subscribers/count", controllers.GetApiSubscribersCount(db))

			apiSubscribersRoutes.PUT("/subscriber/:id", controllers.RegenerateKey(db))

			apiSubscribersRoutes.DELETE("/subscriber/:id", controllers.DeleteApiSubscriberByID(db))
			apiSubscribersRoutes.DELETE("/subscriber/bulk/:id", controllers.DeleteApiSubscriberByIDBulk(db))

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
			keyPairsRoutes.GET("/", controllers.GetKeyPairs(db))         // get all keypairs
			keyPairsRoutes.GET("/:id", controllers.GetKeyPairsID(db))    // get keypair for the given user id
			keyPairsRoutes.POST("/:id", controllers.AddKeyPair(db))      // add keypair for the given user id
			keyPairsRoutes.PUT("/:id", controllers.UpdateKeyPair(db))    // update keypair for the given user id
			keyPairsRoutes.DELETE("/:id", controllers.DeleteKeyPair(db)) // delete keypair for the given user id
		}

		subscriptionPlansRoutes := api.Group("/subscription_plans") // subscription plans api group
		{
			subscriptionPlansRoutes.GET("/", controllers.GetSubscriptionPlans(db)) // get all subscription plans
		}

		templateRoutes := api.Group("/template") // template api group
		{
			templateRoutes.GET("/", controllers.GetTemplates(db)) // get all templates
		}

		visitorUserRoutes := api.Group("/visitor_user") // visitor user api group
		{
			visitorUserRoutes.GET("/", controllers.GetVisitorUsers(db)) // get all visitor users
		}
		rateLimitRouts := api.Group("/ratelimit") // visitor user api group
		{

			rateLimitRouts.POST("/addratelimit", controllers.AddRatelimit(db)) // add webpage

			rateLimitRouts.GET("/ratelimits/:count/:page", controllers.GetRateLimits(db))                    // get all webpages
			rateLimitRouts.GET("/ratelimit/:id", controllers.GetRatelimitById(db))                           // get a webpage by id
			rateLimitRouts.GET("/ratelimits/status/:count/:page", controllers.GetRatelimitsByStatus(db))     // get all webpages by status
			rateLimitRouts.GET("/ratelimits/status/count", controllers.GetRatelimitsByStatusCount(db))       // get all webpages by status
			rateLimitRouts.GET("/ratelimits/datetime/:count/:page", controllers.GetRatelimitsByDatetime(db)) // get all webpages by datetime
			rateLimitRouts.GET("/ratelimits/datetime/count", controllers.GetRatelimitsByDatetimeCount(db))   // get all webpages by datetime
			rateLimitRouts.GET("/ratelimits/count", controllers.GetRateLimitCount(db))                       // get all webpages count

			rateLimitRouts.PUT("/ratelimits/status/:id", controllers.UpdateRatelimitStatus(db))          // update webpage status by id
			rateLimitRouts.PUT("/ratelimits/:id", controllers.EditRatelimit(db))                         // edit webpage by id
			rateLimitRouts.PUT("/ratelimits/status/bulk/:id", controllers.UpdateRatelimitStatusBulk(db)) // update webpage status by id (bulk)
			//
			rateLimitRouts.DELETE("/ratelimits/:id", controllers.DeleteRatelimitByID(db))          // delete webpage by ID
			rateLimitRouts.DELETE("/ratelimits/bulk/:id", controllers.DeleteRatelimitByIDBulk(db)) // delete webpage by ID (bulk)
		}

		webContentRoutes := api.Group("/webcontent")
		//apply ratelimiter for webcontent subgrooup
		webContentRoutes.Use(rateLimiter.Limit()) //this also possible
		webContentRoutes.Use(middleware.AuthMiddleware(db))
		{
			webContentRoutes.GET("/webcontents", controllers.GetAllWebContents(db)) // get all webcontent
			webContentRoutes.GET("/webcontents/updated", controllers.GetUpdatedWebContents(db))
		}
	}
}
