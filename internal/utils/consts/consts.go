package consts

const (
	Construction string = "Construction"
	Delivery     string = "Delivery"
	Manufacture  string = "Manufacture"

	TenderCreated   string = "Created"
	TenderPublished string = "Published"
	TenderClosed    string = "Closed"

	AuthorTypeOrganization string = "Organization"
	AuthorTypeUser         string = "User"

	BidCreated   string = "Created"
	BidPublished string = "Published"
	BidCanceled  string = "Canceled"

	BidApproved string = "Approved"
	BidRejected string = "Rejected"

	IncorrectRequestBody    string = "Некорректное тело запроса"
	InternalServerError     string = "Неизвестная ошибка сервера"
	InsufficientPermissions string = "У вас недостаточно прав для выполнения данного запроса"
	FailedToWriteResponse   string = "Не удалось сформировать ответ на запрос"

	IncorrectLimitOffsetParams string = "Некорректно задан limit или/и offset"
	NoUsernameParamPresent     string = "Не задан username пользователя"
	IncorrectDecision          string = "Некорректно задано решение, варианты: Approved, Rejected"
	IncorrectVersion           string = "Версия указана не числом"
	IncorrectTenderId          string = "Id тендера должно быть в формате uuid"
	IncorrectBidId             string = "Id предложения должно быть в формате uuid"
	IncorrectStatus            string = "Статус не указан либо указан некорректный"

	UserNotExists            string = "Пользователь не существует или некорректен для данного запроса"
	UserOrOrgNotExists       string = "Пользователь или организация не существует"
	TenderNotExists          string = "Тендер не существует"
	BidNotExists             string = "Предложение не существует"
	TenderOrVersionNotExists string = "Тендер или его версия не существует"
	VersionNotExists         string = "Версия не существует"
)
