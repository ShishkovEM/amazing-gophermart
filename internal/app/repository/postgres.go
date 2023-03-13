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
	dsn  string
	pool *pgxpool.Pool
}

func NewPostgresDB(dsn string) *PostgresDB {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	log.Printf("Connected!")
	return &PostgresDB{dsn: dsn, pool: pool}
}

func (pdb *PostgresDB) GetConn() (*pgxpool.Conn, error) {
	if pdb.dsn == "" {
		return nil, exceptions.ErrNoDatabaseDSN
	}

	conn, err := pdb.pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v\n", err)
		return nil, err
	}

	return conn, nil
}

func (pdb *PostgresDB) Close() error {
	if pdb.pool == nil {
		return nil
	}

	pdb.pool.Close()

	return nil
}

func (pdb *PostgresDB) MigrateToTheLatestSchema() error {
	db, err := sql.Open("postgres", pdb.dsn)
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

	databasePath, err := url.Parse(pdb.dsn)
	if err != nil {
		panic(err)
	}
	databaseName := databasePath.Path[1:]

	m, mErr := migrate.NewWithDatabaseInstance(
		"file://./schema",
		databaseName, driver,
	)
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

	_, err := pdb.pool.Exec(context.Background(), "INSERT INTO users (id, username, pass, cookie, cookie_expires) VALUES ($1, $2, $3, $4, $5) on conflict (username) do nothing", user.ID, user.Username, user.Password, user.Cookie, user.CookieExpires)

	if err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
		return exceptions.ErrDuplicatePK
	}

	return nil
}

func (pdb *PostgresDB) ReadUser(username string) (*models.User, error) {
	var user models.User
	err := pdb.pool.QueryRow(context.Background(), "SELECT id, username, pass, cookie, cookie_expires FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password, &user.Cookie, &user.CookieExpires)
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
	_, err := pdb.pool.Exec(context.Background(), "INSERT INTO orders (user_id, order_num, status) VALUES ($1, $2, $3) ON CONFLICT (order_num) DO NOTHING", order.UserID, order.OrderNum, order.Status)

	if err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
		return exceptions.ErrDuplicatePK
	}

	return nil
}

func (pdb *PostgresDB) ReadOrders(userID uuid.UUID) ([]*models.OrderDB, error) {
	orders := make([]*models.OrderDB, 0)
	rows, err := pdb.pool.Query(context.Background(), "SELECT order_num, accrual, status, created_at FROM orders where user_id=$1 order by created_at;", userID)
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
	err := pdb.pool.QueryRow(context.Background(), "SELECT withdraw, \"current\" FROM balance WHERE user_id=$1;", userID).Scan(&balance.Withdraw, &balance.Current)
	if err != nil {
		log.Println(err)
		return &balance, err
	}

	return &balance, nil
}

func (pdb *PostgresDB) CreateWithdrawal(withdraw *models.Withdraw) error {
	_, err := pdb.pool.Exec(context.Background(), "INSERT INTO withdrawals (user_id, order_num, withdraw) VALUES ($1, $2, $3) ON CONFLICT (order_num) DO NOTHING", withdraw.UserID, withdraw.OrderNum, withdraw.Withdraw)

	if err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
		return exceptions.ErrDuplicatePK
	}

	return nil
}

func (pdb *PostgresDB) ReadAllWithdrawals(userID uuid.UUID) ([]*models.WithdrawDB, error) {
	var withdrawals []*models.WithdrawDB
	rows, err := pdb.pool.Query(context.Background(), "SELECT order_num, withdraw, created_at FROM withdrawals WHERE user_id=$1 ORDER BY created_at", userID)
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

	rows, err := pdb.pool.Query(context.Background(), "SELECT order_num FROM orders WHERE status IN ('NEW','PROCESSING') ORDER BY created_at")
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
