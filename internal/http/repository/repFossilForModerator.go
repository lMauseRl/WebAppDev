package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/lud0m4n/WebAppDev/internal/model"
)

func (r *Repository) GetFossilForModerator(searchSpecies, startFormationDate, endFormationDate, fossilStatus string, moderatorID uint) ([]model.FossilRequest, error) {
	searchSpecies = strings.ToUpper(searchSpecies + "%")
	fossilStatus = strings.ToLower(fossilStatus + "%")

	// Построение основного запроса для получения ископаемых.
	query := r.db.Table("fossils").
		Select("DISTINCT fossils.id_fossil, fossils.species, fossils.creation_date, fossils.formation_date, fossils.completion_date, fossils.status").
		Joins("JOIN fossilperiods ON fossils.id_fossil = fossilperiods.fossil_id").
		Joins("JOIN periods ON periods.id_period = fossilperiods.fossil_id").
		Where("fossils.status LIKE ? AND fossils.species LIKE ? AND fossils.moderator_id = ?", fossilStatus, searchSpecies, moderatorID)
	// Добавление условия фильтрации по дате формирования, если она указана.
	if startFormationDate != "" && endFormationDate != "" {
		query = query.Where("fossil.formation_date BETWEEN ? AND ?", startFormationDate, endFormationDate)
	}

	// Выполнение запроса и сканирование результатов в структуру fossil.
	var fossil []model.FossilRequest
	if err := query.Scan(&fossil).Error; err != nil {
		return nil, errors.New("ошибка получения ископаемых")
	}
	return fossil, nil
}

func (r *Repository) GetFossilByIDForModerator(fossilID int, moderatorID uint) (model.FossilGetResponse, error) {
	var fossil_info model.FossilGetResponse
	var fossil model.FossilRequest
	// Получение информации о останках по fossilID.
	if err := r.db.
		Table("fossils").
		Select("fossils.id_fossil, fossils.species, fossils.creation_date, fossils.formation_date, fossils.completion_date, fossils.status").
		Where("fossils.status != ? AND fossils.id_fossil = ? AND fossils.moderator_id = ?", model.FOSSIL_STATUS_DELETED, fossilID, moderatorID).
		Scan(&fossil).Error; err != nil {
		return model.FossilGetResponse{}, errors.New("ошибка получения останков по ИД")
	}

	// Получение периодов по указанному fossilID.
	var periods []model.Period
	if err := r.db.
		Table("periods").
		Select("periods.id_period, periods.name, periods.description, periods.age, periods.status, periods.photo").
		Joins("JOIN fossilperiods ON periods.id_period = fossilperiods.period_id").
		Where("fossilperiods.fossil_id = ?", fossil.IDFossil).
		Scan(&periods).Error; err != nil {
		return model.FossilGetResponse{}, errors.New("ошибка нахождения списка периодов")
	}
	// Добавление информации о периодах в поле "periods" внутри ископаамых.
	fossil_info.IDFossil = fossil.IDFossil
	fossil_info.CompletionDate = fossil.CompletionDate
	fossil_info.CreationDate = fossil.CreationDate
	fossil_info.FormationDate = fossil.FormationDate
	fossil_info.Species = fossil.Species
	fossil_info.Status = fossil.Status
	fossil_info.Periods = periods
	return fossil_info, nil
}

func (r *Repository) UpdateFossilForModerator(fossilID int, moderatorID uint, updatedFossil *model.FossilUpdateSpeciesRequest) error {
	// Проверяем, существует ли указанные останки в базе данных
	var fossil model.Fossil
	if err := r.db.First(&fossil, fossilID).Error; err != nil {
		return errors.New("данного ископаемого не существует в БД")
	}

	// Проверяем, что ископаемое принадлежит указанному пользователю
	if fossil.ModeratorID != moderatorID {
		return errors.New("текущий модератор не имеет прав изменять вид данного ископаемого")
	}

	// Проверяем, что обновляем только поле Species
	if updatedFossil.Species != "" {
		// Обновляем только поле Species из JSON-запроса
		if err := r.db.Model(&model.Fossil{}).Where("id_fossil = ?", fossilID).Update("species", updatedFossil.Species).Error; err != nil {
			return err
		}
	} else {
		return errors.New("ошибка обновления вида")
	}

	return nil
}
func (r *Repository) UpdateFossilStatusForModerator(fossilID int, moderatorID uint, updateRequest *model.FossilUpdateStatusRequest) error {
	// Проверяем, существует ли указанные останки в базе данных
	var fossil model.Fossil
	if err := r.db.First(&fossil, fossilID).Error; err != nil {
		return errors.New("данных останков не существует в БД")
	}

	// Проверяем, что модератор имеет право на изменение статуса этого ископаемого
	if fossil.ModeratorID != moderatorID {
		return errors.New("текущий модератор не имеет прав на изменение статуса данного ископаемого")
	}

	// Проверяем, что текущий статус ископаемого - "в работе"
	if fossil.Status != model.FOSSIL_STATUS_WORK {
		return errors.New("текущий статус останка еще не в работе")
	}

	// Проверяем, что новый статус является "завершен" или "отклонен"
	if updateRequest.Status != model.FOSSIL_STATUS_COMPLETED && updateRequest.Status != model.FOSSIL_STATUS_REJECTED {
		return errors.New("текущий статус останка не завершен или отклонен")
	}

	// Обновляем только поле Status из JSON-запроса
	fossil.Status = updateRequest.Status

	fossil.CompletionDate = time.Now().In(time.FixedZone("MSK", 3*60*60))

	// Обновляем ископаемые в базе данных
	if err := r.db.Save(&fossil).Error; err != nil {
		return errors.New("ошибка обновления статуса ископаемого в БД")
	}

	return nil
}
