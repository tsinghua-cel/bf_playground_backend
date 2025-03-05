package openapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tsinghua-cel/bf_playground_backend/config"
	"github.com/tsinghua-cel/bf_playground_backend/models/apimodels"
	"github.com/tsinghua-cel/bf_playground_backend/models/dbmodel"
	"net/http"
	"os"
	"strconv"
)

const (
	maxPageSize = 30
	maxTotal    = 1000
)

type apiHandler struct {
	conf *config.Config
}

func (api apiHandler) GetProjectList(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	offset := (page - 1) * pageSize
	if offset > maxTotal {
		offset = maxTotal - pageSize
	}

	total, pList := dbmodel.GetProjectListWithPagination(offset, pageSize)
	res := make([]apimodels.HomeProject, 0)
	for _, p := range pList {
		maxSlot := dbmodel.GetStrategyCount(p.ProjectId)
		res = append(res, apimodels.HomeProject{
			ProjectId:     p.ProjectId,
			TotalSlot:     int64(maxSlot),
			TotalStrategy: int64(p.StrategyCount),
			StartTime:     p.CreatedAt.Unix(),
			EndTime:       p.UpdatedAt.Unix(),
			Category:      p.StrategyCategory,
		})
	}

	m := make(map[string]interface{})
	m["data"] = res
	m["total"] = total
	m["code"] = http.StatusOK
	m["message"] = "success"
	api.response(c, http.StatusOK, m)
}

func (api apiHandler) GetTopStrategies(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	offset := (page - 1) * pageSize
	if offset > maxTotal {
		offset = maxTotal - pageSize
	}

	// Get strategies with pagination
	total, list := dbmodel.GetStrategyListByGreatLostRatio(offset, pageSize)
	res := make([]apimodels.HomeStrategy, 0)
	if total > maxTotal {
		total = maxTotal
	}

	for _, s := range list {
		rate1 := strconv.FormatFloat(s.HonestLoseRateAvg*100, 'f', 4, 64)
		rate2 := strconv.FormatFloat(s.AttackerLoseRateAvg*100, 'f', 4, 64)
		ratio := s.HonestLoseRateAvg / s.AttackerLoseRateAvg
		rate_ratio := strconv.FormatFloat(ratio*100, 'f', 4, 64)
		res = append(res, apimodels.HomeStrategy{
			StrategyId:           s.UUID,
			ProjectId:            s.ProjectId,
			HonestLoseRateAvg:    fmt.Sprintf("%s%%", rate1),
			MaliciousLoseRateAvg: fmt.Sprintf("%s%%", rate2),
			Ratio:                fmt.Sprintf("%s%%", rate_ratio),
			StrategyContent:      s.Content,
			Category:             s.Category,
		})
	}

	m := make(map[string]interface{})
	m["data"] = res
	m["total"] = total
	m["code"] = http.StatusOK
	m["message"] = "success"
	api.response(c, http.StatusOK, m)
}

func (api apiHandler) GetProjectDetail(c *gin.Context) {
	// get project id from url.
	project := c.Param("id")
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	p, _ := dbmodel.GetProjectById(project)
	if p == nil {
		api.errResponse(c, fmt.Errorf("project not found"))
		return
	}

	maxSlot := dbmodel.GetMaxSlotNumber(p.ProjectId)

	t1 := make([]apimodels.StrategyWithReorgCount, 0)
	{
		list := dbmodel.GetStrategyListByReorgCount(p.ProjectId, 0, pageSize)
		for _, s := range list {
			t1 = append(t1, apimodels.StrategyWithReorgCount{
				Category:        s.Category,
				StrategyId:      s.UUID,
				ReorgCount:      strconv.FormatInt(int64(s.ReorgCount), 10),
				StrategyContent: s.Content,
			})
		}
	}

	t2 := make([]apimodels.StrategyWithHonestLose, 0)
	{
		list := dbmodel.GetStrategyListByHonestLoseRateAvg(p.ProjectId, 0, pageSize)
		for _, s := range list {
			rate := strconv.FormatFloat(s.HonestLoseRateAvg*100, 'f', 4, 64)
			t2 = append(t2, apimodels.StrategyWithHonestLose{
				Category:          s.Category,
				StrategyId:        s.UUID,
				HonestLoseRateAvg: fmt.Sprintf("%s%%", rate),
				StrategyContent:   s.Content,
			})
		}
	}
	t3 := make([]apimodels.StrategyWithGreatHonestLose, 0)
	{
		list := dbmodel.GetStrategyListByGreatLostRatioInProject(p.ProjectId, 0, pageSize)
		for _, s := range list {
			rate1 := strconv.FormatFloat(s.HonestLoseRateAvg*100, 'f', 4, 64)
			rate2 := strconv.FormatFloat(s.AttackerLoseRateAvg*100, 'f', 4, 64)
			ratio := s.HonestLoseRateAvg / s.AttackerLoseRateAvg
			rate_ratio := strconv.FormatFloat(ratio*100, 'f', 4, 64)
			t3 = append(t3, apimodels.StrategyWithGreatHonestLose{
				Category:             s.Category,
				StrategyId:           s.UUID,
				HonestLoseRateAvg:    fmt.Sprintf("%s%%", rate1),
				MaliciousLoseRateAvg: fmt.Sprintf("%s%%", rate2),
				Ratio:                fmt.Sprintf("%s%%", rate_ratio),
				StrategyContent:      s.Content,
			})
		}
	}

	detail := apimodels.ProjectDetail{
		Stat: apimodels.ProjectStat{
			Category:      p.StrategyCategory,
			ProjectId:     p.ProjectId,
			TotalSlot:     int64(maxSlot),
			TotalStrategy: int64(p.StrategyCount),
			StartTime:     p.CreatedAt.Unix(),
			EndTime:       p.UpdatedAt.Unix(),
		},
		StrategiesWithReorgCount:      t1,
		StrategiesWithHonestLose:      t2,
		StrategiesWithGreatHonestLose: t3,
	}

	// no param.
	res := detail

	m := make(map[string]interface{})
	m["data"] = res
	m["code"] = http.StatusOK
	m["message"] = "success"
	api.response(c, http.StatusOK, m)
}

