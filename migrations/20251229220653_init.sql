-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscriptions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK(price > 0),
    start_date DATE NOT NULL,
    end_date DATE,
    is_ended BOOLEAN DEFAULT false,
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd
