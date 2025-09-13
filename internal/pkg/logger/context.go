package logger

import (
	"context"
	"fmt"
	"litespend-api/internal/model"
	"log/slog"
	"reflect"
)

type loggerCtx struct {
	UserID   int
	UserRole string

	Username string
	Message  string
	Password string
}

type keyType int

const key = keyType(0)

func WithLogUserID(ctx context.Context, userID int) context.Context {
	return withLog(ctx, "UserID", userID)
}

func WithLogPassword(ctx context.Context, password string) context.Context {
	return withLog(ctx, "Password", password)
}

func WithLogMessage(ctx context.Context, message string) context.Context {
	return withLog(ctx, "Error", message)
}

func WithLogUsername(ctx context.Context, username string) context.Context {
	return withLog(ctx, "Username", username)
}

func WithUserRole(ctx context.Context, role model.UserRole) context.Context {
	return withLog(ctx, "UserRole", role.String())
}

func withLog(ctx context.Context, fieldName string, value any) context.Context {
	loggerContext := &loggerCtx{}
	if c, ok := ctx.Value(key).(loggerCtx); ok {
		loggerContext = &c
	}

	loggerContextVal := reflect.ValueOf(loggerContext).Elem()
	if loggerContextVal.Kind() != reflect.Struct {
		slog.Error("LOGGER", "logger context must be a pointer to a struct")
		return ctx
	}

	field := loggerContextVal.FieldByName(fieldName)

	if !field.IsValid() {
		slog.Error("LOGGER", fmt.Errorf("field %s does not exist", fieldName))
		return ctx
	}
	if !field.CanSet() {
		slog.Error("LOGGER", fmt.Errorf("field %s cannot be set", fieldName))
		return ctx
	}

	newVal := reflect.ValueOf(value)

	// Проверяем, совместим ли тип нового значения с типом поля
	if newVal.Type().AssignableTo(field.Type()) {
		field.Set(newVal)
	} else {
		slog.Error("LOGGER", fmt.Errorf("value type %v is not assignable to field type %v", newVal.Type(), field.Type()))
	}

	return context.WithValue(ctx, key, *loggerContext)
}
