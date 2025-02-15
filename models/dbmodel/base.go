package dbmodel

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type BaseModel struct {
	ID        int64     `orm:"column(id)" db:"id" json:"id" form:"id"`                                 // uniq id
	ProjectId string    `orm:"column(project_id)" db:"project_id" json:"project_id" form:"project_id"` // project id
	CreatedAt time.Time `orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime);column(updated_at)" json:"updated_at"`
}

func ProjectFilter(query orm.QuerySeter, project string) orm.QuerySeter {
	if project != "" {
		return query.Filter("project_id", project)
	}
	return query
}

func ProjectFilterString(project string) string {
	return fmt.Sprintf("project_id = \"%s\"", project)
}
