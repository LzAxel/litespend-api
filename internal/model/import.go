package model

type ExcelColumnMapping struct {
	// Транзакции
	TransactionDescription *string `json:"transaction_description,omitempty"` // столбец для описания транзакции
	TransactionAmount      string  `json:"transaction_amount,omitempty"`      // столбец для суммы транзакции
	TransactionDate        *string `json:"transaction_date,omitempty"`        // столбец для даты транзакции
	TransactionCategory    *string `json:"transaction_category,omitempty"`    // столбец для категории транзакции

	// Категории
	CategoryName *string `json:"category_name,omitempty"` // столбец для названия категории
	CategoryType *string `json:"category_type,omitempty"` // столбец для типа категории

	// Бюджеты
	BudgetName   *string `json:"budget_name,omitempty"`
	BudgetAmount *string `json:"budget_amount,omitempty"`
}

type ExcelFileStructure struct {
	Columns []string `json:"columns"` // названия столбцов из первой строки
	Rows    int      `json:"rows"`    // количество строк данных (без заголовка)
}

type ImportRequest struct {
	Mapping ExcelColumnMapping `json:"mapping" binding:"required"`
}

type ImportResult struct {
	TransactionsCreated int      `json:"transactions_created"`
	CategoriesCreated   int      `json:"categories_created"`
	BudgetsCreated      int      `json:"budgets_created"`
	Errors              []string `json:"errors,omitempty"`
}
