package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/lud0m4n/WebAppDev/internal/app/ds"
)

func (r *Repository) GetPeriods(searchName string, userID uint) (map[string]interface{}, error) {
	searchName = strings.ToUpper(searchName + "%")
	var fossilID uint
	if err := r.db.
		Table("fossils").
		Select("fossils.id_fossil").
		Where("user_id = ? AND status = ?", userID, ds.FOSSIL_STATUS_DRAFT).
		Take(&fossilID).Error; err != nil {
		//return nil, errors.New("ошибка нахождения id_fossil черновика")
	}

	var periods []ds.Period
	if err := r.db.
		Table("periods").
		Select("periods.id_period, periods.name, periods.description, periods.age, periods.status, periods.photo").
		Where("periods.status = ? AND periods.name LIKE ?", ds.PERIOD_STATUS_ACTIVE, searchName).
		Scan(&periods).Error; err != nil {
		return nil, errors.New("ошибка нахождения списка периодов")
	}

	// Создаем объект JSON для включения id_fossil и periods
	result := make(map[string]interface{})
	result["periods"] = periods
	result["id_fossil"] = fossilID

	return result, nil
}

func (r *Repository) GetPeriodByID(periodID int, userID uint) (map[string]interface{}, error) {
	var periods map[string]interface{}
	if err := r.db.
		Table("periods").
		Select("periods.id_period, periods.name, periods.description, periods.age, periods.status, periods.photo").
		Where("periods.status = ? AND periods.id_period = ?", ds.PERIOD_STATUS_ACTIVE, periodID).
		Scan(&periods).Error; err != nil {
		return nil, errors.New("ошибка нахождения периода по ID")
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
		Where("periods.status = ? AND fossils.species = ?", ds.PERIOD_STATUS_ACTIVE, fossilSpecies).
		Scan(&periods).Error; err != nil {
		return nil, errors.New("ошибка нахождения списка периодов по названию ископаемого")
	}

	return periods, nil
}

func (r *Repository) CreatePeriod(periods *ds.Period) error {
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
	return r.db.Exec("UPDATE periods SET status = ? WHERE id_period = ?", ds.PERIOD_STATUS_DELETED, periodID).Error
}

func (r *Repository) UpdatePeriod(periodID int, updatedPeriod *ds.Period) error {
	err := r.db.Model(&ds.Period{}).Where("id_period = ? AND status = ?", periodID, ds.PERIOD_STATUS_ACTIVE).Updates(updatedPeriod).Error
	if err != nil {
		return errors.New("ошибка изменения периода")
	}
	return nil
}

func (r *Repository) AddPeriodToFossil(periodID uint, userID uint, moderatorID uint) error {
	// Проверяем, существует ли указанный период в базе данных
	var periods ds.Period
	if err := r.db.First(&periods, periodID).Error; err != nil {
		return errors.New("недопустимый ID для периода")
	}

	// Получаем последнюю заявку со статусом "черновик" для указанного пользователя, если такая существует
	var latestDraftFossil ds.Fossil
	if err := r.db.Where("status = ? AND user_id = ?", ds.FOSSIL_STATUS_DRAFT, userID).Last(&latestDraftFossil).Error; err != nil {
		// Если нет заявки со статусом "черновик", создаем новую
		currentTime := time.Now().In(time.FixedZone("UTC+3", 3*60*60)) // Часовой пояс Москвы
		latestDraftFossil = ds.Fossil{
			Status:       ds.FOSSIL_STATUS_DRAFT,
			CreationDate: currentTime,
			UserID:       userID, // Устанавливаем ID пользователя для заявки
			ModeratorID:  moderatorID,
		}
		if err := r.db.Create(&latestDraftFossil).Error; err != nil {
			return errors.New("ошибка создания останков со статусом черновик")
		}
	}

	// Создаем связь между периодом и заявкой в промежуточной таблице
	relation := &ds.Fossilperiod{
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
	var relation ds.Fossilperiod

	// Проверяем, принадлежит ли период текущему пользователю и находится ли он в статусе "черновик"
	if err := tx.Joins("JOIN fossils ON fossilperiods.fossil_id = fossils.id_fossil").
		Where("fossilperiods.period_id = ? AND fossils.user_id = ? AND fossils.status = ?", periodID, userID, ds.FOSSIL_STATUS_DRAFT).
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

func (r *Repository) AddPeriodImage(periodID int, imageBytes []byte, contentType string) error {
	// Удаление существующего изображения (если есть)
	err := r.minioClient.RemoveServiceImage(periodID)
	if err != nil {
		return err
	}

	// Загрузка нового изображения в MinIO
	imageURL, err := r.minioClient.UploadServiceImage(periodID, imageBytes, contentType)
	if err != nil {
		return err
	}

	// Обновление информации об изображении в БД (например, ссылки на MinIO)
	err = r.db.Model(&ds.Period{}).Where("id_period = ?", periodID).Update("photo", imageURL).Error
	if err != nil {
		// Обработка ошибки обновления URL изображения в БД, если необходимо
		return errors.New("ошибка обновления url изображения в БД")
	}

	return nil
}
