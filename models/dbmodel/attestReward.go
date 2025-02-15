package dbmodel

import (
	"github.com/astaxie/beego/orm"
)

type AttestReward struct {
	BaseModel
	Epoch          int64 `orm:"column(epoch)" db:"epoch" json:"epoch" form:"epoch"`                                         // epoch
	ValidatorIndex int   `orm:"column(validator_index)" db:"validator_index" json:"validator_index" form:"validator_index"` // validator index
	HeadAmount     int64 `orm:"column(head_amount)" db:"head_amount" json:"head_amount" form:"head_amount"`                 // Head reward amount
	TargetAmount   int64 `orm:"column(target_amount)" db:"target_amount" json:"target_amount" form:"target_amount"`         // Target reward amount
	SourceAmount   int64 `orm:"column(source_amount)" db:"source_amount" json:"source_amount" form:"source_amount"`         // Source reward amount.
	//Head	Target	Source	Inclusion Delay	Inactivity
}

func (AttestReward) TableName() string {
	return "t_attest_reward"
}

type AttestRewardRepository interface {
	GetListByFilter(filters ...interface{}) []*AttestReward
}

type attestRewardRepositoryImpl struct {
	o       orm.Ormer
	project string
}

func NewAttestRewardRepository(o orm.Ormer, project string) AttestRewardRepository {
	return &attestRewardRepositoryImpl{o, project}
}

func (repo *attestRewardRepositoryImpl) GetListByFilter(filters ...interface{}) []*AttestReward {
	list := make([]*AttestReward, 0)
	query := repo.o.QueryTable(new(AttestReward).TableName())
	query = ProjectFilter(query, repo.project)
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	query.OrderBy("-epoch").All(&list)
	return list
}

func GetRewardListByEpoch(project string, epoch int64) []*AttestReward {
	filters := make([]interface{}, 0)
	filters = append(filters, "epoch", epoch)
	return NewAttestRewardRepository(orm.NewOrm(), project).GetListByFilter(filters...)
}

func GetRewardListByValidatorIndex(project string, index int) []*AttestReward {
	filters := make([]interface{}, 0)
	filters = append(filters, "validator_index", index)
	return NewAttestRewardRepository(orm.NewOrm(), project).GetListByFilter(filters...)
}
