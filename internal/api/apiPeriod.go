package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lud0m4n/WebAppDev/internal/app/ds"
	"github.com/lud0m4n/WebAppDev/internal/auth"
)

// методы для таблицы period
func (h *Handler) GetPeriods(c *gin.Context) {
	// Получение экземпляра singleton для аутентификации
	authInstance := auth.GetAuthInstance()
	searchName := c.DefaultQuery("searchName", "")
	periods, err := h.Repo.GetPeriods(searchName, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, periods)
}

func (h *Handler) GetPeriodByID(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	periodID, err := strconv.Atoi(c.Param("id_period"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	period, err := h.Repo.GetPeriodByID(periodID, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, period)
}

func (h *Handler) CreatePeriod(c *gin.Context) {
	authInstance := auth.GetAuthInstance()

	searchName := c.DefaultQuery("searchName", "")
	var period ds.Period
	if err := c.BindJSON(&period); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Repo.CreatePeriod(&period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получаем обновленный список периодов
	periods, err := h.Repo.GetPeriods(searchName, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Период создан успешно", "periods": periods})
}
func (h *Handler) DeletePeriod(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	searchName := c.DefaultQuery("searchName", "")
	periodID, err := strconv.Atoi(c.Param("id_period"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Repo.DeletePeriod(periodID, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получаем обновленный список периодов
	periods, err := h.Repo.GetPeriods(searchName, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Период успешно удален", "periods": periods})
}
func (h *Handler) UpdatePeriod(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	periodID, err := strconv.Atoi(c.Param("id_period"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var udatedPeriodRequest ds.Period
	if err := c.BindJSON(&udatedPeriodRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Попытка обновления периода в репозитории
	err = h.Repo.UpdatePeriod(periodID, &udatedPeriodRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получаем обновленный объект периода (map[string]interface{})
	updatedPeriod, err := h.Repo.GetPeriodByID(periodID, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Период успешно обновлен", "period": updatedPeriod})
}

// м-м
func (h *Handler) AddPeriodToFossil(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	//searchName := c.DefaultQuery("searchName", "")
	// Получаем параметры из URL
	periodID, err := strconv.Atoi(c.Param("id_period"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Попытка обновления связи между периодом и ископаемым в репозитории
	err = h.Repo.AddPeriodToFossil(uint(periodID), authInstance.UserID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получаем обновленный список периодов
	periods, err := h.Repo.GetPeriodByID(periodID, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Период успешно добавлен в ископаемое", "periods": periods})
}

func (h *Handler) RemovePeriodFromFossil(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	searchName := c.DefaultQuery("searchName", "")
	var err error // Объявляем переменную здесь

	// Получаем параметры из URL
	periodID, err := strconv.Atoi(c.Param("id_period"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Попытка удаления связи между периодом и ископаемым в репозитории
	err = h.Repo.RemovePeriodFromFossil(uint(periodID), authInstance.UserID) // Используем объявленную переменную err
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	periods, err := h.Repo.GetPeriods(searchName, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Период успешно удален из ископаемого", "periods": periods})
}
