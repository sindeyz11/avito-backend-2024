package consts

const (
	IncorrectRequestBody    string = "Некорректное тело запроса"
	InternalServerError     string = "Неизвестная ошибка сервера"
	StatusForbidden         string = "У вас нет прав для исполнения запроса"
	InsufficientPermissions string = "У вас недостаточно прав для выполнения данного запроса"
	FailedToWriteResponse   string = "Не удалось сформировать ответ на запрос"

	IncorrectLimitOffsetParams string = "Некорректно задан limit или/и offset"
	NoUsernameParamPresent     string = "Не задан username пользователя"
	IncorrectVersion           string = "Версия указана не числом"
	IncorrectTenderId          string = "Id тендера должно быть в формате uuid"
	IncorrectTenderStatus      string = "Статус не указан либо указан некорректный"

	UserNotExists            string = "Пользователь не существует или некорректен для данного запроса"
	TenderNotExists          string = "Тендер не существует"
	TenderOrVersionNotExists string = "Тендер или его версия не существует"
)
