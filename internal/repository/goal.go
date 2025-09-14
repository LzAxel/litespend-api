package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"litespend-api/internal/model"
)

type GoalRepositoryPostgres struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewGoalRepositoryPostgres(db *sqlx.DB) GoalRepositoryPostgres {
	return GoalRepositoryPostgres{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r GoalRepositoryPostgres) Create(ctx context.Context, goal model.CreateGoalRecord) (int, error) {
	var createdID int
	sql := `INSERT INTO goals(
			user_id,      
			name,         
			target_amount,
			start_amount, 
			frequency,    
			deadline_date,
			created_at   
		) VALUE ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.db.GetContext(ctx, &createdID, sql, goal.UserID, goal.Name, goal.TargetAmount, goal.StartAmount, goal.Frequency, goal.DeadlineDate, goal.CreatedAt)
	if err != nil {
		return 0, err
	}

	return createdID, nil
}

func (r GoalRepositoryPostgres) Replace(ctx context.Context, id int, dto model.UpdateGoalRecord) error {
	sql := `UPDATE goals(name, target_amount, start_amount, frequency, deadline_date) WHERE id=$1 VALUES ($2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, sql, id, dto.Name, dto.Name, dto.TargetAmount, dto.StartAmount, dto.Frequency, dto.DeadlineDate)
	return err
}

func (r GoalRepositoryPostgres) Delete(ctx context.Context, id int) error {
	sql := `DELETE FROM goals WHERE id=$1`

	_, err := r.db.ExecContext(ctx, sql, id)
	return err
}

func (r GoalRepositoryPostgres) GetListByUserID(ctx context.Context, userID int) ([]model.Goal, error) {
	var goals = make([]model.Goal, 0)

	sql := `SELECT * FROM goals WHERE user_id=$1`

	err := r.db.SelectContext(ctx, goals, sql, userID)
	if err != nil {
		return goals, err
	}

	return goals, nil
}

func (r GoalRepositoryPostgres) GetByID(ctx context.Context, id int) (model.Goal, error) {
	var goal model.Goal

	sql := `SELECT * FROM goals WHERE id=$1`

	err := r.db.GetContext(ctx, &goal, sql, id)
	if err != nil {
		return goal, err
	}

	return goal, nil
}
