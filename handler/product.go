package handler

import (
	"strconv"
	"strings"

	"github.com/an-halim/golang-api-product/database"
	"github.com/an-halim/golang-api-product/model"
	"github.com/gofiber/fiber/v2"
)

func CreateProduct(c *fiber.Ctx) error {
	db := database.DB.Db
	product := new(model.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	if err := db.Create(&product).Error; err != nil {
		return c.Status(503).SendString(err.Error())
	}
	return c.Status(201).JSON(fiber.Map{
		"message": "Product successfully created",
		"product": product,
	})
}

func GetProducts(c *fiber.Ctx) error {
	db := database.DB.Db
	var products []model.Product

	var priceFrom string
	var priceTo string
	var Query string

	priceFrom = c.Query("priceFrom")
	priceTo = c.Query("priceTo")

	if priceFrom != "" && priceTo != "" {
		Query = "price BETWEEN " + priceFrom + " AND " + priceTo
	} else if priceFrom !=  "" {
		Query = "price >= " + priceFrom
	} else if priceTo != "" {
		Query = "price <= " + priceTo
	}

	// pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "5"))
	var offset int
	if page == 1 {
		offset = 0
	} else {
		offset = (page - 1) * limit
	}

	// get total data
	var total int64
	db.Model(&model.Product{}).Count(&total)

	// get total pages
	var pages int64
	if total%int64(limit) > 0 {
		pages = (total / int64(limit)) + 1
	} else {
		pages = total / int64(limit)
	}

	// get data
	db.Limit(limit).Offset(offset).Omit("DeletedAt").Where(Query).Find(&products)

	if err := db.Find(&products, Query).Error; err != nil {
		return c.Status(503).SendString(err.Error())
	}
	// validate if products is empty
	if len(products) == 0 {
		return c.Status(404).SendString("No products found")
	}

	return c.JSON(fiber.Map{
		"message": "Successfully get products",
		"data": fiber.Map{
			"data": products,
			"total_data": total,
			"total_page": pages,
			"current_page": page,
			"per_page": limit,
			"next_page": page + 1,
			"prev_page": page - 1,
			"next_page_url": c.BaseURL() + c.Path() + "?page=" + strconv.Itoa(page + 1) + "&limit=" + strconv.Itoa(limit),
		},
	})
}

func GetProduct(c *fiber.Ctx) error {
	db := database.DB.Db
	product := new(model.Product)
	if err := db.Omit("DeletedAt").Where("ID = ?", c.Params("id")).First(&product).Error; err != nil {
		return c.Status(503).SendString(err.Error())
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Product successfully fetched",
		"data": product,
	})
}

func SearchProduct(c *fiber.Ctx) error {
	db := database.DB.Db
	var products []model.Product
	if err := db.Omit("DeletedAt").Where("LOWER(Name) LIKE ?",  "%" + strings.ToLower(c.Query("q")) + "%").Find(&products).Error; err != nil {
		return c.Status(503).JSON(fiber.Map{
			"message": "something error with our server",
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":"Successfully fetch data",
		"data": products,
	})
}


func UpdateProduct(c *fiber.Ctx) error {
	db := database.DB.Db
	product := new(model.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	if err := db.Model(&product).Where("ID = ?", c.Params("id")).Updates(&product).Error; err != nil {
		return c.Status(503).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Product successfully updated",
		"data": product,
	})
}


func DeleteProduct(c *fiber.Ctx) error {
	db := database.DB.Db
	product := new(model.Product)
	if err := db.Where("ID = ?", c.Params("id")).Delete(&product).Error; err != nil {
		return c.Status(503).SendString(err.Error())
	}
	return c.JSON(fiber.Map{
		"message": "Product successfully deleted",
		"data": nil,
	})
}

