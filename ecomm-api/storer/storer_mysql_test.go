package storer

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func withTestDB(t *testing.T, fn func(*sqlx.DB, sqlmock.Sqlmock)){
	t.Helper()
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil{
		t.Fatalf("error creating mock database %v", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	
	fn(db, mock)
}

func TestCreateProduct(t *testing.T){
	// create a test product 
	p := &Product{
		Name:      "test product",
		Image:     "test_image.jpg",
		Category: "test category",
		Description: "test description",
		Rating: 5,
		NumReviews: 10,
		Price:     10.99,
		CountInStock:     100.0,
	}
	
	tcs := []struct{
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1,1))
				cp, err := st.CreateProduct(context.Background(), p)
				require.NoError(t, err)
				
				require.Equal(t, int64(1), cp.ID)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "error insert product",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnError(errors.New("error"))
				cp, err := st.CreateProduct(context.Background(), p)
				require.Error(t, err)
				require.Nil(t, cp)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}
	
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T){
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}

}

func TestGetProduct(t *testing.T){
	p := &Product{
		Name:      "test product",
		Image:     "test_image.jpg",
		Category: "test category",
		Description: "test description",
		Rating: 5,
		NumReviews: 10,
		Price:     10.99,
		CountInStock:     100.0,
	}
	tcs := []struct{
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"})
				rows.AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt)
				
				mock.ExpectQuery("SELECT * FROM products WHERE id = ?").WithArgs(1).WillReturnRows(rows)
				
				gp, err := st.GetProduct(context.Background(), 1)
				require.NoError(t, err)
				require.Equal(t, int64(1), gp.ID)
				
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "not found",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM products WHERE id = ?").WithArgs(1).WillReturnError(fmt.Errorf("error geting product"))
				
				gp, err := st.GetProduct(context.Background(), 1)
				require.Error(t, err)
				require.Nil(t, gp)
				
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}
	
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T){
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}
