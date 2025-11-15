package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"litespend-api/internal/model"
	"litespend-api/internal/repository"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidFileFormat = errors.New("invalid file format")
	ErrEmptyFile         = errors.New("file is empty")
)

type ImportService struct {
	transactionRepo       repository.TransactionRepository
	categoryRepo          repository.CategoryRepository
	prescribedExpanseRepo repository.PrescribedExpanseRepository
}

func NewImportService(
	transactionRepo repository.TransactionRepository,
	categoryRepo repository.CategoryRepository,
	prescribedExpanseRepo repository.PrescribedExpanseRepository,
) *ImportService {
	return &ImportService{
		transactionRepo:       transactionRepo,
		categoryRepo:          categoryRepo,
		prescribedExpanseRepo: prescribedExpanseRepo,
	}
}

func (s *ImportService) ParseExcelFile(fileData []byte) (model.ExcelFileStructure, error) {
	f, err := excelize.OpenReader(bytes.NewReader(fileData))
	if err != nil {
		return model.ExcelFileStructure{}, ErrInvalidFileFormat
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return model.ExcelFileStructure{}, ErrEmptyFile
	}

	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) == 0 {
		return model.ExcelFileStructure{}, ErrEmptyFile
	}

	columns := rows[0]
	if len(rows) < 2 {
		return model.ExcelFileStructure{
			Columns: columns,
			Rows:    0,
		}, nil
	}

	return model.ExcelFileStructure{
		Columns: columns,
		Rows:    len(rows) - 1,
	}, nil
}

func (s *ImportService) ImportData(ctx context.Context, logined model.User, fileData []byte, mapping model.ExcelColumnMapping) (model.ImportResult, error) {
	result := model.ImportResult{
		Errors: []string{},
	}

	f, err := excelize.OpenReader(bytes.NewReader(fileData))
	if err != nil {
		return result, ErrInvalidFileFormat
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) < 2 {
		return result, ErrEmptyFile
	}

	headers := rows[0]
	columnIndexes := s.getColumnIndexes(headers, mapping)

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		// Импорт транзакций
		if mapping.TransactionAmount != nil {
			if err := s.importTransaction(ctx, logined, row, columnIndexes, mapping); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Строка %d (транзакция): %v", i+1, err))
			} else {
				result.TransactionsCreated++
			}
		}

		// Импорт категорий
		if mapping.CategoryName != nil {
			created, err := s.importCategory(ctx, logined, row, columnIndexes, mapping)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Строка %d (категория): %v", i+1, err))
			} else if created {
				result.CategoriesCreated++
			}
		}

		// Импорт обязательных трат
		if mapping.PrescribedExpanseAmount != nil {
			if err := s.importPrescribedExpanse(ctx, logined, row, columnIndexes, mapping); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Строка %d (обязательная трата): %v", i+1, err))
			} else {
				result.PrescribedExpansesCreated++
			}
		}
	}

	return result, nil
}

type columnIndexes struct {
	transactionDescription       int
	transactionAmount            int
	transactionType              int
	transactionDate              int
	transactionCategory          int
	categoryName                 int
	categoryType                 int
	prescribedExpanseDescription int
	prescribedExpanseAmount      int
	prescribedExpanseFrequency   int
	prescribedExpanseDate        int
	prescribedExpanseCategory    int
}

func (s *ImportService) getColumnIndexes(headers []string, mapping model.ExcelColumnMapping) columnIndexes {
	indexes := columnIndexes{
		transactionDescription:       -1,
		transactionAmount:            -1,
		transactionType:              -1,
		transactionDate:              -1,
		transactionCategory:          -1,
		categoryName:                 -1,
		categoryType:                 -1,
		prescribedExpanseDescription: -1,
		prescribedExpanseAmount:      -1,
		prescribedExpanseFrequency:   -1,
		prescribedExpanseDate:        -1,
		prescribedExpanseCategory:    -1,
	}

	for i, header := range headers {
		if mapping.TransactionDescription != "" && strings.EqualFold(header, mapping.TransactionDescription) {
			indexes.transactionDescription = i
		}
		if mapping.TransactionAmount != nil && strings.EqualFold(header, *mapping.TransactionAmount) {
			indexes.transactionAmount = i
		}
		if mapping.TransactionType != nil && strings.EqualFold(header, *mapping.TransactionType) {
			indexes.transactionType = i
		}
		if mapping.TransactionDate != nil && strings.EqualFold(header, *mapping.TransactionDate) {
			indexes.transactionDate = i
		}
		if mapping.TransactionCategory != nil && strings.EqualFold(header, *mapping.TransactionCategory) {
			indexes.transactionCategory = i
		}
		if mapping.CategoryName != nil && strings.EqualFold(header, *mapping.CategoryName) {
			indexes.categoryName = i
		}
		if mapping.CategoryType != nil && strings.EqualFold(header, *mapping.CategoryType) {
			indexes.categoryType = i
		}
		if mapping.PrescribedExpanseDescription != "" && strings.EqualFold(header, mapping.PrescribedExpanseDescription) {
			indexes.prescribedExpanseDescription = i
		}
		if mapping.PrescribedExpanseAmount != nil && strings.EqualFold(header, *mapping.PrescribedExpanseAmount) {
			indexes.prescribedExpanseAmount = i
		}
		if mapping.PrescribedExpanseFrequency != nil && strings.EqualFold(header, *mapping.PrescribedExpanseFrequency) {
			indexes.prescribedExpanseFrequency = i
		}
		if mapping.PrescribedExpanseDate != nil && strings.EqualFold(header, *mapping.PrescribedExpanseDate) {
			indexes.prescribedExpanseDate = i
		}
		if mapping.PrescribedExpanseCategory != nil && strings.EqualFold(header, *mapping.PrescribedExpanseCategory) {
			indexes.prescribedExpanseCategory = i
		}
	}

	return indexes
}

