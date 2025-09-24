-- +goose Up

CREATE TABLE tickets (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	event_id UUID REFERENCES events(id) ON DELETE CASCADE,
	price BIGINT NOT NULL,
    total_quantities INT NOT NULL,
    remaining_quantities INT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS tickets;
