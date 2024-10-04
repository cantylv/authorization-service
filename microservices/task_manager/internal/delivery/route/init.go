package route

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/manager"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/middlewares"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHttpHandlers(r *mux.Router, logger *zap.Logger) http.Handler {
	initRegexpMapping()
	manager.InitHandlers(r, logger)
	h := middlewares.Init(r, logger)
	return h
}

// карта, содержащая пути запроса и соответствующие агенты
var MapRequest map[string][]*regexp.Regexp

// nameRegexp     = regexp.MustCompile(`^[A-ZА-ЯЁ][a-zA-Zа-яА-ЯёЁ\s-]{1,50}$`)
func initRegexpMapping() {
	apiVersion := `api/v1`
	name := `[a-zA-Zа-яА-ЯёЁ]+`
	MapRequest = map[string][]*regexp.Regexp{
		// microservice 'privelege'
		"privelege": {
			// агенты
			regexp.MustCompile(fmt.Sprintf(`^%s/ping$`, apiVersion)),
			regexp.MustCompile(fmt.Sprintf(`^%s/agents/%s/who_creates/%s$`, apiVersion, name, name)),
			regexp.MustCompile(fmt.Sprintf(`^%s/agents/%s/who_deletes/%s$`, apiVersion, name, name)),
			regexp.MustCompile(fmt.Sprintf(`^%s/agents/who_reads/%s$`, apiVersion, name)),
			// пользователи
			regexp.MustCompile(fmt.Sprintf(`^%s/users$`, apiVersion)),
			regexp.MustCompile(fmt.Sprintf(`^%s/users/%s$`, apiVersion, name)),
			regexp.MustCompile(fmt.Sprintf(`^%s/users/%s/who_deletes/%s$`, apiVersion, name, name)),
			regexp.MustCompile(fmt.Sprintf(`^%s/openid/callback$`, apiVersion)),
			// группы
			regexp.MustCompile(fmt.Sprintf(`^%s/groups/%s/add_user/%s/who_invites/%s$`, apiVersion, name, name, name)),
			regexp.MustCompile(fmt.Sprintf(`^%s/users/%s/groups/who_asks/%s$`, apiVersion, name, name)),
			regexp.MustCompile(fmt.Sprintf(`^%s/groups/%s/kick_user/%s/who_kicks/%s$`, apiVersion, name, name, name)),
			regexp.MustCompile(fmt.Sprintf(`^%s/groups/%s/who_adds/%s$`, apiVersion, name, name)),
			regexp.MustCompile(fmt.Sprintf(`^%s/users/%s/groups/%s/who_change_status/%s$`, apiVersion, name, name)),
			regexp.MustCompile(fmt.Sprintf(`^%s/groups/%s/users/%s/who_change_owner/%s$`, apiVersion, name, name)),
			// привелегии
			regexp.MustCompile(fmt.Sprintf(`^%s/groups/%s/priveleges/new/agents/%s/who_adds/%s$`, apiVersion, name, name, name)),
			regexp.MustCompile(fmt.Sprintf(`^%s/groups/%s/priveleges/delete/agents/%s/who_deletes/%s$`, apiVersion, name, name, name)),
			regexp.MustCompile(fmt.Sprintf(`^%s/groups/%s/priveleges/who_asks/%s$`, apiVersion, name)),
			regexp.MustCompile(fmt.Sprintf(`^%s/users/%s/priveleges/new/agents/%s/who_adds/%s$`, apiVersion, name, name, name)),
			regexp.MustCompile(fmt.Sprintf(`^%s/users/%s/priveleges/delete/agents/%s/who_deletes/%s$`, apiVersion, name, name, name)),
			regexp.MustCompile(fmt.Sprintf(`^%s/users/%s/check_access/agents/%s$`, apiVersion, name, name)),
		},
	}

}
