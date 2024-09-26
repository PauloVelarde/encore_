// Service productos implements a CRUD for products with a SQL database.
package productos

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// Product represents a product in the system.
type Product struct {
	ID    int64   `json:"id"`    // El campo ID será parte del JSON de respuesta
	Namep string  `json:"namep"` // El campo Name será parte del JSON de respuesta
	Price float64 `json:"price"` // El campo Price será parte del JSON de respuesta
	Stock int64   `json:"stock"` // El campo Stock será parte del JSON de respuesta
}

// CreateProductParams defines the parameters to create a product.
type CreateProductParams struct {
	Namep string  `json:"namep"`  // Solicitud para crear el producto (cuerpo JSON)
	Price float64 `json:"price"`  // Solicitud para crear el producto (cuerpo JSON)
	Stock int64   `json:"stock"`  // Solicitud para crear el producto (cuerpo JSON)
}

// CreateProductResponse defines the response after creating a product.
type CreateProductResponse struct {
	ID int64 `json:"id"` // El campo ID será parte del JSON de respuesta
}

// UpdateProductParams defines the parameters to update a product.
type UpdateProductParams struct {
	ID    int64   `json:"id"`    // Identificador del producto a actualizar
	Namep string  `json:"namep"` // Actualización del nombre
	Price float64 `json:"price"` // Actualización del precio
	Stock int64   `json:"stock"` // Actualización del stock
}

// GetProductResponse defines the response after fetching a product.
type GetProductResponse struct {
	Product Product `json:"product"` // El producto completo será parte del cuerpo de respuesta
}

// ListProductsResponse defines the response for listing products.
type ListProductsResponse struct {
	Products []Product `json:"products"` // Lista de productos en el cuerpo de respuesta
}

// CreateProduct creates a new product and stores it in the database.
//
//encore:api public
func CreateProduct(ctx context.Context, params *CreateProductParams) (*CreateProductResponse, error) {
	var id int64
	err := db.QueryRow(ctx, `
		INSERT INTO products (namep, price, stock)
		VALUES ($1, $2, $3)
		RETURNING id
	`, params.Namep, params.Price, params.Stock).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("could not create product: %v", err)
	}
	return &CreateProductResponse{ID: id}, nil
}

// GetProduct retrieves a product by its ID.
//
//encore:api public
func GetProduct(ctx context.Context, params *Product) (*GetProductResponse, error) {
	var product Product
	err := db.QueryRow(ctx, `
		SELECT id, namep, price, stock FROM products WHERE id = $1
	`, params.ID).Scan(&product.ID, &product.Namep, &product.Price, &product.Stock)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve product: %v", err)
	}
	return &GetProductResponse{Product: product}, nil
}

// UpdateProduct updates an existing product.
//
//encore:api public
func UpdateProduct(ctx context.Context, params *UpdateProductParams) error {
	_, err := db.Exec(ctx, `
		UPDATE products
		SET namep = $1, price = $2, stock = $3
		WHERE id = $4
	`, params.Namep, params.Price, params.Stock, params.ID)
	if err != nil {
		return fmt.Errorf("could not update product: %v", err)
	}
	return nil
}

// DeleteProduct deletes a product by its ID.
//
//encore:api public
func DeleteProduct(ctx context.Context, params *Product) error {
	_, err := db.Exec(ctx, `DELETE FROM products WHERE id = $1`, params.ID)
	if err != nil {
		return fmt.Errorf("could not delete product: %v", err)
	}
	return nil
}

// ListProducts lists all products.
//
//encore:api public
func ListProducts(ctx context.Context) (*ListProductsResponse, error) {
	rows, err := db.Query(ctx, `SELECT id, namep, price, stock FROM products`)
	if err != nil {
		return nil, fmt.Errorf("could not list products: %v", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Namep, &product.Price, &product.Stock); err != nil {
			return nil, fmt.Errorf("could not scan product: %v", err)
		}
		products = append(products, product)
	}
	return &ListProductsResponse{Products: products}, nil
}

// Define a database named 'productos', using the database migrations
// in the "./migrations" folder.
var db = sqldb.NewDatabase("productos", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
