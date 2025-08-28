package database

import (
	"context"

	"fmt"
	model "myapp/src/internal/order_model"
	generator "myapp/src/testutils"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	db        *gorm.DB
	container testcontainers.Container
	repo      *OrderRepository
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	postgresPort := "5432/tcp"
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{postgresPort},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
		},
		WaitingFor: wait.ForAll(wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort(nat.Port(postgresPort)),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(suite.T(), err)
	suite.container = container

	host, err := container.Host(ctx)
	assert.NoError(suite.T(), err)
	port, err := container.MappedPort(ctx, nat.Port(postgresPort))
	assert.NoError(suite.T(), err)
	dsn := fmt.Sprintf("host=%s user=testuser password=testpassword dbname=testdb port=%d sslmode=disable", host, port.Int())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(suite.T(), err)
	suite.db = db
	err = db.AutoMigrate(&model.Order{}, &model.Item{})
	assert.NoError(suite.T(), err)
}

func (suite *OrderRepositoryTestSuite) TearDownSuite() {
	ctx := context.Background()
	_ = suite.container.Terminate(ctx)
}

func (suite *OrderRepositoryTestSuite) SetupTest() {
	suite.repo = NewOrderRepository(suite.db)
	suite.db.Exec("TRUNCATE TABLE orders, items RESTART IDENTITY CASCADE")
}

func (suite *OrderRepositoryTestSuite) TestCreateOrder() {
	order := generator.GenerateOrder("wb1")
	err := suite.repo.Create(&order)
	assert.NoError(suite.T(), err)
}

func (suite *OrderRepositoryTestSuite) TestCreateTwoSameOrder() {
	order := generator.GenerateOrder("wb1")
	err := suite.repo.Create(&order)
	assert.NoError(suite.T(), err)
	err = suite.repo.Create(&order)
	assert.Error(suite.T(), err)

}

func (suite *OrderRepositoryTestSuite) TestGetOrder(){
	order := generator.GenerateOrder("wb1")
	err := suite.repo.Create(&order)
	assert.NoError(suite.T(), err)
	order2 , err:= suite.repo.GetByTrackNumber("wb1")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), order2.TrackNumber, order.TrackNumber)
}

func (suite * OrderRepositoryTestSuite) TestGetCache(){
	orders := generator.GenerateOrders(100)
	for _ ,order := range orders{
		err := suite.repo.Create(&order)
		assert.NoError(suite.T(), err)
	}
	cache, err := suite.repo.FindForCache(50)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(),50, len(cache))
}

func TestOrderRepositorySuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}
