-- +goose Up

CREATE TABLE tickets (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	event_id UUID REFERENCES users(id) ON DELETE CASCADE,
	price NUMERIC(10,2) NOT NULL,
	total_quantity INT NOT NULL,
	remaining_quantity INT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE events;
