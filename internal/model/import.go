package model

type ExcelColumnMapping struct {
	// Транзакции
	TransactionDescription string  `json:"transaction_description,omitempty"` // столбец для описания транзакции
	TransactionAmount      *string `json:"transaction_amount,omitempty"`      // столбец для суммы транзакции
	TransactionType        *string `json:"transaction_type,omitempty"`        // столбец для типа транзакции (доход/расход)
	TransactionDate        *string `json:"transaction_date,omitempty"`        // столбец для даты транзакции
	TransactionCategory    *string `json:"transaction_category,omitempty"`    // столбец для категории транзакции

	// Категории
	CategoryName *string `json:"category_name,omitempty"` // столбец для названия категории
	CategoryType *string `json:"category_type,omitempty"` // столбец для типа категории

	// Обязательные траты
	PrescribedExpanseDescription string  `json:"prescribed_expanse_description,omitempty"` // столбец для описания обязательной траты
	PrescribedExpanseAmount      *string `json:"prescribed_expanse_amount,omitempty"`      // столбец для суммы обязательной траты
	PrescribedExpanseFrequency   *string `json:"prescribed_expanse_frequency,omitempty"`   // столбец для частоты обязательной траты
	PrescribedExpanseDate        *string `json:"prescribed_expanse_date,omitempty"`        // столбец для даты обязательной траты
	PrescribedExpanseCategory    *string `json:"prescribed_expanse_category,omitempty"`    // столбец для категории обязательной траты
}

type ExcelFileStructure struct {
	Columns []string `json:"columns"` // названия столбцов из первой строки
	Rows    int      `json:"rows"`    // количество строк данных (без заголовка)
}

type ImportRequest struct {
	Mapping ExcelColumnMapping `json:"mapping" binding:"required"`
}

type ImportResult struct {
	TransactionsCreated       int      `json:"transactions_created"`
	CategoriesCreated         int      `json:"categories_created"`
	PrescribedExpansesCreated int      `json:"prescribed_expanses_created"`
	Errors                    []string `json:"errors,omitempty"`
}
