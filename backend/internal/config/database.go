package config

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type NativeDatabase struct {
	DB *sql.DB
}

func Connect(conf *Configuration) (*NativeDatabase, error) {
	conn, err := sql.Open("postgres", conf.Database.GetConnectionString())

	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &NativeDatabase{conn}, nil
}

func (d *NativeDatabase) Close() {
	_ = d.DB.Close()
}
