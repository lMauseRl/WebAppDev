package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/lud0m4n/WebAppDev/internal/model"
)

func (r *Repository) GetFossilForUser(searchSpecies, startFormationDate, endFormationDate, fossilStatus string, userID uint) ([]model.FossilRequest, error) {
	searchSpecies = strings.ToLower(searchSpecies + "%")
	fossilStatus = strings.ToLower(fossilStatus + "%")

	// Построение основного запроса для получения ископаемых.
	query := r.db.Table("fossils").
		Select("fossils.id_fossil, fossils.species, fossils.creation_date, fossils.formation_date, fossils.completion_date, fossils.status, users.full_name").
		Joins("JOIN users ON users.id_user = fossils.user_id").
		Where("fossils.status LIKE ? AND fossils.species LIKE ? AND fossils.user_id = ? AND fossils.status != ?", fossilStatus, searchSpecies, userID, model.FOSSIL_STATUS_DELETED)

	// Добавление условия фильтрации по дате формирования, если она указана.
	if startFormationDate != "" && endFormationDate != "" {
		query = query.Where("fossils.formation_date BETWEEN ? AND ?", startFormationDate, endFormationDate)
	}

	// Выполнение запроса и сканирование результатов в слайс fossil.
	var fossils []model.FossilRequest
	if err := query.Scan(&fossils).Error; err != nil {
		return nil, errors.New("ошибка получения ископаемых")
	}
	return fossils, nil
}

func (r *Repository) GetFossilByIDForUser(fossilID int, userID uint) (model.FossilGetResponse, error) {
	var fossil model.FossilGetResponse
	// Получение информации о ископаемых по fossilID.
	if err := r.db.
		Table("fossils").
		Select("fossils.id_fossil, fossils.species, fossils.creation_date, fossils.formation_date, fossils.completion_date, fossils.status").
		Where("fossils.status != ? AND fossils.id_fossil = ? AND fossils.user_id = ?", model.FOSSIL_STATUS_DELETED, fossilID, userID).
		Scan(&fossil).Error; err != nil {
		return model.FossilGetResponse{}, errors.New("ошибка получения останков по ИД")
	}
	var periods []model.Period
	if err := r.db.
		Table("periods").
		Select("periods.id_period, periods.name, periods.description, periods.age, periods.status, periods.photo").
		Joins("JOIN fossilperiods ON periods.id_period = fossilperiods.period_id").
		Where("periods.id_period = fossilperiods.period_id").
		Scan(&periods).Error; err != nil {
		return model.FossilGetResponse{}, errors.New("ошибка нахождения списка периодов")
	}
	// Получение периодов по указанному fossilID.
	// Добавление информации о периоде в поле "periods" внутри останков.
	fossil.Period = periods
	return fossil, nil
}

func (r *Repository) DeleteFossilForUser(fossilID int, userID uint) error {
	// Проверяем, существует ли указанные останки в базе данных
	var fossil model.Fossil
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
	if err := tx.Where("fossil_id = ?", fossilID).Delete(&model.Fossilperiod{}).Error; err != nil {
		tx.Rollback()
		return errors.New("ошибка удаления связей из таблицы-множества")
	}

	// Обновляем статус останков на "удален" с использованием GORM
	err := r.db.Model(&model.Fossil{}).Where("id_fossil = ?", fossilID).Update("status", model.FOSSIL_STATUS_DELETED).Error
	if err != nil {
		return errors.New("ошибка обновления статуса на удален")
	}
	// Фиксируем транзакцию
	tx.Commit()

	return nil
}

func (r *Repository) UpdateFossilForUser(fossilID int, userID uint, updatedFossil *model.FossilUpdateSpeciesRequest) error {
	// Проверяем, существует ли указанные останки в базе данных
	var fossil model.Fossil
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
		if err := r.db.Model(&model.Fossil{}).Where("id_fossil = ?", fossilID).Update("species", updatedFossil.Species).Error; err != nil {
			return errors.New("ошибка обновления вида")
		}
	} else {
		return errors.New("можно обновлять только вид")
	}

	return nil
}

func (r *Repository) UpdateFossilStatusForUser(fossilID int, userID uint) error {
	// Проверяем, существует ли указанные останки в базе данных
	var fossil model.Fossil
	if err := r.db.First(&fossil, fossilID).Error; err != nil {
		return errors.New("данные останки не существует")
	}

	// Проверяем, что пользователь имеет право на изменение статуса этого останка
	if fossil.UserID != userID {
		return errors.New("пользователь не является создателем этого останка")
	}

	// Проверяем, что текущий статус останков - "черновик"
	if fossil.Status == model.FOSSIL_STATUS_DRAFT {
		// Обновляем статус останков на "в работе"
		fossil.Status = model.FOSSIL_STATUS_WORK

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
