-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS manager_kpi (
    manager_id UUID NOT NULL,
    shift_date DATE NOT NULL,
    plan_revenue DECIMAL(15,2) NOT NULL,
    fact_revenue DECIMAL(15,2) NOT NULL DEFAULT 0,
    PRIMARY KEY (manager_id, shift_date)
);

CREATE TABLE IF NOT EXISTS sales_reports (
    id BIGSERIAL PRIMARY KEY,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    total_revenue DECIMAL(15,2) NOT NULL,
    total_orders INT NOT NULL,
    average_check DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sales_reports;
DROP TABLE IF EXISTS manager_kpi;
-- +goose StatementEnd
