package app

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lMauseRl/WebAppDev/internal/app/ds"

	_ "github.com/lib/pq"
)

// var services = []Service{
// 	{
// 		Id:              1,
// 		Name:            "Сидерий",
// 		Description:     `Данный период знаменит двумя важными событиями: Кислородной катастрофой и Гуронским оледенением`,
// 		Age:             "От 2,5 до 2,3 млрд лет назад",
// 		LivingOrganisms: "Анаэробные водоросли",
// 		Image:           "/image/s.jpg",
// 	},
// 	{
// 		Id:              2,
// 		Name:            "Рясий",
// 		Description:     `Образуется Бушвельдский комплекс и другие похожие интрузии`,
// 		Age:             "От 2,3 до 2,05 млрд лет назад",
// 		LivingOrganisms: "Первые внутриклеточные организмы",
// 		Image:           "/image/r.png",
// 	},
// 	{
// 		Id:              3,
// 		Name:            "Орозирий",
// 		Description:     `В орозирии Земля испытала два крупнейших из известных астероидных ударов`,
// 		Age:             "От 2,05 до 1,8 млрд лет назад",
// 		LivingOrganisms: "Фотосинтезирующие цианобактерии",
// 		Image:           "/image/o.jpg",
// 	},
// 	{
// 		Id:              4,
// 		Name:            "Статерий",
// 		Description:     `Период характеризуется появлением новых платформ. Формируется суперконтинент Колумбия`,
// 		Age:             "От 1,8 до 1,6 млрд лет назад",
// 		LivingOrganisms: "Наиболее ранние свидетельства присутствия эукариот",
// 		Image:           "/image/stat.png",
// 	},
// 	{
// 		Id:   5,
// 		Name: "Калимий",
// 		Description: `Период характеризуется появлением новых континентальных плит
// 		Распался суперконтинент Колумбия`,
// 		Age:             "От 1,6 до 1,4 млрд лет назад",
// 		LivingOrganisms: "Цианобактерии",
// 		Image:           "/image/kalim.jpg",
// 	},
// 	{
// 		Id:              6,
// 		Name:            "Эктазий",
// 		Description:     `Название период получил из-за продолжавшегося осадконакопления и расширения осадочных чехлов`,
// 		Age:             "От 1,4 до 1,2 млрд лет назад",
// 		LivingOrganisms: "В породах возрастом 1200 миллионов лет с канадского острова Сомерсет были обнаружены ископаемые красные водоросли — древнейшие из известных сохранившихся многоклеточных",
// 		Image:           "/image/ectas.jpg",
// 	},
// 	{
// 		Id:              7,
// 		Name:            "Стений",
// 		Description:     `В стении сложился суперконтинент Родиния.`,
// 		Age:             "От 1,2 до 1 млрд лет назад",
// 		LivingOrganisms: "Наиболее ранние ископаемые остатки эукариот, размножавшихся половым путём",
// 		Image:           "/image/steny.gif",
// 	},
// 	{
// 		Id:              8,
// 		Name:            "Тоний",
// 		Description:     `В этом периоде распался суперконтинент Родиния и началась эволюция животных`,
// 		Age:             "От 1 млрд до 720 млн лет назад",
// 		LivingOrganisms: "Первый представитель царства Животные, скорее всего являвшийся губкой",
// 		Image:           "/image/tony.webp",
// 	},
// 	{
// 		Id:              9,
// 		Name:            "Криогений",
// 		Description:     `Этот период характеризовался значительными, вплоть до экватора, оледенениями Земли`,
// 		Age:             "От 720 до 635 млн лет назад",
// 		LivingOrganisms: "Водоросли, страменопилы, инфузории, динофлагелляты и амёбы",
// 		Image:           "/image/krio.jpg",
// 	},
// 	{
// 		Id:              10,
// 		Name:            "Эдиакарий",
// 		Description:     `Начало периода совпадает с окончанием глобального оледенения, а конец — с началом кембрийского взрыва`,
// 		Age:             "От 635 до 538,8 млн лет назад",
// 		LivingOrganisms: "Сегментированные червеобразные животные",
// 		Image:           "/image/ediacar.jpg",
// 	},
// 	{
// 		Id:              11,
// 		Name:            "Кембрий",
// 		Description:     `Развитие беспозвоночных морских обитателей`,
// 		Age:             "От 538,8 до 485 млн лет назад",
// 		LivingOrganisms: "Беспозвоночные морские обитатели",
// 		Image:           "/image/kembriy.jpg",
// 	},
// 	{
// 		Id:              12,
// 		Name:            "Ордовик",
// 		Description:     `Разнообразная морская жизнь, включая позвоночных`,
// 		Age:             "От 485 до 444 млн лет назад",
// 		LivingOrganisms: "Растения, размножающиеся спорами",
// 		Image:           "/image/ordovik.jpg",
// 	},
// 	{
// 		Id:   13,
// 		Name: "Силурийский период",
// 		Description: `Коралловые рифы
// 		Гигантские скорпионы
// 		Первая челюстная рыба`,
// 		Age:             "От 444 до 419 млн лет назад",
// 		LivingOrganisms: "Первая челюстная рыба",
// 		Image:           "/image/silur.jpg",
// 	},
// 	{
// 		Id:              14,
// 		Name:            "Девонский период",
// 		Description:     `Этот период богат биотическими событиями`,
// 		Age:             "От 419 до 359 млн лет назад",
// 		LivingOrganisms: "Многочисленные рыбы, первые бескрылые насекомые",
// 		Image:           "/image/devon.png",
// 	},
// 	{
// 		Id:              15,
// 		Name:            "Каменноугольный период",
// 		Description:     `Максимальное образование угля в болотистых лесах`,
// 		Age:             "От 359 до 299 млн лет назад",
// 		LivingOrganisms: "Развиваются насекомые, земноводные, рептилии, рыбы",
// 		Image:           "/image/ugol.jpg",
// 	},
// 	{
// 		Id:              16,
// 		Name:            "Пермский период",
// 		Description:     `В результате извержения сибирских траппов вымерло 81 % всех морских и 70 % всех наземных видов организмов`,
// 		Age:             "От 299 до 252 млн лет назад",
// 		LivingOrganisms: "Крупные рептилии, амфибии",
// 		Image:           "/image/perm.jpg",
// 	},
// 	{
// 		Id:              17,
// 		Name:            "Триас",
// 		Description:     `Ранние динозавры, крокодилы, черепахи`,
// 		Age:             "От 252 до 201 млн лет назад",
// 		LivingOrganisms: "Первые млекопитающие",
// 		Image:           "/image/trias.jpg",
// 	},
// 	{
// 		Id:              18,
// 		Name:            "Юрский период",
// 		Description:     `В юрский период достигают расцвета такие группы животных, как динозавры, а также ихтиозавры, птерозавры и плезиозавры`,
// 		Age:             "От 201 до 145 млн лет назад",
// 		LivingOrganisms: "Ранние крупные динозаврыПозднее летающие рептилии (птерозавры), самые ранние известные птицы",
// 		Image:           "/image/jura.jpg",
// 	},
// 	{
// 		Id:   19,
// 		Name: "Меловой период",
// 		Description: `Название происходит от писчего мела, который добывается из осадочных отложений
// 		Доминируют динозавры и другие рептилии`,
// 		Age:             "От 145 до 66 млн лет назад",
// 		LivingOrganisms: "Появляются семеноносные растения",
// 		Image:           "/image/mel.jpg",
// 	},
// 	{
// 		Id:              20,
// 		Name:            "Палеоген",
// 		Description:     `Богатая фауна насекомых, ранние летучие мыши, разнообразные виды млекопитающих и птиц`,
// 		Age:             "От 66 до 23,03 млн лет назад",
// 		LivingOrganisms: "В этом периоде начался бурный расцвет млекопитающих",
// 		Image:           "/image/paleo.jpg",
// 	},
// 	{
// 		Id:   21,
// 		Name: "Неоген",
// 		Description: `Сформировался Панамский перешеек, соединяя Северную и Южную Америки
// 		Формируются Гималаи`,
// 		Age:             "От 23,03 до 2,58 млн лет назад",
// 		LivingOrganisms: "Млекопитающие продолжают развиваться",
// 		Image:           "/image/neo.jpg",
// 	},
// 	{
// 		Id:              22,
// 		Name:            "Четвертичный период",
// 		Description:     `Характеризуется началом Четвертичного оледенения — нынешнего ледникового периода`,
// 		Age:             "От 2,58 млн лет до настоящего времени",
// 		LivingOrganisms: "Появление человека",
// 		Image:           "/image/4.jpg",
// 	},
// }

