-- +goose Up
CREATE TABLE events (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	title VARCHAR(255) NOT NULL,
	description TEXT ,
	venue TEXT NOT NULL ,
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	start_time TIMESTAMP NOT NULL,
	end_time TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE events;
