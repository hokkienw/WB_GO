package database

import(
	"database/sql"
	"fmt"
	"log"
	"level_0/pkg/cache"
	_"github.com/lib/pq"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1"
	dbname   = "postgres"
)

func InitDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}

func CreateTableDB(db *sql.DB) error {
	create_table_query := "CREATE TABLE IF NOT EXISTS orders (id SERIAL PRIMARY KEY, message varchar NOT NULL);"
	_, err := db.Exec(create_table_query)
	return err
}

func InsertToDB(db *sql.DB, message string) error {
	insert_query := "INSERT INTO orders (message) VALUES ($1);"
	_, err := db.Exec(insert_query, message)
	return err
}

func GetFromDB(db *sql.DB, caache_instanse *cache.Cache) (error) {
	selectQuery := `
		SELECT id, message FROM orders;
	`

	rows, err := db.Query(selectQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	var id string
	var msg string

	for rows.Next() {
		if err := rows.Scan(&id, &msg); err != nil {
			return  err
		}
		caache_instanse.Set(id, msg, 20*time.Minute)
	}

	if err := rows.Err(); err != nil {
		return  err
	}

	return  nil
}

func CountDB(db *sql.DB) (int, error){
	count_query := "SELECT COUNT(*) FROM temp"
	var resault int
	err := db.QueryRow(count_query).Scan(&resault)
	if err != nil {
		return 0, err
	}

	return resault, nil
}