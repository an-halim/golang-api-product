package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/an-halim/golang-api-product/database"
	"github.com/an-halim/golang-api-product/model"
	"github.com/an-halim/golang-api-product/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)



func CreateProduct(c *fiber.Ctx) error {
	tin := time.Now()

	db := database.DB.Db
	product := new(model.Product)
	if err := c.BodyParser(product); err != nil {
		return utils.Failed(c, 400, err.Error(), tin)
	}	

	// validate
	if err := utils.ValidateProduct(&utils.Product{
		Name: product.Name,
		Price: product.Price,
		Stock: product.Stock,
		Image_url: product.Image_url,
	}); err != nil {
		return utils.Failed(c, 400, err.Error(), tin)
	}

	if err := db.Create(&product).Error; err != nil {
		return utils.Failed(c, 500, err.Error(), tin)
	}

	return utils.Success(c, 201, product, tin)
}


func GetProducts(c *fiber.Ctx) error {	
	tin := time.Now()

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

	if strings.ToLower(c.Query("q")) != "" {
		if Query != "" {
			Query += " AND LOWER(Name) LIKE '%" + strings.ToLower(c.Query("q")) + "%'"
		} else {
			Query = "LOWER(Name) LIKE '%" + strings.ToLower(c.Query("q")) + "%'"
		}
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
	db.Limit(limit).Offset(offset).Where(Query).Find(&products)

	if err := db.Find(&products, Query).Error; err != nil {
		utils.Failed(c, 500, err.Error(), tin)
	}
	// validate if products is empty
	if len(products) == 0 {
		return utils.Failed(c, 404, "Products not found", tin)
	}

	finishTime := time.Now()
	duration := strconv.FormatFloat(float64(finishTime.UnixNano() - tin.UnixNano()) / 1000000, 'f', 6, 64) + " ms"

	return c.JSON(fiber.Map{
		"correlationid": uuid.New(),
		"success": true,
		"error": "",
		"tin": tin,
		"tout": time.Now(),
		"data": fiber.Map{
			"list": products,
			"total_items": total,
			"total_pages": pages,
			"page": page,
			"page_size": limit,
			"start": tin,
			"finish": finishTime,
			"duration": duration,
		},
	})
}

func GetProduct(c *fiber.Ctx) error {
	tin := time.Now()

	db := database.DB.Db
	product := new(model.Product)
	if err := db.Omit("DeletedAt").Where("ID = ?", c.Params("id")).First(&product).Error; err != nil {
		err_code_string := strings.Split(err.Error(), ":")
		var err_code uint64

		if err_code_string[0] == "record not found" {
			err_code = 404
		} 
		 if err_code_string[0] == "invalid input syntax for type uuid" {
			err_code = 400
		}

		return utils.Failed(c, int(err_code), err.Error(), tin)
	}
	
	return utils.Success(c, 200, product, tin)
}

func UpdateProduct(c *fiber.Ctx) error {
	tin := time.Now()

	db := database.DB.Db
	product := new(model.Product)
	if err := c.BodyParser(product); err != nil {
		return utils.Failed(c, 400, err.Error(), tin)
	}

	if err := db.Model(&product).Where("ID = ?", c.Params("id")).Updates(&product).Error; err != nil {
		return utils.Failed(c, 500, err.Error(), tin)
	}

	return utils.Success(c, 200, product, tin)
}


func DeleteProduct(c *fiber.Ctx) error {
	tin := time.Now()

	db := database.DB.Db
	product := new(model.Product)
	if err := db.Where("ID = ?", c.Params("id")).Delete(&product).Error; err != nil {
		return utils.Failed(c, 500, err.Error(), tin)
	}
	
	return utils.Success(c, 200, nil, tin)
}

