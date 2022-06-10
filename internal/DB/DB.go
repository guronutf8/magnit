package DB

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"magnit/internal/conf"
	"strconv"
)

type DB struct {
	conn *sqlx.DB
}

// Init Инициализация соединения с БД
func Init(cfg conf.Config) *DB {
	conn, err := sqlx.Open("postgres", "user="+cfg.DBLogin+" password="+cfg.DBPassword+" dbname=demo sslmode=disable")

	if err != nil {
		panic(err)
		//os.Exit(1)
	}

	conn.SetMaxIdleConns(10)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxIdleTime(5 * time.Second)
	conn.SetConnMaxLifetime(5 * time.Second)

	db := DB{conn: conn}

	return &db
}

// AddUrl Добавление ссылки в бд. Отдаем ID строки, потом это ID превратим в короткую ссылку
func (r DB) AddUrl(domain string, path string) int {
	var lastID struct{ Id int }
	// todo надо ли проверять на дубли? доделаю если скажут
	err := r.conn.QueryRowx("insert into \"ShortUrl\" (domain, path) VALUES ('" + domain + "','" + path + "') RETURNING id;").StructScan(&lastID)

	if err != nil {
		fmt.Println("err:", err)
	}
	return lastID.Id
}

// GetOriginal Получаем оригинальную ссылку, ищем сразу ID, если ее нет то ошибка
func (r DB) GetOriginal(index uint64) (string, error) {
	var data struct {
		Domain string
		Path   string
	}
	// todo проверять на дубли
	err := r.conn.QueryRowx("select Domain,Path from \"ShortUrl\" where id = " + strconv.Itoa(int(index))).StructScan(&data)
	if err != nil {
		fmt.Println("err:", err)
	}

	if data.Domain != "" && data.Path != "" {
		return data.Domain + data.Path, nil
	}

	return "", errors.New("not found")
}
