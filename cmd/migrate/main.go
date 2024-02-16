package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/lud0m4n/WebAppDev/internal/model"
)

const numRecords = 1000000

func main() {
	dsn := "host=127.0.0.1 port=5433 user=postgres password=1111 dbname=Geo_periods sslmode=disable"
	// Инициализация подключения к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	_ = godotenv.Load()
	// Явно мигрировать только нужные таблицы
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic("cant migrate db")
	}
	db.AutoMigrate(&model.Period{})

	// Загрузка 10 тысяч записей
	err = insertRandomPeriods(db, numRecords)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d записей успешно загружены в таблицу baggages\n", numRecords)
}

// Вставка случайных записей в таблицу baggages
func insertRandomPeriods(db *gorm.DB, numRecords int) error {
	for i := 1; i <= numRecords; i++ {
		period := model.Period{
			Name:        generatePeriodCode(i),
			Description: fmt.Sprintf("Something"),
			Age:         fmt.Sprintf("%dx%dx%d", rand.Intn(100), rand.Intn(100), rand.Intn(100)),
			Photo:       fmt.Sprintf("http://localhost:9000/images-bucket/period/1/image"),
		}

		// Начало транзакции
		tx := db.Begin()

		if err := tx.Create(&period).Error; err != nil {
			// Откат транзакции при ошибке
			tx.Rollback()
			return err
		}

		// Фиксация транзакции
		tx.Commit()

		// Небольшая задержка для имитации реального использования
		time.Sleep(time.Millisecond)
	}

	return nil
}

func generatePeriodCode(i int) string {
	rand.Seed(time.Now().UnixNano())

	// Генерация трех случайных букв капсом
	randomLetters := make([]byte, 3)
	for j := range randomLetters {
		randomLetters[j] = byte('A' + rand.Intn('Z'-'A'+1))
	}

	return fmt.Sprintf("%s%d", randomLetters, i)
}
