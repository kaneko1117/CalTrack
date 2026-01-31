package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
)

// RunMigrations はマイグレーションを実行する
func RunMigrations() error {
	// DB接続情報を取得
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		getEnv("DB_USER", "caltrack"),
		getEnv("DB_PASSWORD", "caltrack"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_NAME", "caltrack"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("マイグレーション用DB接続エラー: %w", err)
	}
	defer db.Close()

	// マイグレーションソースを設定
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	// マイグレーションを実行
	n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("マイグレーション実行エラー: %w", err)
	}

	log.Printf("マイグレーション完了: %d件適用", n)
	return nil
}
