package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/lud0m4n/WebAppDev/internal/model"
)

func (r *Repository) GetPeriods(searchName string, userID uint) (model.PeriodGetResponse, error) {
	searchName = strings.Title(searchName + "%")
	var fossilID uint
	if err := r.db.
		Table("fossils").
		Select("fossils.id_fossil").
		Where("user_id = ? AND status = ?", userID, model.FOSSIL_STATUS_DRAFT).
		Take(&fossilID).Error; err != nil {
		//return nil, errors.New("ошибка нахождения id_fossil черновика")
	}

	var periods []model.Period
	if err := r.db.
		Table("periods").
		Select("periods.id_period, periods.name, periods.description, periods.age, periods.status, periods.photo").
		Where("periods.status = ? AND periods.name LIKE ?", model.PERIOD_STATUS_ACTIVE, searchName).
		Order("id_period").
		Scan(&periods).Error; err != nil {
		return model.PeriodGetResponse{}, errors.New("ошибка нахождения списка периодов")
	}

	// Создаем объект JSON для включения id_fossil и periods
	periodResponse := model.PeriodGetResponse{
		Period:   periods,
		IDFossil: fossilID,
	}

	return periodResponse, nil
}

func (r *Repository) GetPeriodByID(periodID int, userID uint) (model.Period, error) {
	var periods model.Period
	if err := r.db.
		Table("periods").
		Select("periods.id_period, periods.name, periods.description, periods.age, periods.status, periods.photo").
		Where("periods.status = ? AND periods.id_period = ?", model.PERIOD_STATUS_ACTIVE, periodID).
		Scan(&periods).Error; err != nil {
		return model.Period{}, errors.New("ошибка нахождения периода по ID")
	}
	return periods, nil
}

func (r *Repository) GetPeriodsBySpecies(fossilSpecies string) ([]map[string]interface{}, error) {

	var periods []map[string]interface{}
	// Выполнение запроса к базе данных для получения периода с указанными параметрами.
	if err := r.db.
		Table("periods").
		Select("periods.id_period, periods.name, periods.description, periods.age, periods.status, periods.photo").
		Joins("JOIN fossilperiods ON periods.id_period = fossilperiods.period_id").
		Joins("JOIN fossils ON fossilperiods.fossil_id = fossils.id_fossil").
		Where("periods.status = ? AND fossils.species = ?", model.PERIOD_STATUS_ACTIVE, fossilSpecies).
		Scan(&periods).Error; err != nil {
		return nil, errors.New("ошибка нахождения списка периодов по названию ископаемого")
	}

	return periods, nil
}

func (r *Repository) CreatePeriod(userID uint, periods *model.Period) error {
	// Создаем период
	if err := r.db.Create(periods).Error; err != nil {
		return errors.New("ошибка создания периода")
	}

	return nil
}

func (r *Repository) DeletePeriod(periodID int, userID uint) error {
	// Удаление изображения из MinIO
	// err := r.minioClient.RemoveServiceImage(periodID)
	// if err != nil {
	// 	// Обработка ошибки удаления изображения из MinIO, если необходимо
	// 	return err
	// }
	return r.db.Exec("UPDATE periods SET status = ? WHERE id_period = ?", model.PERIOD_STATUS_DELETED, periodID).Error
}

func (r *Repository) UpdatePeriod(periodID uint, userID uint, updatedPeriod *model.Period) error {
	err := r.db.Model(&model.Period{}).Where("id_period = ? AND status = ?", periodID, model.PERIOD_STATUS_ACTIVE).Updates(updatedPeriod).Error
	if err != nil {
		return errors.New("ошибка изменения периода")
	}
	return nil
}

func (r *Repository) AddPeriodToFossil(periodID uint, userID uint, moderatorID uint) error {
	// Проверяем, существует ли указанный период в базе данных
	var periods model.Period
	if err := r.db.First(&periods, periodID).Error; err != nil {
		return errors.New("недопустимый ID для периода")
	}

	// Получаем последнюю заявку со статусом "черновик" для указанного пользователя, если такая существует
	var latestDraftFossil model.Fossil
	if err := r.db.Where("status = ? AND user_id = ?", model.FOSSIL_STATUS_DRAFT, userID).Last(&latestDraftFossil).Error; err != nil {
		// Если нет заявки со статусом "черновик", создаем новую
		currentTime := time.Now().In(time.FixedZone("UTC+3", 3*60*60)) // Часовой пояс Москвы
		latestDraftFossil = model.Fossil{
			Status:       model.FOSSIL_STATUS_DRAFT,
			CreationDate: currentTime,
			UserID:       userID, // Устанавливаем ID пользователя для заявки
			ModeratorID:  moderatorID,
		}
		if err := r.db.Create(&latestDraftFossil).Error; err != nil {
			return errors.New("ошибка создания останков со статусом черновик")
		}
	}

	// Создаем связь между периодом и заявкой в промежуточной таблице
	relation := &model.Fossilperiod{
		PeriodID: periodID,
		FossilID: latestDraftFossil.IDFossil,
	}

	// Начинаем транзакцию
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Создаем связь в таблице delivery_periods
	if err := tx.Create(relation).Error; err != nil {
		tx.Rollback()
		return errors.New("ошибка создания связи между периодом и останками")
	}

	// Фиксируем транзакцию
	tx.Commit()

	return nil
}

func (r *Repository) RemovePeriodFromFossil(periodID uint, userID uint) error {
	// Начинаем транзакцию
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Поиск связи между периодом и ископаемым в базе данных
	var relation model.Fossilperiod

	// Проверяем, принадлежит ли период текущему пользователю и находится ли он в статусе "черновик"
	if err := tx.Joins("JOIN fossils ON fossilperiods.fossil_id = fossils.id_fossil").
		Where("fossilperiods.period_id = ? AND fossils.user_id = ? AND fossils.status = ?", periodID, userID, model.FOSSIL_STATUS_DRAFT).
		First(&relation).Error; err != nil {
		tx.Rollback()
		return errors.New("период не принадлежит пользователю или находится не в статусе черновик")
	}

	// Удаление связи из базы данных
	if err := tx.Delete(&relation).Error; err != nil {
		tx.Rollback()
		return errors.New("ошибка удаления связи между периодом и останками")
	}

	// Фиксируем транзакцию
	tx.Commit()

	return nil
}

func (r *Repository) AddPeriodImage(periodID int, userID uint, imageURL string) error {
	err := r.db.Table("periods").Where("id_period = ?", periodID).Update("photo", imageURL).Error
	if err != nil {
		return errors.New("ошибка обновления url изображения в БД")
	}

	return nil
}
