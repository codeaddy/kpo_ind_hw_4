CREATE TABLE order_dish (
                            id SERIAL PRIMARY KEY,
                            order_id INT NOT NULL,
                            dish_id INT NOT NULL,
                            quantity INT NOT NULL,
                            price DECIMAL(10, 2) NOT NULL,
                            FOREIGN KEY (order_id) REFERENCES "order"(id),
                            FOREIGN KEY (dish_id) REFERENCES "dish"(id)
);
