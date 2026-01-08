package repository

import (
	"context"
	"database/sql"
	"products-api/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProductRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo *ProductRepository
}

var (
	product models.Product = models.Product{
		BaseModel: models.BaseModel{ID: "1"},
		Name:      "Test Product",
		Price:     9.99,
		Quantity:  100,
	}
)

func (suite *ProductRepositoryTestSuite) SetupSuite() {
	db, mock, err := sqlmock.New()
	suite.NoError(err)
	suite.db = db
	suite.mock = mock
	suite.repo = NewProductRepository(db)
}

func (suite *ProductRepositoryTestSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *ProductRepositoryTestSuite) SetupTest() {
	// Reset expectations if needed, but for simplicity, set per test
}

func (suite *ProductRepositoryTestSuite) TearDownTest() {
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func setupProductMock(mock sqlmock.Sqlmock) {
	fixedTime := time.Now()
	product := models.Product{
		BaseModel: models.BaseModel{ID: "1"},
		Name:      "Test Product",
		Price:     9.99,
		Quantity:  100,
	}

	mock.ExpectQuery("INSERT INTO products").WithArgs(product.ID, product.Name, product.Price, product.SellerID, product.Quantity).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "seller_id", "quantity", "created_at", "updated_at"}).AddRow(product.ID, "Test Product", 9.99, "", 0, fixedTime, fixedTime))
	mock.ExpectQuery("SELECT .* FROM products WHERE .*").WithArgs(product.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "seller_id", "quantity", "created_at", "updated_at"}).AddRow(product.ID, "Test Product", 9.99, "", 0, fixedTime, fixedTime))
}

func (suite *ProductRepositoryTestSuite) TestCreateProduct() {
	setupProductMock(suite.mock)

	err := suite.repo.Create(context.Background(), &product)
	suite.NoError(err, "expected no error while creating product")

	storedProduct, err := suite.repo.GetProductByID(context.Background(), product.ID)
	suite.NoError(err, "expected no error while retrieving product")
	assert.Equal(suite.T(), product, storedProduct, "expected retrieved product to match created product")
}

func (suite *ProductRepositoryTestSuite) TestUpdateProductCount() {
	fixedTime := time.Now()
	suite.mock.ExpectExec("UPDATE .*").WithArgs(95, "1").WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectQuery("SELECT .* FROM products WHERE .*").WithArgs("1").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "seller_id", "quantity", "created_at", "updated_at"}).AddRow("1", "Test Product", 9.99, "", 95, fixedTime, fixedTime))
	err := suite.repo.UpdateProductCount(context.Background(), &product, 5)
	suite.NoError(err, "expected no error while updating product count")
	assert.Equal(suite.T(), 95, product.Quantity, "expected product quantity to be updated correctly")

	// Additionally check via GetProductByID
	resultProduct, err := suite.repo.GetProductByID(context.Background(), "1")
	suite.NoError(err, "expected no error while retrieving product")
	assert.Equal(suite.T(), 95, resultProduct.Quantity, "expected retrieved product quantity to be updated")
}

func (suite *ProductRepositoryTestSuite) TestGetAllProducts() {
	expectedProducts := []models.Product{
		{BaseModel: models.BaseModel{ID: "1"}, Name: "Product 1", Price: 10.0, SellerID: "seller1", Quantity: 5},
		{BaseModel: models.BaseModel{ID: "2"}, Name: "Product 2", Price: 20.0, SellerID: "seller2", Quantity: 3},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "price", "seller_id", "quantity", "created_at", "updated_at"})
	for _, p := range expectedProducts {
		rows.AddRow(p.ID, p.Name, p.Price, p.SellerID, p.Quantity, time.Now(), time.Now())
	}

	suite.mock.ExpectQuery("SELECT .* FROM products$").WillReturnRows(rows)

	result, err := suite.repo.GetAll(context.Background())
	suite.NoError(err, "expected no error while getting all products")
	suite.Len(result, 2, "expected two products")
	// Note: Exact match may fail due to time fields, so check key fields
	for i, p := range result {
		assert.Equal(suite.T(), expectedProducts[i].ID, p.ID)
		assert.Equal(suite.T(), expectedProducts[i].Name, p.Name)
		assert.Equal(suite.T(), expectedProducts[i].Price, p.Price)
		assert.Equal(suite.T(), expectedProducts[i].SellerID, p.SellerID)
		assert.Equal(suite.T(), expectedProducts[i].Quantity, p.Quantity)
	}

}

func (suite *ProductRepositoryTestSuite) TestDeleteProduct() {
	suite.mock.ExpectExec("DELETE FROM products WHERE .*").WithArgs("1").WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs(product.ID)

	err := suite.repo.DeleteProduct(context.Background(), "1")
	suite.NoError(err, "expected no error while deleting product")

	suite.mock.ExpectQuery("SELECT .* FROM products WHERE .*").WithArgs("1").WillReturnError(sql.ErrNoRows)

	deletedProduct, err := suite.repo.GetProductByID(context.Background(), "1")
	suite.Error(err, "expected error while retrieving deleted product")
	suite.Nil(deletedProduct, "expected no product to be returned after deletion")
}

func TestProductRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ProductRepositoryTestSuite))
}
