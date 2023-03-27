package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/exceptions"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/models"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresDB struct {
	pool *pgxpool.Pool
}

func NewPostgresDB(dsn string) *PostgresDB {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	log.Printf("Connected!")
	return &PostgresDB{pool: pool}
}

func (pdb *PostgresDB) Close() error {
	if pdb.pool == nil {
		return nil
	}

	pdb.pool.Close()

	return nil
}

func (pdb *PostgresDB) MigrateToTheLatestSchema(dsn string, sourceDDL string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	databasePath, err := url.Parse(dsn)
	if err != nil {
		panic(err)
	}
	databaseName := databasePath.Path[1:]

	m, mErr := migrate.NewWithDatabaseInstance(sourceDDL, databaseName, driver)
	if mErr != nil {
		panic(mErr)
	}

	errOnMigrate := m.Up()

	if errOnMigrate != nil && !errors.Is(errOnMigrate, migrate.ErrNoChange) {
		panic(errOnMigrate)
	}

	return nil
}

func (pdb *PostgresDB) Ping() bool {
	return pdb.pool.Ping(context.Background()) == nil
}

func (pdb *PostgresDB) CreateUser(user *models.User) error {

	_, err := pdb.pool.Exec(context.Background(), "INSERT INTO users (id, username, pass, token, token_expires) VALUES ($1, $2, $3, $4, $5) on conflict (username) do nothing", user.ID, user.Username, user.Password, user.Token, user.TokenExpires)

	if err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
		return exceptions.ErrDuplicatePK
	}

	return nil
}

func (pdb *PostgresDB) ReadUser(username string) (*models.User, error) {
	var user models.User
	err := pdb.pool.QueryRow(context.Background(), "SELECT id, username, pass, token, token_expires FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password, &user.Token, &user.TokenExpires)
	if err != nil {
		log.Println(err)
		return &models.User{}, err
	}
	return &user, err
}

func (pdb *PostgresDB) CheckOrder(orderNum string) (*models.Order, error) {
	var order models.Order
	err := pdb.pool.QueryRow(context.Background(), "SELECT user_id, order_num FROM orders WHERE order_num=$1", orderNum).Scan(&order.UserID, &order.OrderNum)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &order, nil
		}
		log.Println(err)
		return &order, err
	}

	return &order, nil
}

func (pdb *PostgresDB) CreateOrder(order *models.Order) error {
	order.Status = "NEW"
	_, err := pdb.pool.Exec(context.Background(), "INSERT INTO orders (user_id, order_num, status, created_at) VALUES ($1, $2, $3, now()) ON CONFLICT (order_num) DO NOTHING", order.UserID, order.OrderNum, order.Status)

	if err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
		return exceptions.ErrDuplicatePK
	}

	return nil
}

func (pdb *PostgresDB) ReadOrders(userID uuid.UUID) ([]*models.OrderDB, error) {
	orders := make([]*models.OrderDB, 0)
	rows, err := pdb.pool.Query(context.Background(), "SELECT order_num, accrual, status, created_at FROM orders WHERE (user_id=$1 AND created_at IS NOT NULL) ORDER BY created_at;", userID)
	if err != nil {
		log.Println(err)
		return orders, err
	}

	for rows.Next() {
		var order models.OrderDB

		err := rows.Scan(&order.OrderNum, &order.Accrual, &order.Status, &order.Created)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	if len(orders) == 0 {
		return orders, exceptions.ErrNoValues
	}

	return orders, nil
}

func (pdb *PostgresDB) ReadBalance(userID uuid.UUID) (*models.Balance, error) {
	var balance models.Balance
	err := pdb.pool.QueryRow(context.Background(), `
        SELECT
            COALESCE(SUM(CASE WHEN orders.created_at IS NULL AND withdrawn_at IS NOT NULL THEN withdrawal ELSE 0 END), 0)::NUMERIC(10, 2) withdraw,
            COALESCE(SUM(CASE WHEN orders.created_at IS NOT NULL AND withdrawn_at IS NULL THEN accrual ELSE 0 END), 0)::NUMERIC(10, 2) - COALESCE(SUM(CASE WHEN orders.created_at IS NULL AND withdrawn_at IS NOT NULL THEN withdrawal ELSE 0 END), 0)::NUMERIC(10, 2) "current"
        FROM
            users
            LEFT JOIN orders ON orders.user_id = users.id
        WHERE
            users.id = $1
        GROUP BY
            users.id;
    `, userID).Scan(&balance.Withdraw, &balance.Current)
	if err != nil {
		log.Println(err)
		return &balance, err
	}
	return &balance, err
}

func (pdb *PostgresDB) CreateWithdrawal(withdraw *models.Withdraw) error {
	_, err := pdb.pool.Exec(context.Background(), "INSERT INTO orders (user_id, order_num, withdrawal, withdrawn_at) VALUES ($1, $2, $3, now()) ON CONFLICT (order_num) DO NOTHING", withdraw.UserID, withdraw.OrderNum, withdraw.Withdraw)

	if err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
		return exceptions.ErrDuplicatePK
	}

	return nil
}

func (pdb *PostgresDB) ReadAllWithdrawals(userID uuid.UUID) ([]*models.WithdrawDB, error) {
	var withdrawals []*models.WithdrawDB
	rows, err := pdb.pool.Query(context.Background(), "SELECT order_num, withdrawal, withdrawn_at FROM orders WHERE (user_id=$1 AND withdrawn_at IS NOT NULL) ORDER BY withdrawn_at", userID)
	if err != nil {
		log.Println(err)
		return withdrawals, err
	}

	for rows.Next() {
		var withdrawal models.WithdrawDB

		err := rows.Scan(&withdrawal.OrderNum, &withdrawal.Withdraw, &withdrawal.Created)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, &withdrawal)
	}

	if len(withdrawals) == 0 {
		return withdrawals, exceptions.ErrNoValues
	}

	return withdrawals, nil
}

func (pdb *PostgresDB) ReadOrdersForProcessing() ([]string, error) {
	var orders []string

	rows, err := pdb.pool.Query(context.Background(), "SELECT order_num FROM orders WHERE (status IN ('NEW','PROCESSING') AND created_at IS NOT NULL) ORDER BY created_at")
	if err != nil {
		log.Println(err)
		return orders, err
	}

	for rows.Next() {
		var order string

		err := rows.Scan(&order)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return orders, exceptions.ErrNoValues
	}

	return orders, nil
}

func (pdb *PostgresDB) UpdateOrder(order models.ProcessingOrder) {
	var query string

	if order.Accrual != nil {
		query = fmt.Sprintf("UPDATE orders SET status = '%s', accrual = '%f' WHERE order_num = '%s';", order.Status, order.Accrual.(float64), order.OrderNum)
	} else {
		query = fmt.Sprintf("UPDATE orders SET status = '%s' WHERE order_num = '%s';", order.Status, order.OrderNum)
	}

	res, err := pdb.pool.Exec(context.Background(), query)
	if err != nil {
		log.Printf("update failed, err:%v\n", err)
		return
	}

	fmt.Printf("update success, affected rows:%d\n", res.RowsAffected())
}
