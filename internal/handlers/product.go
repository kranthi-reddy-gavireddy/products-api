package handlers

import (
	"fmt"
	apperrors "products-api/app-errors"
	"products-api/dtos"
	"products-api/internal/services"
	"products-api/logger"
	"products-api/utils/helpers"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IProductHandler interface {
	Create(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Filter(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type ProductHandler struct {
	productService services.IProductService
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {

	var req dtos.ProductRequest

	if err := c.BodyParser(&req); err != nil {
		return apperrors.BadRequestError(err.Error())
	}

	if err := helpers.DataValidator(&req); err != nil {
		return apperrors.ValidationError(err.(validator.ValidationErrors))
	}

	res, err := h.productService.Create(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *ProductHandler) Get(c *fiber.Ctx) error {

	pagnation := dtos.PageNation{}

	err := c.QueryParser(&pagnation)
	if err != nil {
		return apperrors.BadRequestError(err.Error())
	}

	pagnation.ApplyDefaults()

	products, err := h.productService.Get(c.Context(), pagnation.Limit, pagnation.Offset)
	if err != nil {
		return err
	}

	return c.JSON(products)
}

func (h *ProductHandler) Filter(c *fiber.Ctx) error {

	filter := dtos.FilterParams{}

	err := c.QueryParser(&filter)
	if err != nil {
		return apperrors.BadRequestError(err.Error())
	}

	if err := helpers.DataValidator(&filter); err != nil {
		return apperrors.ValidationError(err.(validator.ValidationErrors))
	}

	filter.ApplyDefaults()

	var sortClauses []dtos.SortClause
	if filter.Sort != "" {
		err, sortClauses = sortingParser(filter.Sort)
	}
	if err != nil {
		return err
	}

	products, err := h.productService.Filter(c.Context(), filter.MinPrice, filter.MaxPrice, filter.Category, filter.Limit, filter.Offset, sortClauses)
	if err != nil {
		return err
	}

	return c.JSON(products)
}

func (h *ProductHandler) GetByID(c *fiber.Ctx) error {

	id := c.Params("id")
	if id == "" {
		return apperrors.BadRequestError(`Id Cannot be an empty value`)
	}

	product, err := h.productService.GetByID(c.Context(), id)
	if err != nil {
		return err
	}
	if product == nil {
		return apperrors.NotFoundError(fmt.Sprintf(apperrors.PRODUCT_NOT_FOUND, id))
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {

	id := c.Params("id")
	if id == "" {
		return apperrors.BadRequestError(fmt.Sprintf(apperrors.INVALID_PARAMS, "id"))
	}

	err := h.productService.Delete(c.Context(), id)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Product deleted successfully %s", id)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": msg, "id": id})

}

func sortingParser(sortField string) (error, []dtos.SortClause) {

	var sortClauses []dtos.SortClause

	for _, data := range strings.Split(sortField, ",") {

		logger.Infof("Field is %v", data)

		if sortClause, ok := helpers.SortMapper[data]; ok {
			sortClauses = append(sortClauses, sortClause)
		} else {
			return apperrors.BadRequestError(fmt.Sprintf(apperrors.INVALID_SORT_FIELD, sortField)), nil
		}

	}

	return nil, sortClauses
}

func New(productService services.IProductService) IProductHandler {
	return &ProductHandler{productService: productService}
}
