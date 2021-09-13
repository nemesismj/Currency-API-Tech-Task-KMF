package store

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"techtask/internal/app/model"
	"time"
)
// RCurrencyRepository struct
type RCurrencyRepository struct {
	store *Store
}
// Create - create entity func
func (r *RCurrencyRepository) Create(u *model.RCurrency) (*model.RCurrency, error) {
	pTime := u.A_DATE.Format("2006-01-02")
	tsql := fmt.Sprintf("INSERT INTO TEST.dbo.R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (N'%s','%s', %f, '%s');",
		u.TITLE, u.CODE, u.VALUE, pTime)
	_, err := r.store.db.Exec(tsql)
	if err != nil {
		logrus.Error(fmt.Sprintf("Произошла ошибка при вызова метода INSERT %s", err.Error()))
		return nil, err
	}
	return u, nil
}
// FindByDateAndCode - get list of entities by date and code func
func (r *RCurrencyRepository) FindByDateAndCode(date time.Time, code string) ([]model.RCurrency, error) {
	var items []model.RCurrency
	var tsql string
	if len(code) > 1 {
		tsql = fmt.Sprintf("SELECT * FROM TEST.dbo.R_CURRENCY where A_DATE = '%s' and CODE LIKE '%s';",
			strings.Split(date.String(), " ")[0], code)
	} else {
		tsql = fmt.Sprintf("SELECT * FROM TEST.dbo.R_CURRENCY where A_DATE = '%s';",
			strings.Split(date.String(), " ")[0])
	}
	if err := r.store.db.Select(&items, tsql); err != nil {
		logrus.Error(fmt.Sprintf("Произошла ошибка при вызова метода SELECT %s", err.Error()))
		return nil, err
	}
	return items, nil
}
