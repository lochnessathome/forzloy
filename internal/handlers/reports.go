package handlers

import (
	"net/http"

	"billing/internal/domain/reports"

	"github.com/labstack/echo/v4"
)

func (h *Handler) ReportsPurchse(c echo.Context) error {
	reportId := c.Param("report_id")
	userId := ParseJWTSubject(c)

	r := reports.New(h.pgPool, h.mnDatabase)

	paid, negativeBalance, err := r.Purchase(reportId, userId)
	if err != nil && !negativeBalance {
		return err
	}
	if err != nil && negativeBalance {
		return c.JSON(http.StatusPaymentRequired, false)
	}

	return c.JSON(http.StatusOK, paid)
}
