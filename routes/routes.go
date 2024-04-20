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
			webRoutes.POST("/site", controllers.AddSite(db))                                    // add site
			webRoutes.GET("/sites", middleware.UserAuthMiddleware(), controllers.ReadSites(db)) // read all sites
			webRoutes.GET("/site/:id", controllers.GetSiteById(db))                             // read site by id
			webRoutes.PUT("/site/:id", controllers.EditSite(db))                                // edit site by id
			webRoutes.DELETE("/site/:id", controllers.DeleteSite(db))                           // delete site by id

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

			adminUserRoutes.POST("/addAdmin", controllers.AddAdminUser(db)) // add admin

			adminUserRoutes.GET("/admins/:count/:page", controllers.GetAdmins(db))                    // get all admins
			adminUserRoutes.GET("/admin/:id", controllers.GetAdminById(db))                           // get an admin by id
			adminUserRoutes.GET("/admins/status/:count/:page", controllers.GetAdminsByStatus(db))     // get all admins by status
			adminUserRoutes.GET("/admins/status/count", controllers.GetAdminsByStatusCount(db))       // get all admins by status
			adminUserRoutes.GET("/admins/datetime/:count/:page", controllers.GetAdminsByDatetime(db)) // get all admins by datetime
			adminUserRoutes.GET("/admins/datetime/count", controllers.GetAdminsByDatetimeCount(db))   // get all admins by datetime count
			adminUserRoutes.GET("/admins/count", controllers.GetAdminsCount(db))                      // get all admins count

			adminUserRoutes.PUT("/admins/status/:id", controllers.UpdateAdminStatus(db))          // update admin status by id
			adminUserRoutes.PUT("/admins/:id", controllers.EditAdmin(db))                         // edit admin by id
			adminUserRoutes.PUT("/admins/status/bulk/:id", controllers.UpdateAdminStatusBulk(db)) // update admin status by id (bulk)

			adminUserRoutes.DELETE("/admins/:id", controllers.DeleteAdminByID(db))          // delete admin by ID
			adminUserRoutes.DELETE("/admins/bulk/:id", controllers.DeleteAdminByIDBulk(db)) // delete admin by ID (bulk)

		}

		adminSitesRoutes := api.Group("/admin_sites") // admin site user api group
		{

			adminSitesRoutes.GET("/sites/:count/:page", controllers.GetSites(db))                    // get all sites
			adminSitesRoutes.GET("/sites/status/:count/:page", controllers.GetSitesByStatus(db))     // get all site by status
			adminSitesRoutes.GET("/sites/status/count", controllers.GetSitesByStatusCount(db))       // get all site by status
			adminSitesRoutes.GET("/sites/datetime/:count/:page", controllers.GetSitesByDatetime(db)) // get all site by datetime
			adminSitesRoutes.GET("/sites/datetime/count", controllers.GetSitesByDatetimeCount(db))   // get all site by datetime count
			adminSitesRoutes.GET("/sites/count", controllers.GetSitesCount(db))                      // get all site count

			adminSitesRoutes.PUT("/sites/status/:id", controllers.UpdateSiteStatus(db))           // update site status by id
			adminSitesRoutes.PUT("/sites/status/bulk/:id", controllers.UpdateSitesStatusBulk(db)) // update site status by id (bulk)

		}

		userRoutes := api.Group("/user") // user api group
		{
			userRoutes.GET("/users/:count/:page", controllers.GetUsers(db)) // get all users
			//userRoutes.GET("/user/:id", controllers.GetUserById(db))                        // get a user by id
			userRoutes.GET("/users/status/:count/:page", controllers.GetUsersByStatus(db))     // get all users by status
			userRoutes.GET("/users/status/count", controllers.GetUsersByStatusCount(db))       // get all users by status
			userRoutes.GET("/users/datetime/:count/:page", controllers.GetUsersByDatetime(db)) // get all users by datetime
			userRoutes.GET("/users/datetime/count", controllers.GetUsersByDatetimeCount(db))   // get all users by datetime
			userRoutes.GET("/users/count", controllers.GetUsersCount(db))                      // get all users count

			userRoutes.PUT("/users/status/:id", controllers.UpdateUserStatus(db))          // update user status by id
			userRoutes.PUT("/users/status/bulk/:id", controllers.UpdateUserStatusBulk(db)) // update users status by id (bulk)

			userRoutes.DELETE("/users/:id", controllers.DeleteUserByID(db))          // delete user by ID
			userRoutes.DELETE("/users/bulk/:id", controllers.DeleteUserByIDBulk(db)) // delete user by ID (bulk)
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

		autoRespondRoutes := api.Group("/chat") // auto respond api group
		{
			autoRespondRoutes.GET("/auto_respond/:count/:page", controllers.GetAutoResponds(db)) // get all auto responds
			autoRespondRoutes.POST("/auto_respond", controllers.AddAutoRespond(db))
			autoRespondRoutes.GET("/auto_respond/id/:id", controllers.GetAutoRespondsById(db))                      // get a webpage by id
			autoRespondRoutes.GET("/auto_respond/status/:count/:page", controllers.GetAutoRespondsByStatus(db))     // get all webpages by status
			autoRespondRoutes.GET("/auto_respond/status/count", controllers.GetAutoRespondsByStatusCount(db))       // get all webpages by status
			autoRespondRoutes.GET("/auto_respond/datetime/:count/:page", controllers.GetAutoRespondsByDatetime(db)) // get all webpages by datetime
			autoRespondRoutes.GET("/auto_respond/datetime/count", controllers.GetAutoRespondsByDatetimeCount(db))   // get all webpages by datetime
			autoRespondRoutes.GET("/auto_respond/count", controllers.GetAutoRespondsCount(db))                      // get all webpages count

			autoRespondRoutes.PUT("/auto_respond/status/:id", controllers.UpdateAutoRespondsStatus(db)) // update webpage status by id
			autoRespondRoutes.PUT("/auto_respond/:id", controllers.EditAutoResponds(db))                // edit webpage by id
			autoRespondRoutes.PUT("/auto_respond/status/bulk/:id", controllers.UpdateAutoRespondsStatusBulk(db))

			autoRespondRoutes.DELETE("/auto_respond/:id", controllers.DeleteAutoRespondsID(db))            // delete webpage by ID
			autoRespondRoutes.DELETE("/auto_respond/bulk/:id", controllers.DeleteAutoRespondsByIDBulk(db)) // delete webpage by ID (bulk)
		}

		analyticalAlertsRoutes := api.Group("/analytical_alerts") // analytical alerts api group
		{
			analyticalAlertsRoutes.GET("/", controllers.GetAnalyticalAlerts(db)) // get all analytical alerts

			analyticalAlertsRoutes.GET("/source/:id", controllers.GetSource(db))
			analyticalAlertsRoutes.GET("/sessions/:id", controllers.GetSessions(db))
			analyticalAlertsRoutes.GET("/devices/:id", controllers.GetDevices(db))
			analyticalAlertsRoutes.GET("/country/:id", controllers.GetCountry(db))
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

			rateLimitRouts.POST("/addratelimit", controllers.AddRatelimit(db))

			rateLimitRouts.GET("/ratelimits/:count/:page", controllers.GetRateLimits(db))
			rateLimitRouts.GET("/ratelimit/:id", controllers.GetRatelimitById(db))
			rateLimitRouts.GET("/ratelimits/status/:count/:page", controllers.GetRatelimitsByStatus(db))
			rateLimitRouts.GET("/ratelimits/status/count", controllers.GetRatelimitsByStatusCount(db))
			rateLimitRouts.GET("/ratelimits/datetime/:count/:page", controllers.GetRatelimitsByDatetime(db))
			rateLimitRouts.GET("/ratelimits/datetime/count", controllers.GetRatelimitsByDatetimeCount(db))
			rateLimitRouts.GET("/ratelimits/count", controllers.GetRateLimitCount(db))

			rateLimitRouts.PUT("/ratelimits/status/:id", controllers.UpdateRatelimitStatus(db))
			rateLimitRouts.PUT("/ratelimits/:id", controllers.EditRatelimit(db))
			rateLimitRouts.PUT("/ratelimits/status/bulk/:id", controllers.UpdateRatelimitStatusBulk(db))

			rateLimitRouts.DELETE("/ratelimits/:id", controllers.DeleteRatelimitByID(db))
			rateLimitRouts.DELETE("/ratelimits/bulk/:id", controllers.DeleteRatelimitByIDBulk(db))
		}

		webContentRoutes := api.Group("/webcontent")
		//apply ratelimiter for webcontent subgrooup
		webContentRoutes.Use(rateLimiter.Limit()) //this also possible
		webContentRoutes.Use(middleware.AuthMiddleware(db))
		{
			webContentRoutes.GET("/webcontents", controllers.GetAllWebContents(db)) // get all webcontent
			webContentRoutes.GET("/webcontents/updated", controllers.GetUpdatedWebContents(db))
		}

		BillingRoutes := api.Group("/billing") // web api group
		{
			BillingRoutes.POST("/profiles", controllers.AddBillingProfile(db)) // add transaction
			BillingRoutes.POST("/subscription", controllers.Subscribe(db))     // add transaction

			BillingRoutes.GET("/profiles/:count/:page", controllers.GetBillingProfiles(db))                 // get all transactions
			BillingRoutes.GET("/profile/:id", controllers.GetBillingProfileById(db))                        // get a transactions by id
			BillingRoutes.GET("/profiles/status/:count/:page", controllers.GetBillingProfileByStatus(db))   // get all transactions by status
			BillingRoutes.GET("/profiles/status/count", controllers.GetBillingProfileByStatusCount(db))     // get all transactions by status
			BillingRoutes.GET("/profiles/datetime/:count/:page", controllers.GetBillingProfileDateTime(db)) // get all transactions by datetime
			BillingRoutes.GET("/profiles/datetime/count", controllers.GetBillingProfileByDatetimeCount(db)) // get all transactions by datetime
			BillingRoutes.GET("/profiles/count", controllers.GetBillingProfileCount(db))                    // get all transactions count
			BillingRoutes.GET("/profile/check/:web_id", controllers.CheckBillingProfileExists(db))          // get all transactions total
			BillingRoutes.GET("/subscription/check/:web_id", controllers.CheckSubscriptionExists(db))       // get all transactions total

			BillingRoutes.PUT("/profiles/status/:id", controllers.UpdateBillingProfileStatus(db))          // update transactions status by id
			BillingRoutes.PUT("/profiles/:id", controllers.EditBillingProfile(db))                         // edit transactions by id
			BillingRoutes.PUT("/profiles/status/bulk/:id", controllers.UpdateBillingProfileStatusBulk(db)) // update transactions status by id (bulk)

			BillingRoutes.DELETE("/profiles/:id", controllers.DeleteBillingProfileByID(db))          // delete transactions by ID
			BillingRoutes.DELETE("/profiles/bulk/:id", controllers.DeleteBillingProfileByIDBulk(db)) // delete transactions by ID (bulk)

		}

		SubscriptionRoutes := api.Group("/web") // subscription api group
		{
			//SubscriptionRoutes.POST("/subscriptions", middleware.UserAuthMiddleware(), controllers.AddSubscription(db)) // add subscription
			//SubscriptionRoutes.GET("/subscriptions/:count/:page", controllers.GetSubscriptions(db))                     // get all subscriptions
			SubscriptionRoutes.GET("/subscription/:id", controllers.GetSubscriptionByID(db))

			SubscriptionRoutes.DELETE("/subscription/:id", controllers.DeleteSubscriptionByID(db)) // delete subscription by ID

		}

	}
}
