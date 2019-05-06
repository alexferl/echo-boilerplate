package main

import (
	"fmt"

	"github.com/admiralobvious/echo-logrusmiddleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/admiralobvious/echo-boilerplate/database"
	"github.com/admiralobvious/echo-boilerplate/handlers"
	"github.com/admiralobvious/echo-boilerplate/models"
	"github.com/admiralobvious/echo-boilerplate/util"
)

func init() {
	cnf := NewConfig()
	cnf.BindFlags()
	InitLogging()
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.POST, echo.OPTIONS, echo.DELETE, echo.PUT},
	}))

	if !viper.GetBool("log-requests-disabled") {
		e.Logger = logrusmiddleware.Logger{Logger: logrus.StandardLogger()}
		e.Use(logrusmiddleware.Hook())
	}

	// Database
	db := database.NewDB(viper.GetString("mongodb-uri"))
	db.Dial()
	err := db.Init()
	if err != nil {
		//fail to init db
	}

	createSysAdminUser(db)

	// Routes
	h := &handlers.Handler{DB: db}
	e.GET("/", h.Root)
	e.POST("/login", h.Login)

	// Restricted routes
	r := e.Group("")
	r.Use(middleware.JWT([]byte(viper.GetString("token-secret"))))
	r.GET("/users", h.GetUsers)
	r.POST("/users", h.CreateUser)
	r.GET("/users/:id", h.GetUser)
	r.PUT("/users/:id", h.UpdateUser)
	r.DELETE("/users/:id", h.DeleteUser)

	// Start server
	addr := viper.GetString("address") + ":" + viper.GetString("port")
	e.Logger.Fatal(e.Start(addr))
}

// createSysAdminUser will create a system administrator user if one doesn't exist already
func createSysAdminUser(db database.DB) {
	user := &models.User{}
	db.FindUserByEmail(viper.GetString("sysadmin-email"), user)

	if user.Id == "" {
		genPassword := false
		p := viper.GetString("sysadmin-password")
		if p == "" {
			p = util.RandomString(24)
			genPassword = true
		}

		email := viper.GetString("sysadmin-email")
		u := models.NewUser(email, p)
		u.Type = "sysadmin"
		u.Name = "System Administrator"
		u.Username = viper.GetString("sysadmin-username")
		if len(u.Email) < 1 || len(u.Username) < 1 {
			e := "The system administrator username or email cannot be empty!"
			logrus.Panic(e)
			panic(e)
		}

		db.CreateUser(u)
		m := fmt.Sprintf("System administrator user '%s' successfully created", u.Username)
		if genPassword {
			m += fmt.Sprintf(" with password '%s'", p)
		}
		logrus.Infof(m)
	}
}
