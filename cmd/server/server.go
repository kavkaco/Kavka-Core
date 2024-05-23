package main

func main() {
	// // Init Zap Logger
	// logger := logs.InitZapLogger()

	// // Define paths
	// TemplatesPath := config.ProjectRootPath + "/app/views/mail/"

	// // Load Configs
	// configs := config.Read()

	// // Init MongoDB
	// mongoDB, mongoErr := database.GetMongoDBInstance(
	// 	database.NewMongoDBConnectionString(
	// 		configs.Mongo.Host,
	// 		configs.Mongo.Port,
	// 		configs.Mongo.Username,
	// 		configs.Mongo.Password,
	// 	),
	// 	configs.Mongo.DBName,
	// )
	// if mongoErr != nil {
	// 	panic(mongoErr)
	// }

	// // Init RedisDB
	// redisClient := database.GetRedisDBInstance(configs.Redis)

	// // Init WebServer
	// ginEngine := gin.New()
	// api := ginEngine.Group("/api/v1")

	// // Cors
	// api.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{configs.App.Server.CORS.AllowOrigins},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Refresh", "Authorization"},
	// 	ExposeHeaders:    []string{"Refresh", "Authorization"},
	// 	AllowCredentials: true,
	// }))

	// // Initializing various services and repositories used in the application
	// session := session.NewSession(logger, redisClient, configs.App.Auth)
	// smsService := sms_service.NewSmsService(logger, &configs.SMS, TemplatesPath)

	// userRepo := userRepository.NewRepository(logger, mongoDB)
	// userService := service.NewUserService(logger, userRepo, session, smsService)

	// chatRepo := chatRepository.NewRepository(logger, mongoDB)
	// chatService := service.NewChatService(logger, chatRepo, userRepo)

	// messageRepo := messageRepository.NewRepository(logger, mongoDB)
	// messageService := service.NewMessageService(logger, messageRepo, chatRepo)

	// // Init routes
	// router.NewUserRouter(logger, api.Group("/users"), userService, chatService)

	// // Init websocket server
	// websocketAdapter := adapters.NewWebsocketAdapter(logger)

	// handlerServices := handlers.HandlerServices{
	// 	UserService: userService,
	// 	ChatService: chatService,
	// 	MsgService:  messageService,
	// }

	// api.GET("/ws", middleware.AuthenticatedMiddleware(userService), router.WebsocketRoute(logger, websocketAdapter, handlerServices))

	// // Everything almost done!
	// err := ginEngine.Run(configs.App.HTTP.Address)
	// if err != nil {
	// 	panic(err)
	// }
}
