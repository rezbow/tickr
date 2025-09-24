-- +goose Up
CREATE TABLE payment (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	ticket_id UUID REFERENCES tickets(id) ON DELETE CASCADE,
	quantity INT NOT NULL,
    paid_amount BIGINT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS payment;