func (a *Application) StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.Static("/image", "./resources/image")
	r.Static("/css", "./resources/css")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		var periods []ds.Periods
		periods, err := a.repository.GetAllPeriods()
		if err != nil { // если не получилось
			log.Printf("cant get product by id %v", err)
			return
		}
		searchQuery := c.DefaultQuery("period", "")

		if searchQuery == "" {
			c.HTML(http.StatusOK, "main_page.tmpl", gin.H{
				"services": periods,
			})
			return
		}

		var result []ds.Periods

		for _, period := range periods {
			if strings.Contains(strings.ToLower(period.Name), strings.ToLower(searchQuery)) {
				result = append(result, period)
			}
		}

		c.HTML(http.StatusOK, "main_page.tmpl", gin.H{
			"services":    result,
			"search_text": searchQuery,
		})
	})

	r.POST("/delete/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			// Обработка ошибки
			log.Printf("cant get consultation by id %v", err)
			c.Redirect(http.StatusMovedPermanently, "/")
		}
		a.repository.DeletePeriods(id)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	// r.GET("/search", func(c *gin.Context) {

	// 	searchQuery := c.DefaultQuery("fsearch","")

	// 	var result []Service

	// 	for _, service := range services {
	// 		if strings.Contains(strings.ToLower(service.Name), strings.ToLower(searchQuery)) {
	// 			result = append(result, service)
	// 		}
	// 	}

	// 	c.HTML(http.StatusOK, "main_page.tmpl", gin.H {
	// 		"services": result,
	// 		"search_text": searchQuery,
	// 	})
	// })

	r.GET("/period/:id", func(c *gin.Context) {
		var periods *ds.Periods

		id, err := strconv.Atoi(c.Param("id"))
		periods, err = a.repository.GetPeriodsByID(id)
		if err != nil {
			// Обработка ошибки
			log.Printf("cant get period by id %v", err)
			return
		}

		c.HTML(http.StatusOK, "period.tmpl", periods)
	})

	r.Run()

	log.Println("Server down")
}
