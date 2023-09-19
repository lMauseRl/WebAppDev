package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var services = []Service{
	{
		Id:   1,
		Name: "Архей",
		Description: `Образование океанов, атмосферы и континентов
		Появление первых бактерий`,
		Age:             "От 4,0 до 2,5 млрд лет назад",
		LivingOrganisms: "Первые бактерии",
		Image:           "/image/arhey.jpg",
	},
	{
		Id:   2,
		Name: "Протерозой",
		Description: `Накопление кислорода
		Появление многоклеточных организмов`,
		Age:             "От 2,5 до 2,3 млрд лет назад",
		LivingOrganisms: "Первые многоклеточные организмы",
		Image:           "/image/proterozoy.jpg",
	},
	{
		Id:              3,
		Name:            "Кембрий",
		Description:     `Развитие беспозвоночных морских обитателей`,
		Age:             "540 млн лет назад",
		LivingOrganisms: "Первые бактерии",
		Image:           "/image/kembriy.jpg",
	},
	{
		Id:   4,
		Name: "Ордовик",
		Description: `Разнообразная морская жизнь, включая позвоночных
		Растения, размножающиеся спорами`,
		Age:             "500 млн лет назад",
		LivingOrganisms: "Первые бактерии",
		Image:           "/image/ordovik.jpg",
	},
	{
		Id:   5,
		Name: "Силурийский период",
		Description: `Коралловые рифы
		Гигантские скорпионы
		Первая челюстная рыба`,
		Age:             "435 млн лет назад",
		LivingOrganisms: "Первые бактерии",
		Image:           "/image/silur.jpg",
	},
	{
		Id:              6,
		Name:            "Девонский период",
		Description:     `Многочисленные рыбы, первые бескрылые насекомые`,
		Age:             "400 млн лет назад",
		LivingOrganisms: "Первые бактерии",
		Image:           "/image/devon.png",
	},
	{
		Id:   7,
		Name: "Каменноугольный период",
		Description: `Максимальное образование угля в болотистых лесах
		Развиваются насекомые, земноводные, рептилии, рыбы`,
		Age:             "345 млн лет назад",
		LivingOrganisms: "Первые бактерии",
		Image:           "/image/ugol.jpg",
	},
	{
		Id:   8,
		Name: "Пермский период",
		Description: `Крупные рептилии, амфибии
		Большинство видов вымерли`,
		Age:             "280 млн лет назад",
		LivingOrganisms: "Первые бактерии",
		Image:           "/image/perm.jpg",
	},
	{
		Id:   9,
		Name: "Триас",
		Description: `Ранние динозавры, крокодилы, черепахи
		Первые млекопитающие`,
		Age:             "248 млн лет назад",
		LivingOrganisms: "Первые бактерии",
		Image:           "/image/trias.jpg",
	},
	{
		Id:   10,
		Name: "Юрский период",
		Description: `Морские рептилии
		Ранние крупные динозаврыПозднее летающие рептилии (птерозавры), самые ранние известные птицы`,
		Age:             "190 млн лет назад",
		LivingOrganisms: "Первые бактерии",
		Image:           "/image/jura.jpg",
	},
	{
		Id:   11,
		Name: "Меловой период",
		Description: `Доминируют динозавры и другие рептилии
		Появляются семеноносные растения`,
		Age:             "136 млн лет назад",
		LivingOrganisms: "Первые бактерии",
		Image:           "/image/mel.jpg",
	},
	{
		Id:              12,
		Name:            "Палеоген",
		Description:     `Богатая фауна насекомых, ранние летучие мыши, разнообразные виды млекопитающих и птиц`,
		Age:             "65 млн лет назад",
		LivingOrganisms: "Первые бактерии",
		Image:           "/image/paleo.jpg",
	},
	{
		Id:   13,
		Name: "Неоген",
		Description: `Дальнейшее развитие млекопитающих и птиц
		Различные формы человека, включая Homo sapiens`,
		Age:             "25 млн лет назад",
		LivingOrganisms: "Первые бактерии",
		Image:           "/image/neo.jpg",
	},
}

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"services": services,
		})
	})

	r.GET("/service/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			// Обработка ошибки
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		service := services[id-1]
		c.HTML(http.StatusOK, "info.tmpl", service)
	})

	r.GET("/search", func(c *gin.Context) {
		searchQuery := c.DefaultQuery("fsearch", "")
		var result []Service

		for _, service := range services {
			if strings.Contains(strings.ToLower(service.Name), strings.ToLower(searchQuery)) {
				result = append(result, service)
			}
		}

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"services": result,
		})
	})
	r.Static("/image", "./resources/image")
	r.Static("/css", "./resources/css")

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
