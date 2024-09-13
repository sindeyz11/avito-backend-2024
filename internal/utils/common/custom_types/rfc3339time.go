package custom_types

import (
	"time"
)

type RFC3339Time time.Time

// MarshalJSON переопределяет поведение сериализации для RFC3339Time
func (t RFC3339Time) MarshalJSON() ([]byte, error) {
	formatted := time.Time(t).Format(time.RFC3339)
	return []byte(`"` + formatted + `"`), nil
}

// ConvertToTime преобразует RFC3339Time в time.Time для работы с базой данных
func (t RFC3339Time) ConvertToTime() time.Time {
	return time.Time(t)
}
