package helpers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
)

type incrementRepo struct {
	db *pgxpool.Pool
}

func NewIncrementRepo(db *pgxpool.Pool) *incrementRepo {
	return &incrementRepo{
		db: db,
	}
}

func IncrementId(n int) string {
	t := "0000000"
	if len(strconv.Itoa(n+1)) == len(strconv.Itoa(n)) {
		return t[len(strconv.Itoa(n)):] + strconv.Itoa(n+1)
	}
	return t[len(strconv.Itoa(n))+1:] + strconv.Itoa(n+1)
}

func (i *incrementRepo) GetLast(ctx context.Context, tableName, columnName string) (string, error) {
	var last_Id string

	query := `SELECT MAX("` + columnName + `") FROM "` + tableName + `"`
	fmt.Println(query)
	err := i.db.QueryRow(ctx, query).Scan(&last_Id)

	if last_Id == "" {
		fmt.Println(last_Id)
		return "0000001", nil
	}
	if err != nil {
		fmt.Println(">>>>>>>>>>>>>>..", err)
		return "", err
	}
	number, err := strconv.Atoi(last_Id[2:])
	if err != nil {
		return "", err
	}
	incrementNumber := IncrementId(number)

	return incrementNumber, nil
}