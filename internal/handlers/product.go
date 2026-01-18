package handlers

import (
	"fmt"
	"log"
	"products-api/internal/models"
	"products-api/internal/services"
	"products-api/internal/utils"

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	// Set unique ID
	product.SetID()
	if err := utils.DataValidator(&product); err != nil {
		log.Printf("Validation error: %v", err.(validator.ValidationErrors))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err := h.productService.Create(c.Context(), &product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	pagnation := struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}{}
	err := c.QueryParser(&pagnation)
	if err != nil {
		log.Printf("Error parsing query params: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}
	products, err := h.productService.GetProducts(c.Context(), pagnation.Limit, pagnation.Offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve products"})
	}
	return c.JSON(products)
}

func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Product ID is required"})
	}

	product, err := h.productService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve product"})
	}
	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}
	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Product ID is required"})
	}

	err := h.productService.Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete product"})
	}
	msg := fmt.Sprintf("Product deleted successfully %s", id)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": msg, "id": id})
}

func (h *ProductHandler) FilterProducts(c *fiber.Ctx) error {
	filter := struct {
		MinPrice float64 `json:"min_price" validate:"gte=0"`
		MaxPrice float64 `json:"max_price" validate:"gte=0"`
		Category string  `json:"category"`
		Limit    int     `json:"limit" validate:"gte=0"`
		Offset   int     `json:"offset" validate:"gte=0"`
	}{}
	err := c.QueryParser(&filter)
	if err != nil {
		log.Printf("Error parsing filter query params: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid filter query parameters"})
	}
	if err := utils.DataValidator(&filter); err != nil {
		log.Printf("Validation error: %v", err.(validator.ValidationErrors))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	products, err := h.productService.FilterProducts(c.Context(), filter.MinPrice, filter.MaxPrice, filter.Category, filter.Limit, filter.Offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to filter products"})
	}
	return c.JSON(products)
}
