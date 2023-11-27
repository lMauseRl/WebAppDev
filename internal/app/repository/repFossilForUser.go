package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/lud0m4n/WebAppDev/internal/app/ds"
)

func (r *Repository) GetFossilForUser(searchSpecies, startFormationDate, endFormationDate, fossilStatus string, userID uint) ([]ds.FossilPeriod, error) {
	searchSpecies = strings.ToUpper(searchSpecies + "%")
	fossilStatus = strings.ToLower(fossilStatus + "%")

	// Построение основного запроса для получения ископаемых.
	query := r.db.Table("fossil").
		Select("DISTINCT fossil.id_fossil, fossil.genus, fossil.species, fossil.creation_date, fossil.formation_date, fossil.completion_date, fossil.status, users.full_name").
		Joins("JOIN fossilperiod ON fossil.id_fossil = fossilperiods.id_fossil").
		Joins("JOIN periods ON periods.id_period = fossilperiods.period_id").
		Joins("JOIN users ON users.user_id = fossil.user_id").
		Where("fossil.status LIKE ? AND fossil.species LIKE ? AND fossil.user_id = ? AND fossil.status != ?", fossilStatus, searchSpecies, userID, ds.FOSSIL_STATUS_DELETED)

	// Добавление условия фильтрации по дате формирования, если она указана.
	if startFormationDate != "" && endFormationDate != "" {
		query = query.Where("fossil.formation_date BETWEEN ? AND ?", startFormationDate, endFormationDate)
	}

	// Выполнение запроса и сканирование результатов в слайс fossil.
	var fossils []ds.FossilPeriod
	if err := query.Find(&fossils).Error; err != nil {
		return nil, errors.New("ошибка получения ископаемых")
	}
	return fossils, nil
}

func (r *Repository) GetFossilByIDForUser(fossilID int, userID uint) (map[string]interface{}, error) {
	var fossil map[string]interface{}
	// Получение информации о ископаемых по fossilID.
	if err := r.db.
		Table("fossil").
		Select("fossil.id_fossil, fossil.flight_number, fossil.creation_date, fossil.formation_date, fossil.completion_date, fossil.status").
		Where("fossil.status != ? AND fossil.id_fossil = ? AND fossil.user_id = ?", ds.FOSSIL_STATUS_DELETED, fossilID, userID).
		Scan(&fossil).Error; err != nil {
		return nil, errors.New("ошибка получения останков по ИД")
	}

	// Получение периодов по указанному fossilID.
	periods, err := r.GetPeriodsBySpecies(fossil["species"].(string))
	if err != nil {
		return nil, err
	}
	// Добавление информации о периоде в поле "periods" внутри останков.
	fossil["periods"] = periods

	return fossil, nil
}

func (r *Repository) DeleteFossilForUser(fossilID int, userID uint) error {
	// Проверяем, существует ли указанные останки в базе данных
	var fossil ds.Fossil
	if err := r.db.First(&fossil, fossilID).Error; err != nil {
		return errors.New("данные останки не существуют")
	}

	// Проверяем, что пользователь является создателем этого останка
	if fossil.UserID != userID {
		return errors.New("пользователь не является создателем этого останка")
	}

	// Начинаем транзакцию для атомарности операций
	tx := r.db.Begin()

	// Удаляем связанные записи из таблицы-множества (fossilperiods)
	if err := tx.Where("id_fossil = ?", fossilID).Delete(&ds.FossilPeriod{}).Error; err != nil {
		tx.Rollback()
		return errors.New("ошибка удаления связей из таблицы-множества")
	}

	// Обновляем статус останков на "удален" с использованием GORM
	err := r.db.Model(&ds.Fossil{}).Where("id_fossil = ?", fossilID).Update("status", ds.FOSSIL_STATUS_DELETED).Error
	if err != nil {
		return errors.New("ошибка обновления статуса на удален")
	}
	// Фиксируем транзакцию
	tx.Commit()

	return nil
}

func (r *Repository) UpdateFossilForUser(fossilID int, userID uint, updatedFossil *ds.Fossil) error {
	// Проверяем, существует ли указанные останки в базе данных
	var fossil ds.Fossil
	if err := r.db.First(&fossil, fossilID).Error; err != nil {
		return errors.New("данные останки не существует")
	}

	// Проверяем, что останки принадлежат указанному пользователю
	if fossil.UserID != userID {
		return errors.New("пользователь не является создателем этого останка")
	}

	// Проверяем, что обновляем только поле Species
	if updatedFossil.Species != "" {
		// Обновляем только поле Species из JSON-запроса
		if err := r.db.Model(&ds.Fossil{}).Where("id_fossil = ?", fossilID).Update("species", updatedFossil.Species).Error; err != nil {
			return errors.New("ошибка обновления вида")
		}
	} else {
		return errors.New("можно обновлять только вид")
	}

	return nil
}

func (r *Repository) UpdateFossilStatusForUser(fossilID int, userID uint) error {
	// Проверяем, существует ли указанные останки в базе данных
	var fossil ds.Fossil
	if err := r.db.First(&fossil, fossilID).Error; err != nil {
		return errors.New("данные останки не существует")
	}

	// Проверяем, что пользователь имеет право на изменение статуса этого останка
	if fossil.UserID != userID {
		return errors.New("пользователь не является создателем этого останка")
	}

	// Проверяем, что текущий статус останков - "черновик"
	if fossil.FossilStatus == ds.FOSSIL_STATUS_DRAFT {
		// Обновляем статус останков на "в работе"
		fossil.FossilStatus = ds.FOSSIL_STATUS_WORK

		// Обновляем дату формирования на текущее московское время
		moscowTime, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			return err
		}
		fossil.FormationDate = time.Now().In(moscowTime)
	} else {
		return errors.New("останки должны иметь статус черновик")
	}

	// Обновляем останки в базе данных
	if err := r.db.Save(&fossil).Error; err != nil {
		return errors.New("ошибка обновления статуса")
	}

	return nil
}
