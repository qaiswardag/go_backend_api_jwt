package main

import (
	"fmt"

	"github.com/qaiswardag/go_backend_api_jwt/database"
	"github.com/qaiswardag/go_backend_api_jwt/internal/config"
	"github.com/qaiswardag/go_backend_api_jwt/internal/model"
)

func main() {
	// Load environment variables file
	config.LoadEnvironmentFile()

	db, err := database.InitDB()
	if err != nil {
		panic("failed to connect database")
	}

	// Create tables
	database.DropTables(db)

	// Drop all tables
	database.CreateTables(db)

	// Create 10 fake users
	for i := 1; i <= 10; i++ {
		user := model.User{
			UserName:  fmt.Sprintf("user%d", i),
			FirstName: fmt.Sprintf("FirstName%d", i),
			LastName:  fmt.Sprintf("LastName%d", i),
			Email:     fmt.Sprintf("user%d@example.com", i),
		}
		db.Create(&user)
	}
	// Create 20 fake jobs
	for i := 1; i <= 20; i++ {
		job := model.Job{
			Title: fmt.Sprintf("job%d", i),
			Description: fmt.Sprintln("Consectetuer adipiscing elit. Ac per bibendum quis nec tristique porttitor. Maecenas eros maximus augue, nostra facilisi metus magna. Consequat condimentum mollis luctus molestie turpis et tortor vivamus. Elementum himenaeos potenti tempus nascetur ultrices per. Lacinia tortor eget mus felis magnis luctus. Tellus dis donec erat condimentum per nostra nibh dignissim. Purus sapien finibus mauris vivamus etiam pretium. Hac curae porttitor elementum eget lobortis lobortis. \n\nElementum non sagittis feugiat condimentum dui bibendum ultricies torquent. Sem platea bibendum blandit viverra id urna pellentesque. Phasellus tristique in sodales leo fermentum; cursus dictum. Aptent parturient mus eleifend orci ac. Amet lectus vehicula lacus ac velit. Tortor ex ipsum; fusce hac gravida sagittis porttitor. Ac mollis risus suscipit sodales libero metus magnis.",
				i),
		}
		db.Create(&job)
	}
}
