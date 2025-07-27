package handlers

import (
	"net/http"

	"billing/internal/domain/reports"

	"github.com/labstack/echo/v4"
)

type ReportsLinkAnonymousRequest struct {
	ClientGeneratedId string `json:"client_generated_id" validate:"required"`
}

type ReportsListRequest struct {
	Limit  int64 `json:"limit" validate:"required"`
	Offset int64 `json:"offset"`
}

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

func (h *Handler) ReportsLinkAnonymous(c echo.Context) error {
	req := new(ReportsLinkAnonymousRequest)

	err := c.Bind(req)
	if err != nil {
		return err
	}

	err = c.Validate(req)
	if err != nil {
		return err
	}

	userId := ParseJWTSubject(c)

	r := reports.New(h.pgPool, h.mnDatabase)

	err = r.LinkAnonymous(req.ClientGeneratedId, userId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, true)
}

func (h *Handler) ReportsList(c echo.Context) error {
	req := new(ReportsListRequest)

	err := c.Bind(req)
	if err != nil {
		return err
	}

	err = c.Validate(req)
	if err != nil {
		return err
	}

	userId := ParseJWTSubject(c)

	r := reports.New(h.pgPool, h.mnDatabase)

	page, err := r.List(userId, req.Limit, req.Offset)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, page)
}
