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
			webRoutes.POST("/site", middleware.UserAuthMiddleware(), controllers.AddSite(db))   // add site
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

		adminDashboardRoutes := api.Group("/admin_dashboard") // admin dashboard api group
		{
			adminDashboardRoutes.GET("/usersTotalCount", controllers.GetTotalUserCount(db)) // get total user count
			adminDashboardRoutes.GET("/websitesTotalCount", controllers.GetTotalWebsitesCount(db))
			adminDashboardRoutes.GET("/apiSubscribersTotalCount", controllers.GetTotalApiSubscribersCount(db))
			adminDashboardRoutes.GET("/marketplaceUsersTotalCount", controllers.GetTotalMarketplaceUsersCount(db))
			adminDashboardRoutes.GET("/sites/storage", controllers.GetSitesStorage(db)) // get all sites storage
			adminDashboardRoutes.GET("/sites/totalStorage", controllers.GetTotalUsedStorage(db))
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
			autoRespondRoutes.GET("/auto_respond/:count/:page/:webId", controllers.GetAutoResponds(db)) //Get a list of auto-responses with pagination and webId filtering
			autoRespondRoutes.POST("/auto_respond/:webId", controllers.AddAutoRespond(db))
			autoRespondRoutes.GET("/auto_respond/id/:id/:webId", controllers.GetAutoRespondsById(db))                      // get a AutoResponds by id
			autoRespondRoutes.GET("/auto_respond/status/:count/:page/:webId", controllers.GetAutoRespondsByStatus(db))     // get all AutoResponds by status
			autoRespondRoutes.GET("/auto_respond/status/count/:webId", controllers.GetAutoRespondsByStatusCount(db))       // get all AutoResponds by status
			autoRespondRoutes.GET("/auto_respond/datetime/:count/:page/:webId", controllers.GetAutoRespondsByDatetime(db)) // get all AutoResponds by datetime
			autoRespondRoutes.GET("/auto_respond/datetime/count/:webId", controllers.GetAutoRespondsByDatetimeCount(db))   // get all AutoResponds by datetime
			autoRespondRoutes.GET("/auto_respond/count/:webId", controllers.GetAutoRespondsCount(db))                      // get all AutoResponds count

			autoRespondRoutes.PUT("/auto_respond/status/:id/:webId", controllers.UpdateAutoRespondsStatus(db)) // update AutoResponds status by id
			autoRespondRoutes.PUT("/auto_respond/:id/:webId", controllers.EditAutoResponds(db))                // edit AutoResponds by id
			autoRespondRoutes.PUT("/auto_respond/status/bulk/:id/:webId", controllers.UpdateAutoRespondsStatusBulk(db))

			autoRespondRoutes.DELETE("/auto_respond/:id/:webId", controllers.DeleteAutoRespondsID(db)) // delete AutoResponds by ID
			autoRespondRoutes.DELETE("/auto_respond/bulk/:id/:webId", controllers.DeleteAutoRespondsByIDBulk(db))
			autoRespondRoutes.GET("/auto_respond/get/:webId", controllers.GetAutoRespondsByWebID(db)) // delete AutoResponds by ID (bulk)
		}

		analyticalAlertsRoutes := api.Group("/analytics") // analytical alerts api group
		{
			analyticalAlertsRoutes.GET("/visitorsInfo/:count/:page/:id", controllers.GetVisitorInfo(db)) // get all webpages
			analyticalAlertsRoutes.GET("/visitorInfo/:id", controllers.GetVisitorInfoById(db))           // get a webpage by id
			//analyticalAlertsRoutes.GET("/", controllers.GetAnalyticalAlerts(db)) // get all analytical alerts
			analyticalAlertsRoutes.GET("/visitorInfo/datetime/:count/:page", controllers.GetVisitorInfoByDatetime(db)) // get all webpages by datetime
			analyticalAlertsRoutes.GET("/visitorInfo/datetime/count", controllers.GetVisitorByDatetimeCount(db))       // get all webpages by datetime
			analyticalAlertsRoutes.GET("/visitorInfo/count/:id", controllers.GetVisitorInfoCount(db))                  // get all webpages count

			analyticalAlertsRoutes.GET("/source/:id", controllers.GetSource(db))
			analyticalAlertsRoutes.GET("/sessions/:id", controllers.GetSessions(db))
			analyticalAlertsRoutes.GET("/devices/:id", controllers.GetDevices(db))
			analyticalAlertsRoutes.GET("/country/:id", controllers.GetCountry(db))
			analyticalAlertsRoutes.POST("/Alert", controllers.CreateNewAlert(db)) //create new alert

			analyticalAlertsRoutes.GET("/Alerts/:count/:page/:id", controllers.GetAllAlert(db))         // get all alerts
			analyticalAlertsRoutes.GET("/Alert/:id", controllers.GetAlertbyId(db))                      // get alert by id
			analyticalAlertsRoutes.GET("/Alert/status/:count/:page", controllers.GetAlertsByStatus(db)) // get all webpages by status
			analyticalAlertsRoutes.GET("/Alert/status/count", controllers.GetAlertsByStatusCount(db))   // get all webpages by status

			analyticalAlertsRoutes.GET("/Alert/count/:id", controllers.GetAlertsCount(db))              // get all webpages count
			analyticalAlertsRoutes.PUT("/Alert/status/:id", controllers.UpdateAlertStatus(db))          // update webpage status by id
			analyticalAlertsRoutes.PUT("/Alert/:id", controllers.EditAlert(db))                         // edit webpage by id
			analyticalAlertsRoutes.PUT("/Alert/status/bulk/:id", controllers.UpdateAlertStatusBulk(db)) // update webpage status by id (bulk)
			//
			analyticalAlertsRoutes.DELETE("/Alert/:id", controllers.DeleteAlertByID(db))          // delete webpage by ID
			analyticalAlertsRoutes.DELETE("/Alert/bulk/:id", controllers.DeleteAlertByIDBulk(db)) // delete webpage by ID (bulk)
			analyticalAlertsRoutes.POST("/Alert/sessionrecord", controllers.SessionRecord(db))    //add session record
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

		templateRoutes := api.Group("/marketplace") // marketplace api group
		{
			templateRoutes.POST("/template", middleware.UserAuthMiddleware(), controllers.AddTemplate(db))                 // add new template
			templateRoutes.POST("/template/rating", middleware.UserAuthMiddleware(), controllers.AddRatings(db))           // add new template rating
			templateRoutes.POST("/template/upload", middleware.UserAuthMiddleware(), controllers.UploadTemplate(db))       // upload template
			templateRoutes.POST("/template/image/upload", middleware.UserAuthMiddleware(), controllers.UploadThumbImg(db)) // upload template thumbnail image

			templateRoutes.GET("/templates/:count/:page", controllers.GetTemplates(db))                                            // get all
			templateRoutes.GET("/template/:id", controllers.GetTemplatesById(db))                                                  // get by id
			templateRoutes.GET("/templates/status/:count/:page", controllers.GetTemplatesByStatus(db))                             // get all by status
			templateRoutes.GET("/templates/status/count", controllers.GetTemplatesByStatusCount(db))                               // get all by status
			templateRoutes.GET("/templates/datetime/:count/:page", controllers.GetTemplatesByDatetime(db))                         // get all by datetime
			templateRoutes.GET("/templates/datetime/count", controllers.GetTemplatesByDatetimeCount(db))                           // get all by datetime
			templateRoutes.GET("/templates/count", controllers.GetTemplatesCount(db))                                              // get all count
			templateRoutes.GET("/templat/:id", controllers.DownloadById(db))                                                       // download templates by id
			templateRoutes.GET("/templates/user/:count/:page", middleware.UserAuthMiddleware(), controllers.GetTemplatesBydid(db)) // get all templates by user id
			templateRoutes.GET("/templates/acceptstatus/:count/:page", controllers.GetAcceptedTemplates(db))                       // get all the accepted templates
			templateRoutes.GET("/templates/filter/:count/:page/:category", controllers.GetTemplatesByCategory(db))                 // filter function

			templateRoutes.PUT("/templates/status/:id", controllers.UpdateTemplatesStatus(db))                    // update status by id
			templateRoutes.PUT("/templates/:id", middleware.UserAuthMiddleware(), controllers.EditTemplatesD(db)) // edit template details by id
			templateRoutes.PUT("/templates/status/bulk/:id", controllers.UpdateTemplatesStatusBulk(db))           // update status by id (bulk)

			templateRoutes.DELETE("/templates/:id", controllers.DeleteTemplateByID(db))          // delete by ID
			templateRoutes.DELETE("/templates/bulk/:id", controllers.DeleteTemplateByIDBulk(db)) // delete by ID (bulk)

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
		webContentRoutes.GET("/allSites", controllers.GetAllDpacksSites(db)) // get all webcontent
		webContentRoutes.GET("/webcontents/updated/:limit", controllers.GetUpdatedWebContents(db))
	}

		BillingRoutes := api.Group("/billing") // web api group
		{
			BillingRoutes.POST("/profiles", controllers.AddBillingProfile(db)) // add transaction

			BillingRoutes.GET("/profiles/:count/:page", controllers.GetBillingProfiles(db))                 // get all transactions
			BillingRoutes.GET("/profile/:id", controllers.GetBillingProfileById(db))                        // get a transactions by id
			BillingRoutes.GET("/profiles/status/:count/:page", controllers.GetBillingProfileByStatus(db))   // get all transactions by status
			BillingRoutes.GET("/profiles/status/count", controllers.GetBillingProfileByStatusCount(db))     // get all transactions by status
			BillingRoutes.GET("/profiles/datetime/:count/:page", controllers.GetBillingProfileDateTime(db)) // get all transactions by datetime
			BillingRoutes.GET("/profiles/datetime/count", controllers.GetBillingProfileByDatetimeCount(db)) // get all transactions by datetime
			BillingRoutes.GET("/profiles/count", controllers.GetBillingProfileCount(db))                    // get all transactions count
			BillingRoutes.GET("/profile/check/:web_id", controllers.CheckBillingProfileExists(db))          // get all transactions total

			BillingRoutes.PUT("/profiles/status/:id", controllers.UpdateBillingProfileStatus(db))          // update transactions status by id
			BillingRoutes.PUT("/profiles/:id", controllers.EditBillingProfile(db))                         // edit transactions by id
			BillingRoutes.PUT("/profiles/status/bulk/:id", controllers.UpdateBillingProfileStatusBulk(db)) // update transactions status by id (bulk)

			BillingRoutes.DELETE("/profiles/:id", controllers.DeleteBillingProfileByID(db))          // delete transactions by ID
			BillingRoutes.DELETE("/profiles/bulk/:id", controllers.DeleteBillingProfileByIDBulk(db)) // delete transactions by ID (bulk)
		}

		SubscriptionRoutes := api.Group("/subscription") // subscription api group
		{
			SubscriptionRoutes.POST("/", controllers.Subscribe(db)) // add transaction

			SubscriptionRoutes.PUT("/", controllers.UpdateSubscribe(db)) // add transaction

			SubscriptionRoutes.GET("/check/:web_id", controllers.CheckSubscriptionExists(db)) // get all transactions total
			SubscriptionRoutes.GET("/:id", controllers.GetSubscriptionByID(db))

			SubscriptionRoutes.DELETE("/:id", controllers.DeleteSubscriptionByID(db)) // delete subscription by ID
		}

	}

	StorageRoutes := api.Group("/storage")
	{
		StorageRoutes.GET("/:id", controllers.GetStorageByID(db))
	}
}
