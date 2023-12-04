package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lud0m4n/WebAppDev/internal/app/ds"
	"github.com/lud0m4n/WebAppDev/internal/auth"
)

func (h *Handler) GetFossil(c *gin.Context) {
	// Получение экземпляра singleton для аутентификации
	authInstance := auth.GetAuthInstance()

	searchSpecies := c.DefaultQuery("searchSpecies", "")
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
	fossilStatus := c.DefaultQuery("fossilStatus", "")

	// Выбор соответствующего метода репозитория в зависимости от роли пользователя
	var fossil []ds.Fossilperiod
	var err error
	if authInstance.Role == "moderator" {
		// Получение искпоаемых для модератора
		fossil, err = h.Repo.GetFossilForModerator(searchSpecies, startFormationDate, endFormationDate, fossilStatus, authInstance.UserID)
	} else {
		// Получение искпоаемых для пользователя
		fossil, err = h.Repo.GetFossilForUser(searchSpecies, startFormationDate, endFormationDate, fossilStatus, authInstance.UserID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"fossil": fossil})
}

func (h *Handler) GetFossilByID(c *gin.Context) {
	// Получение экземпляра singleton для аутентификации
	authInstance := auth.GetAuthInstance()

	// Получение идентификатора ископаемого из параметров запроса
	fossilID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД ископаемого"})
		return
	}

	// Получение информации о ископаемом в зависимости от роли пользователя
	var fossil map[string]interface{}
	var repoErr error
	if authInstance.Role == "moderator" {
		// Получение ископаемого для модератора
		fossil, repoErr = h.Repo.GetFossilByIDForModerator(fossilID, authInstance.UserID)
	} else {
		// Получение ископаемого для пользователя
		fossil, repoErr = h.Repo.GetFossilByIDForUser(fossilID, authInstance.UserID)
	}

	if repoErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": repoErr.Error()})
		return
	}
	// Возвращение информации о ископаемом
	c.JSON(http.StatusOK, fossil)
}

func (h *Handler) DeleteFossil(c *gin.Context) {
	// Получение экземпляра singleton для аутентификации
	authInstance := auth.GetAuthInstance()

	// Получение идентификатора ископаемого из параметров запроса
	fossilID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД ископаемого"})
		return
	}
	// Возвращение сообщения об успешном удалении и обновленного списка искпоаемых
	searchSpecies := c.DefaultQuery("searchSpecies", "")
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
	fossilStatus := c.DefaultQuery("fossilStatus", "")

	// Проверка, является ли текущий пользователь пользователем (не модератором)
	if authInstance.Role == "moderator" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "данный запрос недоступен для модератора"})
		return
	}

	// Удаление ископаемого только если оно принадлежит текущему пользователю
	err = h.Repo.DeleteFossilForUser(fossilID, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получаем обновленный список искпоаемых
	fossil, err := h.Repo.GetFossilForUser(searchSpecies, startFormationDate, endFormationDate, fossilStatus, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ископаемое успешно удалено", "fossil": fossil})
}

func (h *Handler) UpdateFossil(c *gin.Context) {
	// Получение идентификатора ископаемого из параметров запроса
	fossilID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД ископаемого"})
		return
	}

	// Привязка JSON-запроса к структуре Fossil
	var updatedFossilRequest ds.Fossil
	if err := c.BindJSON(&updatedFossilRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получение экземпляра singleton для аутентификации
	authInstance := auth.GetAuthInstance()

	// Проверка, является ли пользователь авторизованным и имеет ли права на обновление ископаемого
	var repoErr error
	if authInstance.Role == "moderator" {
		// Обновление ископаемого для модератора
		repoErr = h.Repo.UpdateFossilForModerator(fossilID, authInstance.UserID, &updatedFossilRequest)
		if repoErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": repoErr.Error()})
			return
		}
		// Получение обновленного объекта ископаемого
		updatedFossil, err := h.Repo.GetFossilByIDForUser(fossilID, authInstance.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Ископаемое успешно изменено", "fossil": updatedFossil})
	} else {
		// Обновление ископаемого для пользователя
		repoErr = h.Repo.UpdateFossilForUser(fossilID, authInstance.UserID, &updatedFossilRequest)
		if repoErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": repoErr.Error()})
			return
		}
		updatedFossil, err := h.Repo.GetFossilByIDForUser(fossilID, authInstance.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Доставка успешно изменена", "fossil": updatedFossil})
	}
}

func (h *Handler) UpdateFossilStatusForUser(c *gin.Context) {
	// Получение экземпляра singleton для аутентификации
	authInstance := auth.GetAuthInstance()

	// Получение идентификатора ископаемого из параметров запроса
	fossilID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недоупстимый ИД ископаемого"})
		return
	}

	// Проверка роли пользователя
	if authInstance.Role == "user" {
		// Пользователь может обновлять только свои ископаемого
		err = h.Repo.UpdateFossilStatusForUser(fossilID, authInstance.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Получение обновленного объекта ископаемого
		updatedFossil, err := h.Repo.GetFossilByIDForUser(fossilID, authInstance.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Ископаемое успешно обновлено", "fossil": updatedFossil})
	} else if authInstance.Role == "moderator" {
		// Модератор не имеет права обновлять статус искпоаемых пользователя
		c.JSON(http.StatusUnauthorized, gin.H{"error": "данный запрос доступен только пользователю"})
		return
	}
}

func (h *Handler) UpdateFossilStatusForModerator(c *gin.Context) {
	// Получение экземпляра singleton для аутентификации
	authInstance := auth.GetAuthInstance()

	fossilID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД ископаемого"})
		return
	}

	var updateRequest ds.Fossil
	if err := c.BindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка роли пользователя
	if authInstance.Role == "moderator" {
		// Пользователь может обновлять только свои ископаемого
		err = h.Repo.UpdateFossilStatusForModerator(fossilID, authInstance.UserID, &updateRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Получение обновленного объекта ископаемого
		updatedFossil, err := h.Repo.GetFossilByIDForUser(fossilID, authInstance.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Ископаемое успешно обновлена", "fossil": updatedFossil})
	} else if authInstance.Role == "user" {
		// Модератор не имеет права обновлять статус искпоаемых пользователя
		c.JSON(http.StatusUnauthorized, gin.H{"error": "данный запрос доступен только модератору"})
		return
	}
}
