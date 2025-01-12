-- Таблица клиентов
CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    inn VARCHAR(50)
);

-- Таблица товаров
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    unit VARCHAR(50) NOT NULL
);

-- Таблица заказов
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    client_id INTEGER REFERENCES clients(id),
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    number VARCHAR(50) NOT NULL UNIQUE,
    total_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    is_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Таблица позиций заказа
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id),
    quantity DECIMAL(15,3) NOT NULL CHECK (quantity >= 0),
    price DECIMAL(15,2) NOT NULL CHECK (price >= 0),
    line_amount DECIMAL(15,2) NOT NULL DEFAULT 0
);

-- Таблица агрегированных сумм по клиентам
CREATE TABLE orders_by_client (
    client_id INTEGER REFERENCES clients(id) PRIMARY KEY,
    orders_sum DECIMAL(15,2) NOT NULL DEFAULT 0
);

-- Триггер для автоматического расчета line_amount
CREATE OR REPLACE FUNCTION calculate_line_amount()
RETURNS TRIGGER AS $$
BEGIN
    NEW.line_amount = NEW.quantity * NEW.price;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_calculate_line_amount
BEFORE INSERT OR UPDATE ON order_items
FOR EACH ROW EXECUTE FUNCTION calculate_line_amount(); 