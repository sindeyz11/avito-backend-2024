package consts

const (
	UserNotExistsError  string = "пользователь не существует"
	CannotFindUserError string = "пользователь с заданными параметрами не найден"
	UnknownBDError      string = "неизвестная ошибка при выполнении запроса в бд"
	ElementNotExists    string = "элемент не существует"

	MethodNotAllowed     string = "Некорректный метод запроса"
	IncorrectRequestBody string = "Некорректное тело запроса"
	InternalServerError  string = "Неизвестная ошибка сервера"
	StatusForbidden      string = "У вас нет прав для исполнения запроса"

	FailedToWriteResponse      string = "Не удалось сформировать ответ на запрос"
	IncorrectLimitOffsetParams string = "Некорректно задан limit или/и offset"
	NoUsernameParamPresent     string = "Не задан username пользователя"
	UserNotExists              string = "Пользователь не существует или некорректен для данного запроса"
	TenderNotExists            string = "Тендер не существует"
	IncorrectTenderId          string = "Id тендера должно быть в формате uuid"
	IncorrectTenderStatus      string = "Статус не указан либо указан некорректный"
)
