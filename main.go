package handler

import (
	"database/sql"
	"example/web-service-gin/src/presentation/controller"
	"example/web-service-gin/src/repository/postgres/repository"
	"fmt"
	"net/http"
	"os"

	album_uc "example/web-service-gin/src/usecase/album"

	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func generateDB() (*sql.DB, error) {
	host := os.Getenv("PSQL_HOST")
	dbname := os.Getenv("PSQL_DBNAME")
	user := os.Getenv("PSQL_USER")
	password := os.Getenv("PSQL_PASS")

	return sql.Open(
		"postgres",
		fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", host, dbname, user, password))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}

func main() {
	app = gin.New()
	db, dbErr := generateDB()
	if dbErr != nil {
		panic("failed database connection")
	}
	albumRepo := repository.NewAlbumRepository(db)
	createAlbumUsecase := album_uc.NewCreateAlbumUsecase(albumRepo)
	getAlbumUsecase := album_uc.NewGetAlbumUsecase(albumRepo)
	listAlbumUsecase := album_uc.NewListAlbumUsecase(albumRepo)
	updateAlbumUsecase := album_uc.NewUpdateAlbumUsecase(albumRepo)
	deleteAlbumUsecase := album_uc.NewDeleteAlbumUsecase(albumRepo)
	albumCon := controller.NewAlbumController(
		*createAlbumUsecase,
		*getAlbumUsecase,
		*listAlbumUsecase,
		*updateAlbumUsecase,
		*deleteAlbumUsecase,
	)
	healthCheckCon := controller.NewHealthCheckController()

	app.GET("/albums", albumCon.ListAlbums)
	app.GET("/albums/:id", albumCon.GetAlbumByID)
	app.POST("/albums", albumCon.CreateAlbum)
	app.PUT("/albums", albumCon.UpdateAlbum)
	app.DELETE("/albums/:id", albumCon.DeleteAlbum)
	app.GET("/health", healthCheckCon.HealthCheck)
}

func init() {
	main()
}
