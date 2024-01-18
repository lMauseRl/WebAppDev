package usecase

import (
	"errors"
	"strings"

	"github.com/lud0m4n/WebAppDev/internal/model"
)

type PeriodUseCase interface {
}

func (uc *UseCase) GetPeriods(searchName string, userID uint) (model.PeriodGetResponse, error) {
	if userID < 0 {
		return model.PeriodGetResponse{}, errors.New("недопустимый ИД пользователя")
	}

	searchName = strings.Title(searchName + "%")

	periods, err := uc.Repository.GetPeriods(searchName, userID)
	if err != nil {
		return model.PeriodGetResponse{}, err
	}

	return periods, nil
}

func (uc *UseCase) GetPeriodByID(periodID int, userID uint) (model.Period, error) {
	if periodID <= 0 {
		return model.Period{}, errors.New("недопустимый ИД периода")
	}
	if userID < 0 {
		return model.Period{}, errors.New("недопустимый ИД пользователя")
	}

	period, err := uc.Repository.GetPeriodByID(periodID, userID)
	if err != nil {
		return model.Period{}, err
	}

	return period, nil
}

func (uc *UseCase) CreatePeriod(userID uint, requestPeriod model.PeriodRequest) error {
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	if requestPeriod.Name == "" {
		return errors.New("название должно быть заполнено")
	}
	if requestPeriod.Description == "" {
		return errors.New("описания должно быть заполнено")
	}
	if requestPeriod.Age == "" {
		return errors.New("временной промежуток должен быть заполнен")
	}

	period := model.Period{
		Name:        requestPeriod.Name,
		Description: requestPeriod.Description,
		Age:         requestPeriod.Age,
	}

	err := uc.Repository.CreatePeriod(userID, &period)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) DeletePeriod(periodID int, userID uint) error {
	if periodID <= 0 {
		return errors.New("недопустимый ИД периода")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	err := uc.Repository.DeletePeriod(periodID, userID)
	if err != nil {
		return err
	}

	err = uc.Repository.RemoveServiceImage(periodID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) UpdatePeriod(periodID, userID uint, requestPeriod model.PeriodRequest) error {
	if periodID <= 0 {
		return errors.New("недопустимый ИД периода")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	period := model.Period{
		Name:        requestPeriod.Name,
		Description: requestPeriod.Description,
		Age:         requestPeriod.Age,
	}

	err := uc.Repository.UpdatePeriod(periodID, userID, &period)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) AddPeriodToFossil(periodID, userID, moderatorID uint) error {
	if periodID <= 0 {
		return errors.New("недопустимый ИД периода")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	if moderatorID <= 0 {
		return errors.New("недопустимый ИД модератора")
	}

	err := uc.Repository.AddPeriodToFossil(periodID, userID, moderatorID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) RemovePeriodFromFossil(periodID, userID uint) error {
	if periodID <= 0 {
		return errors.New("недопустимый ИД периода")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	err := uc.Repository.RemovePeriodFromFossil(periodID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) AddPeriodImage(periodID int, userID uint, imageBytes []byte, ContentType string) error {
	if periodID <= 0 {
		return errors.New("недопустимый ИД периода")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	if imageBytes == nil {
		return errors.New("недопустимый imageBytes изображения")
	}
	if ContentType == "" {
		return errors.New("недопустимый ContentType изображения")
	}

	imageURL, err := uc.Repository.UploadServiceImage(periodID, userID, imageBytes, ContentType)
	if err != nil {
		return err
	}

	err = uc.Repository.AddPeriodImage(periodID, userID, imageURL)
	if err != nil {
		return err
	}

	return nil
}
