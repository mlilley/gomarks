package handler
//
//import (
//	"github.com/labstack/echo/v4"
//	"github.com/mlilley/gomarks/mark"
//	"net/http"
//)
//
//type createMarkInputData struct {
//	Title string `json:"title"`
//	URL string `json:"url"`
//}
//
//type updateMarkInputData struct {
//	ID string `json:"id"`
//	Title string `json:"title"`
//	URL string `json:"url"`
//}
//
//func (h *Handler) GetMarks(c echo.Context) error {
//	marks, err := h.MarkRepo.FindAll()
//	if err != nil {
//		return err
//	}
//
//	return c.JSON(http.StatusOK, marks)
//}
//
//func (h *Handler) GetMark(c echo.Context) error {
//	id := c.Param("id")
//
//	m, err := h.MarkRepo.FindByID(id)
//	if err != nil {
//		return err
//	}
//
//	if m == nil {
//		return echo.NewHTTPError(http.StatusNotFound)
//	}
//
//	return c.JSON(http.StatusOK, m)
//}
//
//func (h *Handler) CreateMark(c echo.Context) error {
//	input := new(createMarkInputData)
//	err := c.Bind(input)
//	if err != nil {
//		return err // TODO: probably not an InternalServer error
//	}
//
//	m := mark.Mark{
//		Title: input.Title,
//		URL: input.URL,
//	}
//
//	mm, err := h.MarkRepo.Create(&m)
//	if err != nil {
//		return err // TODO probably not an InternalServer error
//	}
//
//	url := c.Echo().Reverse("GetMark", mm.ID)
//
//	c.Response().Header().Set(echo.HeaderLocation, url)
//	return c.NoContent(http.StatusCreated)
//}
//
//func (h *Handler) UpdateMark(c echo.Context) error {
//	id := c.Param("id")
//
//	input := new(updateMarkInputData)
//	err := c.Bind(input)
//	if err != nil {
//		return err // TODO
//	}
//
//	if id != input.ID {
//		return echo.NewHTTPError(http.StatusBadRequest, "Cannot change id")
//	}
//
//	m := mark.Mark{
//		ID: input.ID,
//		Title: input.Title,
//		URL: input.URL,
//	}
//
//	err = h.MarkRepo.Update(&m)
//	if err != nil {
//		return err // TODO
//	}
//
//	return c.NoContent(http.StatusOK)
//}
//
//func (h *Handler) DeleteMark(c echo.Context) error {
//	id := c.Param("id")
//
//	exists, err := h.MarkRepo.DeleteByID(id)
//	if err != nil {
//		return err
//	}
//
//	if !exists {
//		return echo.NewHTTPError(http.StatusNotFound)
//	}
//
//	return c.NoContent(http.StatusOK)
//}
//
//func (h *Handler) DeleteMarks(c echo.Context) error {
//	err := h.MarkRepo.DeleteAll()
//	if err != nil {
//		return err
//	}
//
//	return c.NoContent(http.StatusOK)
//}
