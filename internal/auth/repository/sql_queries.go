package repository

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/shlembo598/text-lexicon-go/internal/models"
)

func createUserQuery(user *models.User) (string, []interface{}, error) {
	return sq.Insert("users").Columns(
		"first_name", "last_name", "email", "password", "created_at", "updated_at", "login_date",
	).Values(
		&user.FirstName, &user.LastName, &user.Email,
		&user.Password, time.Now(), time.Now(), time.Now(),
	).PlaceholderFormat(sq.Dollar).ToSql()
}

func updateUserQuery(user *models.User) (string, []interface{}, error) {
	return sq.Update("users").Set(
		"first_name", sq.Expr(
			"COALESCE(NULLIF(?, ''), first_name)",
			user.FirstName,
		),
	).Set(
		"last_name", sq.Expr(
			"COALESCE(NULLIF(?, ''), last_name)",
			user.LastName,
		),
	).Set(
		"avatar", sq.Expr(
			"COALESCE(NULLIF(?, ''), avatar)",
			user.Avatar,
		),
	).Set(
		"country", sq.Expr(
			"COALESCE(NULLIF(?, ''), country)",
			user.Country,
		),
	).Set(
		"updated_at", time.Now(),
	).Where(
		"user_id", user.UserID,
	).PlaceholderFormat(sq.Dollar).ToSql()
}

func deleteUserQuery(userID uuid.UUID) (string, []interface{}, error) {
	return sq.Delete("users").Where("user_id", userID).PlaceholderFormat(sq.Dollar).ToSql()
}

func getUserQuery(userID uuid.UUID) (string, []interface{}, error) {
	return sq.Select("users").Columns(
		"user_id", "first_name", "last_name", "email", "avatar", "country", "created_at", "updated_at",
		"login_date",
	).Where("user_id", userID).PlaceholderFormat(sq.Dollar).ToSql()
}

func findUserByEmail(email string) (string, []interface{}, error) {
	return sq.Select("users").Columns(
		"user_id", "first_name", "last_name", "email", "avatar", "country", "created_at", "updated_at",
		"login_date", "password",
	).Where("email", email).PlaceholderFormat(sq.Dollar).ToSql()
}
