CREATE INDEX idx_transactions_user_date ON transactions(user_id, date);
CREATE INDEX idx_transactions_category_date ON transactions(category_id, user_id, date);

CREATE INDEX idx_budgets_user_year_month_category
    ON budgets(user_id, year, month, category_id);