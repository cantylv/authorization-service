package myconstants

type AccessKey string

// Частые переменные
const (
	RequestID = "request_id"
	XRealIP   = "X-Real-IP"
)

// Настройка хэширования с помощью Argon2
const (
	HashTime    = 1
	HashMemory  = 2 * 1024
	HashThreads = 2
	HashKeylen  = 56
	HashLetters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
)

var AllowedStatus = map[string]struct{}{
	"approved": {},
	"rejected": {},
}
