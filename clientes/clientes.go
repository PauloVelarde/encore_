// Service clientes implements a CRUD for clients with a SQL database.
package clientes

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// Client represents a client in the system.
type Client struct {
	ID       int64  `json:"id"`       // El campo ID será parte del JSON de respuesta
	Name     string `json:"name"`     // El nombre del cliente
	Email    string `json:"email"`    // El email del cliente
	Phone    string `json:"phone"`    // El teléfono del cliente
	Address  string `json:"address"`  // La dirección del cliente
}

// CreateClientParams defines the parameters to create a client.
type CreateClientParams struct {
	Name    string `json:"name"`    // Solicitud para crear el cliente (cuerpo JSON)
	Email   string `json:"email"`   // Solicitud para crear el cliente (cuerpo JSON)
	Phone   string `json:"phone"`   // Solicitud para crear el cliente (cuerpo JSON)
	Address string `json:"address"` // Solicitud para crear el cliente (cuerpo JSON)
}

// CreateClientResponse defines the response after creating a client.
type CreateClientResponse struct {
	ID int64 `json:"id"` // El campo ID será parte del JSON de respuesta
}

// UpdateClientParams defines the parameters to update a client.
type UpdateClientParams struct {
	ID       int64  `json:"id"`       // Identificador del cliente a actualizar
	Name     string `json:"name"`     // Actualización del nombre
	Email    string `json:"email"`    // Actualización del email
	Phone    string `json:"phone"`    // Actualización del teléfono
	Address  string `json:"address"`  // Actualización de la dirección
}

// GetClientResponse defines the response after fetching a client.
type GetClientResponse struct {
	Client Client `json:"client"` // El cliente completo será parte del cuerpo de respuesta
}

// ListClientsResponse defines the response for listing clients.
type ListClientsResponse struct {
	Clients []Client `json:"clients"` // Lista de clientes en el cuerpo de respuesta
}

// CreateClient creates a new client and stores it in the database.
//
// encore:api public
func CreateClient(ctx context.Context, params *CreateClientParams) (*CreateClientResponse, error) {
	var id int64
	err := db.QueryRow(ctx, `
		INSERT INTO clients (name, email, phone, address)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, params.Name, params.Email, params.Phone, params.Address).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("could not create client: %v", err)
	}
	return &CreateClientResponse{ID: id}, nil
}

// GetClient retrieves a client by its ID.
//
// encore:api public
func GetClient(ctx context.Context, params *Client) (*GetClientResponse, error) {
	var client Client
	err := db.QueryRow(ctx, `
		SELECT id, name, email, phone, address FROM clients WHERE id = $1
	`, params.ID).Scan(&client.ID, &client.Name, &client.Email, &client.Phone, &client.Address)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve client: %v", err)
	}
	return &GetClientResponse{Client: client}, nil
}

// UpdateClient updates an existing client.
//
// encore:api public
func UpdateClient(ctx context.Context, params *UpdateClientParams) error {
	_, err := db.Exec(ctx, `
		UPDATE clients
		SET name = $1, email = $2, phone = $3, address = $4
		WHERE id = $5
	`, params.Name, params.Email, params.Phone, params.Address, params.ID)
	if err != nil {
		return fmt.Errorf("could not update client: %v", err)
	}
	return nil
}

// DeleteClient deletes a client by its ID.
//
// encore:api public
func DeleteClient(ctx context.Context, params *Client) error {
	_, err := db.Exec(ctx, `DELETE FROM clients WHERE id = $1`, params.ID)
	if err != nil {
		return fmt.Errorf("could not delete client: %v", err)
	}
	return nil
}

// ListClients lists all clients.
//
// encore:api public
func ListClients(ctx context.Context) (*ListClientsResponse, error) {
	rows, err := db.Query(ctx, `SELECT id, name, email, phone, address FROM clients`)
	if err != nil {
		return nil, fmt.Errorf("could not list clients: %v", err)
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.ID, &client.Name, &client.Email, &client.Phone, &client.Address); err != nil {
			return nil, fmt.Errorf("could not scan client: %v", err)
		}
		clients = append(clients, client)
	}
	return &ListClientsResponse{Clients: clients}, nil
}

// Define a database named 'clientes', using the database migrations
// in the "./migrations" folder.
var db = sqldb.NewDatabase("clientes", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
