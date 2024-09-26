CREATE TABLE clients (
    id BIGSERIAL PRIMARY KEY,       -- ID autoincremental
    name VARCHAR(255) NOT NULL, -- Nombre del cliente, no puede ser nulo
    email VARCHAR(255) NOT NULL UNIQUE, -- Email del cliente, debe ser único
    phone VARCHAR(50),           -- Teléfono del cliente
    address TEXT                 -- Dirección del cliente
);
