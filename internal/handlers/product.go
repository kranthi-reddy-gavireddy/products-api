package handlers

import (
	"fmt"
	apperrors "products-api/app-errors"
	"products-api/internal/models"
	"products-api/internal/services"
	"products-api/utils"

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
	if err := utils.DataValidator(&product); err != nil {
		return apperrors.ValidationError(err.(validator.ValidationErrors))
	}
	err := h.productService.Create(c.Context(), &product)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	pagnation := models.PageNation{}
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
	filter := models.FilterParams{}
	err := c.QueryParser(&filter)
	if err != nil {
		return apperrors.BadRequestError(err.Error())
	}
	if err := utils.DataValidator(&filter); err != nil {
		return apperrors.ValidationError(err.(validator.ValidationErrors))
	}
	filter.ApplyDefaults()
	products, err := h.productService.FilterProducts(c.Context(), filter.MinPrice, filter.MaxPrice, filter.Category, filter.Limit, filter.Offset)
	if err != nil {
		return err
	}
	return c.JSON(products)
}
