package store_test

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"techtask/internal/app/model"
	"techtask/internal/app/store"
	"testing"
	"time"
)

var databaseURL = "sqlserver://kursUser:kursPswd@localhost:1400"
// TestRCurrencyRepository_Create func
func TestRCurrencyRepository_Create(t *testing.T) {
	s, _ := store.TestStore(t, databaseURL)

	u, err := s.RCurrency().Create(&model.RCurrency{
		TITLE:  "NEWTEST",
		CODE:   "AUD",
		VALUE:  float32(123.22),
		A_DATE: time.Now(),
	})
	assert.NoError(t, err)
	assert.NotNil(t, u)
	logrus.Info(fmt.Sprintf("Данные были успешно сохранены %+v\\n", u))
}
// TestRCurrencyRepository_FindByDateAndCode func
func TestRCurrencyRepository_FindByDateAndCode(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("TEST.dbo.R_CURRENCY")
	fmt.Println(time.Now().Format("2006-01-02"))
	u, err := s.RCurrency().FindByDateAndCode(
		time.Now(),
		"AUD")
	assert.NoError(t, err)
	assert.NotNil(t, u)
	logrus.Info(fmt.Sprintf("Данные были успешно обработаны %+v\\n", u))
}
