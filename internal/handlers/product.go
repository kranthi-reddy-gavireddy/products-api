package handlers

import (
	"fmt"
	apperrors "products-api/app-errors"
	"products-api/dtos"
	"products-api/internal/models"
	"products-api/internal/services"
	"products-api/logger"
	"products-api/utils/helpers"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return apperrors.BadRequestError(err.Error())
	}
	// Set unique ID
	product.SetID()
	if err := helpers.DataValidator(&product); err != nil {
		return apperrors.ValidationError(err.(validator.ValidationErrors))
	}
	err := h.productService.Create(c.Context(), &product)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	pagnation := dtos.PageNation{}
	err := c.QueryParser(&pagnation)
	if err != nil {
		return apperrors.BadRequestError(err.Error())
	}
	pagnation.ApplyDefaults()
	products, err := h.productService.GetProducts(c.Context(), pagnation.Limit, pagnation.Offset)
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

func (h *ProductHandler) FilterProducts(c *fiber.Ctx) error {
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
	products, err := h.productService.FilterProducts(c.Context(), filter.MinPrice, filter.MaxPrice, filter.Category, filter.Limit, filter.Offset, sortClauses)
	if err != nil {
		return err
	}
	return c.JSON(products)
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
