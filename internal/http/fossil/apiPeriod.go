package fossil

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lud0m4n/WebAppDev/internal/model"
)

// @Summary Получение списка периодов
// @Description Возращает список всех активных периодов
// @Tags Период
// @Produce json
// @Param searchName query string false "Название периода" Format(email)
// @Success 200 {object} model.PeriodGetResponse "Список периодов"
// @Failure 500 {object} model.PeriodGetResponse "Ошибка сервера"
// @Router /period [get]
func (h *Handler) GetPeriods(c *gin.Context) {
	searchName := c.DefaultQuery("searchName", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "30"))
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте пп"})
		return
	}
	userID := ctxUserID.(uint)

	periods, err := h.UseCase.GetPeriods(searchName, userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, periods)
}

// @Summary Получение периода по ID
// @Description Возвращает информацию о периоде по его ID
// @Tags Период
// @Produce json
// @Param id_period path int true "ID периода"
// @Success 200 {object} model.Period "Информация о периоде"
// @Failure 400 {object} model.Period "Некорректный запрос"
// @Failure 500 {object} model.Period "Внутренняя ошибка сервера"
// @Router /period/{id_period} [get]
func (h *Handler) GetPeriodByID(c *gin.Context) {
	periodID, err := strconv.Atoi(c.Param("id_period"))
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД периода"})
		return
	}

	period, err := h.UseCase.GetPeriodByID(int(periodID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, period)
}

// @Summary Создание нового периода
// @Description Создает новый период с предоставленными данными
// @Tags Период
// @Accept json
// @Produce json
// @Param searchName query string false "Имя периода" Format(email)
// @Param period body model.PeriodRequest true "Пользовательский объект в формате JSON"
// @Success 200 {object} model.PeriodGetResponse "Список периодов"
// @Failure 400 {object} model.PeriodGetResponse "Некорректный запрос"
// @Failure 500 {object} model.PeriodGetResponse "Внутренняя ошибка сервера"
// @Router /period/create [post]
func (h *Handler) CreatePeriod(c *gin.Context) {
	searchName := c.DefaultQuery("searchName", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "30"))
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	var period model.PeriodRequest

	if err := c.BindJSON(&period); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось прочитать JSON"})
		return
	}

	err := h.UseCase.CreatePeriod(userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	periods, err := h.UseCase.GetPeriods(searchName, userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, periods)
}

// @Summary Удаление периода
// @Description Удаляет период по его ID
// @Tags Период
// @Produce json
// @Param id_period path int true "ID периода"
// @Param searchName query string false "Имя периода" Format(email)
// @Success 200 {object} model.PeriodGetResponse "Список периодов"
// @Failure 400 {object} model.PeriodGetResponse "Некорректный запрос"
// @Failure 500 {object} model.PeriodGetResponse "Внутренняя ошибка сервера"
// @Router /period/{id_period}/delete [delete]
func (h *Handler) DeletePeriod(c *gin.Context) {
	searchName := c.DefaultQuery("searchName", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "30"))
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	periodID, err := strconv.Atoi(c.Param("id_period"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД периода"})
		return
	}

	err = h.UseCase.DeletePeriod(int(periodID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	periods, err := h.UseCase.GetPeriods(searchName, userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, periods)
}

// @Summary Обновление информации о периоде
// @Description Обновляет информацию о периоде по его ID
// @Tags Период
// @Accept json
// @Produce json
// @Param id_period path int true "ID периода"
// @Success 200 {object} model.Period "Информация о периоде"
// @Failure 400 {object} model.Period "Некорректный запрос"
// @Failure 500 {object} model.Period "Внутренняя ошибка сервера"
// @Router /period/{id_period}/update [put]
func (h *Handler) UpdatePeriod(c *gin.Context) {
	periodID, err := strconv.Atoi(c.Param("id_period"))
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"error": "недопустимый ИД периода"}})
		return
	}

	var period model.PeriodRequest
	if err := c.BindJSON(&period); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось прочитать JSON"})
		return
	}

	err = h.UseCase.UpdatePeriod(uint(periodID), uint(userID), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedPeriod, err := h.UseCase.GetPeriodByID(int(periodID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPeriod)
}

// @Summary Добавление периода к доставке
// @Description Добавляет период к доставке по его ID
// @Tags Период
// @Produce json
// @Param id_period path int true "ID периода"
// @Param searchName query string false "Имя периода" Format(email)
// @Success 200 {object} model.PeriodGetResponse  "Список периодов"
// @Failure 400 {object} model.PeriodGetResponse  "Некорректный запрос"
// @Failure 500 {object} model.PeriodGetResponse  "Внутренняя ошибка сервера"
// @Router /period/{id_period}/fossil [post]
func (h *Handler) AddPeriodToFossil(c *gin.Context) {
	periodID, err := strconv.Atoi(c.Param("id_period"))
	searchName := c.DefaultQuery("searchName", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "30"))
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД периода"})
		return
	}

	err = h.UseCase.AddPeriodToFossil(uint(periodID), uint(userID), 2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	periods, err := h.UseCase.GetPeriods(searchName, uint(userID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, periods)
}

// @Summary Удаление периода из доставки
// @Description Удаляет период из доставки по его ID
// @Tags Период
// @Produce json
// @Param id_period path int true "ID периода"
// @Param searchName query string false "Имя периода" Format(email)
// @Success 200 {object} model.PeriodGetResponse "Список периодов"
// @Failure 400 {object} model.PeriodGetResponse "Некорректный запрос"
// @Failure 500 {object} model.PeriodGetResponse "Внутренняя ошибка сервера"
// @Router /period/{id_period}/fossil/delete [delete]
func (h *Handler) RemovePeriodFromFossil(c *gin.Context) {
	searchName := c.DefaultQuery("searchName", "")
	periodID, err := strconv.Atoi(c.Param("id_period"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "30"))
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД периода"})
		return
	}

	err = h.UseCase.RemovePeriodFromFossil(uint(periodID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	periods, err := h.UseCase.GetPeriods(searchName, uint(userID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, periods)
}

// @Summary Добавление изображения к периоду
// @Description Добавляет изображение к периоду по его ID
// @Tags Период
// @Accept mpfd
// @Produce json
// @Param id_period path int true "ID периода"
// @Param image formData file true "Изображение периода"
// @Success 200 {object} model.Period "Информация о периоде с изображением"
// @Success 200 {object} model.Period
// @Failure 400 {object} model.Period "Некорректный запрос"
// @Failure 500 {object} model.Period "Внутренняя ошибка сервера"
// @Router /period/{id_period}/image [post]
func (h *Handler) AddPeriodImage(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	periodID, err := strconv.Atoi(c.Param("id_period"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД периода"})
		return
	}

	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимое изображение"})
		return
	}

	file, err := image.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось открыть изображение"})
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось прочитать изображение в байтах"})
		return
	}

	contentType := image.Header.Get("Content-Type")

	err = h.UseCase.AddPeriodImage(int(periodID), uint(userID), imageBytes, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	period, err := h.UseCase.GetPeriodByID(int(periodID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, period)
}