func (api apiHandler) DownloadProject(c *gin.Context) {
	// get project id from url.
	project := c.Param("id")
	p, _ := dbmodel.GetProjectById(project)
	if p == nil {
		api.errResponse(c, fmt.Errorf("project not found"))
		return
	}

	strategies, err := dbmodel.GetStrategyListCSV(p.ProjectId)
	if err != nil {
		api.errResponse(c, err)
		return
	}

	// Create a temporary file to store the CSV data
	tmpFile, err := os.CreateTemp("", "strategies-*.csv")
	if err != nil {
		api.errResponse(c, err)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Write the CSV data to the temporary file
	if _, err := tmpFile.Write(strategies); err != nil {
		api.errResponse(c, err)
		return
	}
	if err := tmpFile.Close(); err != nil {
		api.errResponse(c, err)
		return
	}

	// Set the appropriate headers and return the file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", project))
	c.Header("Content-Type", "text/csv")
	c.File(tmpFile.Name())
}

func (api apiHandler) GetStrategyListByHonestLose(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	offset := (page - 1) * pageSize
	if offset > maxTotal {
		offset = maxTotal - pageSize
	}

	// get project id from url.
	project := c.Param("id")
	p, _ := dbmodel.GetProjectById(project)
	if p == nil {
		api.errResponse(c, fmt.Errorf("project not found"))
		return
	}

	count := dbmodel.GetStrategyCount(p.ProjectId)
	list := dbmodel.GetStrategyListByHonestLoseRateAvg(p.ProjectId, offset, pageSize)
	res := make([]apimodels.StrategyWithHonestLose, 0)
	for _, s := range list {
		rate := strconv.FormatFloat(s.HonestLoseRateAvg*100, 'f', 4, 64)
		res = append(res, apimodels.StrategyWithHonestLose{
			Category:          s.Category,
			StrategyId:        s.UUID,
			HonestLoseRateAvg: fmt.Sprintf("%s%%", rate),
			StrategyContent:   s.Content,
		})
	}

	m := make(map[string]interface{})
	m["data"] = res
	m["total"] = count
	m["code"] = http.StatusOK
	m["message"] = "success"
	api.response(c, http.StatusOK, m)
}

func (api apiHandler) GetStrategyListByRatio(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	offset := (page - 1) * pageSize
	if offset > maxTotal {
		offset = maxTotal - pageSize
	}

	// get project id from url.
	project := c.Param("id")
	p, _ := dbmodel.GetProjectById(project)
	if p == nil {
		api.errResponse(c, fmt.Errorf("project not found"))
		return
	}

	count := dbmodel.GetStrategyCount(p.ProjectId)
	list := dbmodel.GetStrategyListByGreatLostRatioInProject(p.ProjectId, offset, pageSize)
	res := make([]apimodels.StrategyWithGreatHonestLose, 0)
	for _, s := range list {
		rate1 := strconv.FormatFloat(s.HonestLoseRateAvg*100, 'f', 4, 64)
		rate2 := strconv.FormatFloat(s.AttackerLoseRateAvg*100, 'f', 4, 64)
		ratio := s.HonestLoseRateAvg / s.AttackerLoseRateAvg
		rate_ratio := strconv.FormatFloat(ratio*100, 'f', 4, 64)
		res = append(res, apimodels.StrategyWithGreatHonestLose{
			Category:             s.Category,
			StrategyId:           s.UUID,
			HonestLoseRateAvg:    fmt.Sprintf("%s%%", rate1),
			MaliciousLoseRateAvg: fmt.Sprintf("%s%%", rate2),
			Ratio:                fmt.Sprintf("%s%%", rate_ratio),
			StrategyContent:      s.Content,
		})
	}

	m := make(map[string]interface{})
	m["data"] = res
	m["total"] = count
	m["code"] = http.StatusOK
	m["message"] = "success"
	api.response(c, http.StatusOK, m)
}

func (api apiHandler) GetStrategyListByReorg(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	offset := (page - 1) * pageSize
	if offset > maxTotal {
		offset = maxTotal - pageSize
	}

	// get project id from url.
	project := c.Param("id")
	p, _ := dbmodel.GetProjectById(project)
	if p == nil {
		api.errResponse(c, fmt.Errorf("project not found"))
		return
	}

	count := dbmodel.GetStrategyCount(p.ProjectId)
	list := dbmodel.GetStrategyListByReorgCount(p.ProjectId, offset, pageSize)
	res := make([]apimodels.StrategyWithReorgCount, 0)
	for _, s := range list {
		res = append(res, apimodels.StrategyWithReorgCount{
			Category:        s.Category,
			StrategyId:      s.UUID,
			ReorgCount:      strconv.FormatInt(int64(s.ReorgCount), 10),
			StrategyContent: s.Content,
		})
	}

	m := make(map[string]interface{})
	m["data"] = res
	m["total"] = count
	m["code"] = http.StatusOK
	m["message"] = "success"
	api.response(c, http.StatusOK, m)
}

func (api apiHandler) response(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

func (api apiHandler) errResponse(c *gin.Context, err error) {
	m := make(map[string]interface{})
	m["code"] = http.StatusInternalServerError
	m["message"] = err.Error()
	api.response(c, http.StatusInternalServerError, m)
}