func (s *ImportService) getCellValue(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func (s *ImportService) importTransaction(ctx context.Context, logined model.User, row []string, indexes columnIndexes, mapping model.ExcelColumnMapping) error {
	// Сумма обязательна для транзакции
	if indexes.transactionAmount < 0 {
		return errors.New("не указан столбец для суммы транзакции")
	}
	amountStr := s.getCellValue(row, indexes.transactionAmount)
	if amountStr == "" {
		return errors.New("сумма транзакции не может быть пустой")
	}

	amount, err := tryParseAmount(amountStr)
	if err != nil {
		return fmt.Errorf("неверный формат суммы: %v", err)
	}

	description := s.getCellValue(row, indexes.transactionDescription)
	if description == "" {
		description = "Импортированная транзакция"
	}

	transactionType := model.TransactionTypeExpanse
	if indexes.transactionType >= 0 {
		typeStr := strings.ToLower(s.getCellValue(row, indexes.transactionType))
		if strings.Contains(typeStr, "доход") || strings.Contains(typeStr, "income") || typeStr == "0" {
			transactionType = model.TransactionTypeIncome
		}
	}

	dateTime := time.Now()
	if indexes.transactionDate >= 0 {
		dateStr := s.getCellValue(row, indexes.transactionDate)
		if dateStr != "" {
			parsedDate, err := s.parseDate(dateStr)
			if err == nil {
				dateTime = parsedDate
			}
		}
	}

	var categoryID uint64 = 0
	if indexes.transactionCategory >= 0 {
		categoryName := s.getCellValue(row, indexes.transactionCategory)
		if categoryName != "" {
			categories, err := s.categoryRepo.GetList(ctx, int(logined.ID))
			if err == nil {
				for _, cat := range categories {
					if strings.EqualFold(cat.Name, categoryName) {
						categoryID = cat.ID
						break
					}
				}
				// Если категория не найдена, создаем её
				if categoryID == 0 {
					catType := model.TransactionTypeExpanse
					if transactionType == model.TransactionTypeIncome {
						catType = model.TransactionTypeIncome
					}
					newCatID, err := s.categoryRepo.Create(ctx, model.CreateCategoryRecord{
						UserID:    logined.ID,
						Name:      categoryName,
						Type:      catType,
						CreatedAt: time.Now(),
					})
					if err == nil {
						categoryID = uint64(newCatID)
					}
				}
			}
		}
	}

	// Если категория не указана, используем первую доступную категорию расходов/доходов
	if categoryID == 0 {
		categories, err := s.categoryRepo.GetListByType(ctx, int(logined.ID), transactionType)
		if err == nil && len(categories) > 0 {
			categoryID = categories[0].ID
		}
	}

	transaction := model.Transaction{
		UserID:      logined.ID,
		CategoryID:  categoryID,
		Description: description,
		Amount:      amount,
		Type:        transactionType,
		DateTime:    dateTime,
		CreatedAt:   time.Now(),
	}

	_, err = s.transactionRepo.Create(ctx, transaction)
	return err
}

func (s *ImportService) importCategory(ctx context.Context, logined model.User, row []string, indexes columnIndexes, mapping model.ExcelColumnMapping) (bool, error) {
	if indexes.categoryName < 0 {
		return false, errors.New("не указан столбец для названия категории")
	}

	name := s.getCellValue(row, indexes.categoryName)
	if name == "" {
		return false, errors.New("название категории не может быть пустым")
	}

	// Проверяем, не существует ли уже такая категория
	categories, err := s.categoryRepo.GetList(ctx, int(logined.ID))
	if err == nil {
		for _, cat := range categories {
			if strings.EqualFold(cat.Name, name) {
				// Категория уже существует, возвращаем false (не создана)
				return false, nil
			}
		}
	}

	catType := model.TransactionTypeExpanse
	if indexes.categoryType >= 0 {
		typeStr := strings.ToLower(s.getCellValue(row, indexes.categoryType))
		if strings.Contains(typeStr, "доход") || strings.Contains(typeStr, "income") || typeStr == "0" {
			catType = model.TransactionTypeIncome
		}
	}

	_, err = s.categoryRepo.Create(ctx, model.CreateCategoryRecord{
		UserID:    logined.ID,
		Name:      name,
		Type:      catType,
		CreatedAt: time.Now(),
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *ImportService) importPrescribedExpanse(ctx context.Context, logined model.User, row []string, indexes columnIndexes, mapping model.ExcelColumnMapping) error {
	// Сумма обязательна для обязательной траты
	if indexes.prescribedExpanseAmount < 0 {
		return errors.New("не указан столбец для суммы обязательной траты")
	}

	amountStr := s.getCellValue(row, indexes.prescribedExpanseAmount)
	if amountStr == "" {
		return errors.New("сумма обязательной траты не может быть пустой")
	}

	amount, err := tryParseAmount(amountStr)
	if err != nil {
		return fmt.Errorf("неверный формат суммы: %v", err)
	}

	description := s.getCellValue(row, indexes.prescribedExpanseDescription)
	if description == "" {
		description = "Импортированная обязательная трата"
	}

	frequency := model.FrequencyTypeMonthly
	if indexes.prescribedExpanseFrequency >= 0 {
		freqStr := strings.ToLower(s.getCellValue(row, indexes.prescribedExpanseFrequency))
		if strings.Contains(freqStr, "ежедневно") || strings.Contains(freqStr, "daily") || freqStr == "1" {
			frequency = model.FrequencyTypeDaily
		} else if strings.Contains(freqStr, "еженедельно") || strings.Contains(freqStr, "weekly") || freqStr == "2" {
			frequency = model.FrequencyTypeWeekly
		} else if strings.Contains(freqStr, "ежеквартально") || strings.Contains(freqStr, "quarterly") || freqStr == "3" {
			frequency = model.FrequencyTypeQuarterly
		}
	}

	dateTime := time.Now()
	if indexes.prescribedExpanseDate >= 0 {
		dateStr := s.getCellValue(row, indexes.prescribedExpanseDate)
		if dateStr != "" {
			parsedDate, err := s.parseDate(dateStr)
			if err == nil {
				dateTime = parsedDate
			}
		}
	}

	var categoryID uint64 = 0
	if indexes.prescribedExpanseCategory >= 0 {
		categoryName := s.getCellValue(row, indexes.prescribedExpanseCategory)
		if categoryName != "" {
			categories, err := s.categoryRepo.GetList(ctx, int(logined.ID))
			if err == nil {
				for _, cat := range categories {
					if strings.EqualFold(cat.Name, categoryName) {
						categoryID = cat.ID
						break
					}
				}
				// Если категория не найдена, создаем её
				if categoryID == 0 {
					newCatID, err := s.categoryRepo.Create(ctx, model.CreateCategoryRecord{
						UserID:    logined.ID,
						Name:      categoryName,
						Type:      model.TransactionTypeExpanse,
						CreatedAt: time.Now(),
					})
					if err == nil {
						categoryID = uint64(newCatID)
					}
				}
			}
		}
	}

	// Если категория не указана, используем первую доступную категорию расходов
	if categoryID == 0 {
		categories, err := s.categoryRepo.GetListByType(ctx, int(logined.ID), model.TransactionTypeExpanse)
		if err == nil && len(categories) > 0 {
			categoryID = categories[0].ID
		}
	}

	prescribedExpanse := model.CreatePrescribedExpanseRecord{
		UserID:      logined.ID,
		CategoryID:  categoryID,
		Description: description,
		Frequency:   frequency,
		Amount:      amount,
		DateTime:    dateTime,
		CreatedAt:   time.Now(),
	}

	_, err = s.prescribedExpanseRepo.Create(ctx, prescribedExpanse)
	return err
}

func (s *ImportService) parseDate(dateStr string) (time.Time, error) {
	// Пробуем различные форматы дат
	formats := []string{
		"2006-01-02",
		"02.01.2006",
		"02/01/2006",
		"2006-01-02 15:04:05",
		"02.01.2006 15:04:05",
		"02/01/2006 15:04:05",
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	// Пробуем парсить как число Excel (дни с 1900-01-01)
	if days, err := strconv.ParseFloat(dateStr, 64); err == nil {
		excelEpoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
		return excelEpoch.AddDate(0, 0, int(days)), nil
	}

	return time.Time{}, fmt.Errorf("не удалось распарсить дату: %s", dateStr)
}

func tryParseAmount(amountStr string) (decimal.Decimal, error) {
	amountSplit := strings.Split(amountStr, " ")
	if len(amountSplit) > 1 {
		amountStr = strings.Join(amountSplit[:len(amountSplit)-1], "")
	}

	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("неверный формат суммы: %v", err)

	}

	return amount, nil
}
