-- +goose Up
CREATE TABLE IF NOT EXISTS events(
    id SERIAL PRIMARY KEY,  
    event_name TEXT,
    all_seats INT NOT NULL,
    booked INT NOT NULL DEFAULT 0,
    event_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bookings(
    id SERIAL PRIMARY KEY,
    event_id  INT NOT NULL REFERENCES events(id),
    status TEXT NOT NULL CHECK(status IN ('created', 'paid')),
    telegram_id INT, 
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bookings_event_id ON bookings(event_id);


-- +goose Down
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS bookings;
