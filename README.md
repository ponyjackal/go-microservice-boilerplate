# Go Microservices Boilerplate

This Go boilerplate project is designed to streamline the development of microservices using gRPC for inter-service communication, Protocol Buffers (protobuf) for data serialization, Gorm as an ORM for database operations, and Gormigrate for database migrations.

### Motivation

The purpose of this boilerplate is to provide a robust starting point for building scalable and efficient microservices in Go. It aims to incorporate best practices and tools for rapid development and high performance.

## Table of Contents

- [Motivation](#motivation)
- [Configuration Manage](#configuration-manage)
  - [ENV Manage](#env-manage)
  - [Server Configuration](#server-configuration)
  - [Database Configuration](#database-configuration)
- [Installation](#installation)
  - [Local Setup Instruction](#local-setup-instruction)
  - [Develop Application in Docker with Live Reload](#develop-application-in-docker-with-live-reload)
- [Middlewares](#middlewares)
- [Boilerplate Structure](#boilerplate-structure)
- [Let's Build an API](#lets-build-an-api)
- [Deployment](#deployment)
  - [Container Development Build](#container-development-build)
  - [Container Production Build and Up](#container-production-build-and-up)
- [Useful Commands](#useful-commands)
- [ENV YAML Configure](#env-yaml-configure)
- [Use Packages](#use-packages)

### Configuration Manage

#### ENV Manage

- Default ENV Configuration Manage from `.env`. sample file `.env.example`

```text
# Server Configuration
SECRET=h9wt*pasj6796j##w(w8=xaje8tpi6h*r&hzgrz065u&ed+k2)
DEBUG=True # `False` in Production
ALLOWED_HOSTS=0.0.0.0
SERVER_HOST=0.0.0.0
SERVER_PORT=8000

# Database Configuration
MASTER_DB_NAME=test_pg_go
MASTER_DB_USER=mamun
MASTER_DB_PASSWORD=123
MASTER_DB_HOST=postgres_db
MASTER_DB_PORT=5432
MASTER_DB_LOG_MODE=True # `False` in Production
MASTER_SSL_MODE=disable

REPLICA_DB_NAME=test_pg_go
REPLICA_DB_USER=mamun
REPLICA_DB_PASSWORD=123
REPLICA_DB_HOST=localhost
REPLICA_DB_PORT=5432
REPLICA_DB_LOG_MODE=True # `False` in Production
REPLICA_SSL_MODE=disable
```

- Server `DEBUG` set `False` in Production
- Database Logger `MASTER_DB_LOG_MODE` and `REPLICA_DB_LOG_MODE` set `False` in production
- If ENV Manage from YAML file add a config.yml file and configuration [db.go](config/db.go) and [server.go](config/server.go). See More [ENV YAML Configure](#env-yaml-configure)

#### Server Configuration

- Use [Gin](https://github.com/gin-gonic/gin) Web Framework

#### Database Configuration

- Use [GORM](https://github.com/go-gorm/gorm) as an ORM
- Use database `MASTER_DB_HOST` value set as `localhost` for local development, and use `postgres_db` for docker development

### Installation

#### Local Setup Instruction

Follow these steps:

- Copy [.env.example](.env.example) as `.env` and configure necessary values
- To add all dependencies for a package in your module `go get .` in the current directory
- Locally run `go run main.go` or `go build main.go` and run `./main`
- Check Application health available on [0.0.0.0:8000/health](http://0.0.0.0:8000/health)

#### Develop Application in Docker with Live Reload

Follow these steps:

- Make sure install the latest version of docker and docker-compose
- Docker Installation for your desire OS https://docs.docker.com/engine/install/ubuntu/
- Docker Composer Installation https://docs.docker.com/compose/install/
- Run and Develop `make dev`
- Check Application health available on [0.0.0.0:8000/health](http://0.0.0.0:8000/health)

### Middlewares

- Use Gin CORSMiddleware

```go
router := gin.New()
router.Use(gin.Logger())
router.Use(gin.Recovery())
router.Use(middleware.CORSMiddleware())
```

### Directory Structure

<pre>├── <font color="#3465A4"><b>internal</b></font>
│   ├── <font color="#3465A4"><b>adapters</b></font>
│   │   ├── <font color="#3465A4"><b>database</b></font>
│   │   │   ├── <font color="#3465A4"><b>migrations</b></font>
│   │   │   │   └── migration.go
│   │   │   ├── <font color="#3465A4"><b>seeds</b></font>
│   │   │   └── database.go
│   ├── <font color="#3465A4"><b>app</b></font>
│   │   ├── <font color="#3465A4"><b>controllers</b></font>
│   │   │   └── tag_controller.go
│   │   ├── <font color="#3465A4"><b>middlewares</b></font>
│   │   │   └── cors.go
│   │   ├── <font color="#3465A4"><b>routers</b></font>
│   │   │   ├── index.go
│   │   │   └── router.go
│   └── <font color="#3465A4"><b>domain</b></font>
│   │   ├── <font color="#3465A4"><b>models</b></font>
│   │   │   └── tag_model.go
│   │   ├── <font color="#3465A4"><b>repositories</b></font>
│   │   │   └── tag_repository.go
│   │   ├── <font color="#3465A4"><b>services</b></font>
│   │   │   └── tag_service.go
├── <font color="#3465A4"><b>pkg</b></font>
│   ├── <font color="#3465A4"><b>config</b></font>
│   │   ├── config.go
│   │   ├── db.go
│   │   └── server.go
│   ├── <font color="#3465A4"><b>constants</b></font>
│   │   └── constants.go
│   ├── <font color="#3465A4"><b>logger</b></font>
│   │   └── logger.go
│   ├── <font color="#3465A4"><b>types</b></font>
│   ├── <font color="#3465A4"><b>utils</b></font>
│   │   └── utils.go
│   └── <font color="#3465A4"><b>workers</b></font>
├── <font color="#3465A4"><b>proto</b></font>
│   ├── <font color="#3465A4"><b>service</b></font>
│   │   ├── service_grpc.pb
│   │   ├── service.pb.go
│   │   ├── service.pb.gw.go
│   │   └── service.proto
│   ├── <font color="#3465A4"><b>shared</b></font>
│   │   ├── shared.pb.go
│   │   └── shared.proto
│   ├── <font color="#3465A4"><b>tag</b></font>
│   │   ├── tag.pb.go
│   │   └── tag.proto
│   ├── buf.lock
│   └── buf.yaml
├── <font color="#3465A4"><b>docs</b></font>
│   ├── docs.go
│   ├── grpc_gateway.go
│   ├── grpc_gateway.swagger.json
│   ├── swagger.json
│   └── swagger.yaml
├── buf.gen.yaml
├── docker-compose-dev.yml
├── docker-compose-prod.yml
├── Dockerfile
├── Dockerfile-dev
├── go.mod
├── go.sum
├── LICENSE
├── main.go
├── Makefile
</pre>

### Let's Build an API

1. [models](models) folder add a new file name `tag_model.go`

```go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID uuid.UUID `gorm:"type:uuid;column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	/* Fields */
	Name string `gorm:"not null;uniqueIndex:unique_tag_name" json:"name"`
	/* Timestamp */
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName is Database TableName of this model
func (e *Tag) TableName() string {
	return "tags"
}

```

2. Add Model to [migration](internal/adapters/database/migratins/migration.go)

```go
package migrations

package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/ponyjackal/go-microservice-boilerplate/internal/adapters/database"
	"github.com/ponyjackal/go-microservice-boilerplate/internal/domain/models"
)

var migrations = []*gormigrate.Migration{}

func Migrate() {
	m := gormigrate.New(database.DB, gormigrate.DefaultOptions, migrations)

	m.InitSchema(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(
			&models.Tag{},
		)
		if err != nil {
			return err
		}

		return nil
	})

	// Run the migrations
	if err := m.Migrate(); err != nil {
		panic(err)
	}
}

```

3. [controller](controllers) folder add a file `tag_controller.go`

- Create API Endpoint
- Write Database Operation in Repository and use them from controller

```go
package controllers

import (
	"net/http"

	"github.com/ponyjackal/go-microservice-boilerplate/internal/domain/services"
	"github.com/ponyjackal/go-microservice-boilerplate/pkg/logger"
	"github.com/ponyjackal/go-microservice-boilerplate/pkg/utils"
	pbTag "github.com/ponyjackal/go-microservice-boilerplate/proto/tag"

	"github.com/bufbuild/protovalidate-go"
	"github.com/gin-gonic/gin"
)

type TagController struct {
	tagService *services.TagService
	validator  *protovalidate.Validator
}

func NewTagController(tagService *services.TagService) *TagController {
	// Create a new validator
	validator, err := protovalidate.New()
	if err != nil {
		logger.Errorf("failed to initialize validator: %v", err)
	}

	return &TagController{
		tagService: tagService,
		validator:  validator,
	}
}

// GetTags godoc
// @Summary Retrieve a list of tags
// @Description Get a list of tags filtered by the name parameter
// @Tags Tags
// @Accept json
// @Produce json
// @Param name query string false "Name of the tag to filter by"
// @Success 200 {object} pbTag.GetTagsResponse "Successful retrieval of tags"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /tags [get]
func (c *TagController) GetTags(ctx *gin.Context) {
	name := ctx.Query("name")

	response, err := c.tagService.GetTags(name)
	if err != nil {
		utils.GRPCErrorHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// GetTagById godoc
// @Summary Retrieve a tag
// @Description Get tag by id from the database
// @Tags Tags
// @Produce json
// @Param id path string true "Tag ID"
// @Success 200 {object} pbTag.Tag "Successfully retrieved a tag"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /tags/{id} [get]
func (c *TagController) GetTagById(ctx *gin.Context) {
	id := ctx.Param("id")

	response, err := c.tagService.GetTagById(&pbTag.TagId{
		Id: id,
	})
	if err != nil {
		utils.GRPCErrorHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// SaveTag godoc
// @Summary Add a new tag
// @Description Add a tag with the provided information
// @Tags Tags
// @Accept json
// @Produce json
// @Param tag body pbTag.SaveTagRequest true "Tag Object"
// @Success 201 {object} models.Tag "Successfully created tag"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /tags [post]
func (c *TagController) SaveTag(ctx *gin.Context) {
	var tagReq pbTag.SaveTagRequest

	if err := ctx.BindJSON(&tagReq); err != nil {
		logger.Errorf("Invalid request: %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err})
		return
	}
	// Validate the request
	if err := c.validator.Validate(&tagReq); err != nil {
		logger.Errorf("Invalid request: %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err})
		return
	}

	response, err := c.tagService.SaveTag(&tagReq)
	if err != nil {
		utils.GRPCErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// UpdateTag godoc
// @Summary Update tag
// @Description Update a tag with the provided information
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path string true "Tag ID"
// @Param tag body pbTag.SaveTagRequest true "Tag Object"
// @Success 200 {object} pbTag.Tag "Successfully updated a tag"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /tags/{id} [put]
func (c *TagController) UpdateTag(ctx *gin.Context) {
	id := ctx.Param("id")

	var tagReq pbTag.SaveTagRequest
	if err := ctx.BindJSON(&tagReq); err != nil {
		logger.Errorf("Invalid request: %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err})
		return
	}

	tag, err := c.tagService.UpdateTag(&pbTag.UpdateTagRequest{
		Id:     id,
		TagReq: &tagReq,
	})
	if err != nil {
		utils.GRPCErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &tag)
}

// DeleteTag godoc
// @Summary Delete tag
// @Description Delete a tag with the provided information
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path string true "Tag ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /tags/{id} [delete]
func (c *TagController) DeleteTag(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.tagService.DeleteTag(&pbTag.TagId{
		Id: id,
	})
	if err != nil {
		utils.GRPCErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}

```

4. [routers](routers) folder add a new route in `index.go`

```go
package routers

import (
	"net/http"

	"github.com/ponyjackal/go-microservice-boilerplate/internal/app/controllers"
	"github.com/ponyjackal/go-microservice-boilerplate/internal/domain/services"

	_ "github.com/ponyjackal/go-microservice-boilerplate/docs"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	route *gin.Engine,
	tagService *services.TagService,
) {
	/* Controllers */
	tagController := controllers.NewTagController(tagService)

	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})

	v1 := route.Group("api/v1")
	{
		// health check
		v1.GET("health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "good"}) })
		// tags
		tags := v1.Group("tags")
		{
			tags.GET("", tagController.GetTags)
			tags.GET(":id", tagController.GetTagById)
			tags.POST("", tagController.SaveTag)
			tags.PUT(":id", tagController.UpdateTag)
			tags.DELETE(":id", tagController.DeleteTag)
		}
	}
}

```

- Congratulation, your new endpoint `0.0.0.0:8000/api/v1/tags`

### Deployment

#### Container Development Build

- Run `make build`

#### Container Production Build and Up

- Run `make production`

#### ENV Yaml Configure

```yaml
database:
  driver: "postgres"
  dbname: "test_pg_go"
  username: "mamun"
  password: "123"
  host: "postgres_db" # use `localhost` for local development
  port: "5432"
  ssl_mode: disable
  log_mode: false

server:
  host: "0.0.0.0"
  port: "8000"
  secret: "secret"
  allow_hosts: "localhost"
  debug: false #use `false` in production
  request:
    timeout: 100
```

- [Server Config](pkg/config/server.go)

```go
func ServerConfig() string {
	appServer := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))
	logger.Infof("Server Running at : %s", appServer)
	return appServer
}

```

- [DB Config](pkg/config/db.go)

```go
func DbConfiguration() (string, string) {
	masterDBName := os.Getenv("MASTER_DB_NAME")
	masterDBUser := os.Getenv("MASTER_DB_USER")
	masterDBPassword := os.Getenv("MASTER_DB_PASSWORD")
	masterDBHost := os.Getenv("MASTER_DB_HOST")
	masterDBPort := os.Getenv("MASTER_DB_PORT")
	masterDBSslMode := os.Getenv("MASTER_SSL_MODE")

	replicaDBName := os.Getenv("REPLICA_DB_NAME")
	replicaDBUser := os.Getenv("REPLICA_DB_USER")
	replicaDBPassword := os.Getenv("REPLICA_DB_PASSWORD")
	replicaDBHost := os.Getenv("REPLICA_DB_HOST")
	replicaDBPort := os.Getenv("REPLICA_DB_PORT")
	replicaDBSslMode := os.Getenv("REPLICA_SSL_MODE")

	masterDBDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		masterDBHost, masterDBUser, masterDBPassword, masterDBName, masterDBPort, masterDBSslMode,
	)

	replicaDBDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		replicaDBHost, replicaDBUser, replicaDBPassword, replicaDBName, replicaDBPort, replicaDBSslMode,
	)
	return masterDBDSN, replicaDBDSN
}
```

### Useful Commands

- `make dev`: make dev for development work
- `make build`: make build container
- `make production`: docker production build and up
- `make protobuf`: generate protobuf go files
- `make doc`: generate swagger doc
- `clean`: clean for all clear docker images

### Use Packages

- [Gorm](https://github.com/go-gorm/gorm) - The fantastic ORM library for Golang
- [Logger](https://github.com/sirupsen/logrus) - Structured, pluggable logging for Go..
- [Air](https://github.com/cosmtrek/air) - Live reload for Go apps (Docker Development)
