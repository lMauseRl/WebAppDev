package fossil

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lud0m4n/WebAppDev/internal/model"
	"github.com/lud0m4n/WebAppDev/internal/pkg/middleware"
)

// GetFossil godoc
// @Summary Получение списка ископаемых
// @Description Возвращает список всех не удаленных ископаемых
// @Tags Останки
// @Produce json
// @Param searchSpecies query string false "Название вида" Format(email)
// @Param startFormationDate query string false "Начало даты формирования" Format(email)
// @Param endFormationDate query string false "Конец даты формирования" Format(email)
// @Param fossilStatus query string false "Статус ископаемого" Format(email)
// @Success 200 {object} model.FossilRequest "Список ископаемых"
// @Failure 500 {object} model.FossilRequest "Ошибка сервера"
// @Router /fossil [get]
func (h *Handler) GetFossil(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	searchSpecies := c.DefaultQuery("searchSpecies", "")
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
	fossilStatus := c.DefaultQuery("fossilStatus", "")

	var fossils []model.FossilRequest
	var err error

	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		fossils, err = h.UseCase.GetFossilForModerator(searchSpecies, startFormationDate, endFormationDate, fossilStatus, userID)
	} else {
		fossils, err = h.UseCase.GetFossilForUser(searchSpecies, startFormationDate, endFormationDate, fossilStatus, userID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fossils": fossils})
}

// GetFossilByID godoc
// @Summary Получение ископаемого по идентификатору
// @Description Возвращает информацию об останке по её идентификатору
// @Tags Останки
// @Produce json
// @Param id path int true "Идентификатор ископаемого"
// @Success 200 {object} model.FossilGetResponse "Информация об останке"
// @Failure 400 {object} model.FossilGetResponse "Недопустимый идентификатор ископаемого"
// @Failure 500 {object} model.FossilGetResponse "Ошибка сервера"
// @Router /fossil/{id} [get]
func (h *Handler) GetFossilByID(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	fossilID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД останка"})
		return
	}

	var fossil model.FossilGetResponse

	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		fossil, err = h.UseCase.GetFossilByIDForModerator(int(fossilID), userID)
	} else {
		fossil, err = h.UseCase.GetFossilByIDForUser(int(fossilID), userID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fossil)
}

// DeleteFossil godoc
// @Summary Удаление ископаемого
// @Description Удаляет ископаемое по её идентификатору
// @Tags Останки
// @Produce json
// @Param id path int true "Идентификатор ископаемого"
// @Param searchSpecies query string false "Название вида" Format(email)
// @Param startFormationDate query string false "Начало даты формирования" Format(email)
// @Param endFormationDate query string false "Конец даты формирования" Format(email)
// @Param fossilStatus query string false "Статус ископаемого" Format(email)
// @Success 200 {object} model.FossilRequest "Список периодов"
// @Failure 400 {object} model.FossilRequest "Недопустимый идентификатор ископаемого"
// @Failure 500 {object} model.FossilRequest "Ошибка сервера"
// @Router /fossil/{id}/delete [delete]
func (h *Handler) DeleteFossil(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	searchSpecies := c.DefaultQuery("searchSpecies", "")
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
	fossilStatus := c.DefaultQuery("fossilStatus", "")
	fossilID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД останка"})
		return
	}

	err = h.UseCase.DeleteFossilForUser(int(fossilID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fossils, err := h.UseCase.GetFossilForUser(searchSpecies, startFormationDate, endFormationDate, fossilStatus, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fossils": fossils})
}

// UpdateFossil godoc
// @Summary Обновление вид ископаемого
// @Description Обновляет вид для ископаемого по её идентификатору
// @Tags Останки
// @Produce json
// @Param id path int true "Идентификатор ископаемого"
// @Param species body model.FossilUpdateSpeciesRequest true "Новый вид ископаемого"
// @Success 200 {object} model.FossilGetResponse "Информация об останке"
// @Failure 400 {object} model.FossilGetResponse "Недопустимый идентификатор ископаемого или ошибка чтения JSON объекта"
// @Failure 500 {object} model.FossilGetResponse "Ошибка сервера"
// @Router /fossil/{id}/update [put]
func (h *Handler) UpdateFossil(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	fossilID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД останка"})
		return
	}

	var species model.FossilUpdateSpeciesRequest
	if err := c.BindJSON(&species); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка чтения JSON объекта"})
		return
	}

	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		err = h.UseCase.UpdateFossilForModerator(int(fossilID), userID, species)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fossil, err := h.UseCase.GetFossilByIDForModerator(int(fossilID), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"fossil": fossil})
	} else {
		err = h.UseCase.UpdateFossilForUser(int(fossilID), userID, species)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fossil, err := h.UseCase.GetFossilByIDForUser(int(fossilID), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"fossil": fossil})
	}
}

// UpdateFossilStatusForUser godoc
// @Summary Обновление статуса ископаемого для пользователя
// @Description Обновляет статус ископаемого для пользователя по идентификатору ископаемого
// @Tags Останки
// @Produce json
// @Param id path int true "Идентификатор ископаемого"
// @Success 200 {object} model.FossilGetResponse "Информация об останке"
// @Failure 400 {object} model.FossilGetResponse "Недопустимый идентификатор ископаемого"
// @Failure 500 {object} model.FossilGetResponse "Ошибка сервера"
// @Router /fossil/{id}/user [put]
func (h *Handler) UpdateFossilStatusForUser(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	fossilID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недоупстимый ИД останка"})
		return
	}

	err = h.UseCase.UpdateFossilStatusForUser(int(fossilID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fossil, err := h.UseCase.GetFossilByIDForUser(int(fossilID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fossil": fossil})
}

// UpdateFossilStatusForModerator godoc
// @Summary Обновление статуса ископаемого для модератора
// @Description Обновляет статус ископаемого для модератора по идентификатору ископаемого
// @Tags Останки
// @Produce json
// @Param id path int true "Идентификатор ископаемого"
// @Param fossilStatus body model.FossilUpdateStatusRequest true "Новый статус ископаемого"
// @Success 200 {object} model.FossilGetResponse "Информация об останке"
// @Failure 400 {object} model.FossilGetResponse "Недопустимый идентификатор ископаемого или ошибка чтения JSON объекта"
// @Failure 500 {object} model.FossilGetResponse "Ошибка сервера"
// @Router /fossil/{id}/status [put]
func (h *Handler) UpdateFossilStatusForModerator(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	fossilID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД останка"})
		return
	}

	var fossilStatus model.FossilUpdateStatusRequest
	if err := c.BindJSON(&fossilStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		err = h.UseCase.UpdateFossilStatusForModerator(int(fossilID), userID, fossilStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fossil, err := h.UseCase.GetFossilByIDForModerator(int(fossilID), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"fossil": fossil})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "данный запрос доступен только модератору"})
		return
	}
}
