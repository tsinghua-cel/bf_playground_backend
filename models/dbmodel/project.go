package dbmodel

import (
	"errors"
	"github.com/astaxie/beego/orm"
)

type Project struct {
	BaseModel
	StrategyCategory string `orm:"column(strategy_category)" db:"strategy_category" json:"strategy_category" form:"strategy_category"` // strategy category
	StrategyCount    int    `orm:"column(strategy_count)" db:"strategy_count" json:"strategy_count" form:"strategy_count"`             // strategy count
	LatestSlot       int64  `orm:"column(latest_slot)" db:"latest_slot" json:"latest_slot" form:"latest_slot"`                         // latest slot
}

func (Project) TableName() string {
	return "project"
}

type ProjectRepository interface {
	GetListByFilter(filters ...interface{}) []*Project
	GetListByFilterWithPagination(offset, limit int, filters ...interface{}) (int64, []*Project)
}

type projectRepositoryImpl struct {
	o orm.Ormer
}

func NewProjectRepository(o orm.Ormer) ProjectRepository {
	return &projectRepositoryImpl{o}
}

func (repo *projectRepositoryImpl) GetListByFilter(filters ...interface{}) []*Project {
	list := make([]*Project, 0)
	query := repo.o.QueryTable(new(Project).TableName())
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	// order by time
	query.OrderBy("-created_at").All(&list)
	return list
}

func (repo *projectRepositoryImpl) GetListByFilterWithPagination(offset, limit int, filters ...interface{}) (int64, []*Project) {
	list := make([]*Project, 0)
	query := repo.o.QueryTable(new(Project).TableName())
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	// get count first.
	count, _ := query.Count()
	// order by time and apply pagination
	query.OrderBy("-created_at").Limit(limit, offset).All(&list)
	return count, list
}

func GetProjectList() []*Project {
	filter := make([]interface{}, 0)
	// strategy_count != 0
	filter = append(filter, "strategy_count__gt", 0)
	return NewProjectRepository(orm.NewOrm()).GetListByFilter(filter...)
}

func GetProjectById(id string) (*Project, error) {
	list := NewProjectRepository(orm.NewOrm()).GetListByFilter("project_id", id)
	if len(list) == 0 {
		return nil, errors.New("project not found")
	}

	return list[0], nil
}
func GetProjectListWithPagination(offset, limit int) (int64, []*Project) {
	filter := make([]interface{}, 0)
	// strategy_count != 0
	//filter = append(filter, "strategy_count__gt", 0)
	return NewProjectRepository(orm.NewOrm()).GetListByFilterWithPagination(offset, limit, filter...)
}
