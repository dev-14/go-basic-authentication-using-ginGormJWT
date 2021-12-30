package controllers

import (
	"fmt"
	"gingorm/models"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddToCart godoc
// @Summary AddToCart endpoint is used to add the book to the cart.
// @Description AddToCart endpoint is used to add the book to the cart.
// @Router /api/v1/auth/cart/add [post]
// @Tags book
// @Accept json
// @Produce json
// @Param title formData string true "title of the book"
func AddToCart(c *gin.Context) {

	//var existingBook models.Book
	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	var User models.User
	var Book models.Book

	username, _ := models.Rdb.HGet("user", "username").Result()
	if username == "" {
		fmt.Println("Redis empty....checking Database for user...")
		err := FillRedis(c)
		if err != nil {
			c.JSON(404, gin.H{
				"error": "something went wrong with redis",
			})
			return
		}
	}

	username, _ = models.Rdb.HGet("user", "username").Result()
	// ID, _ := strconv.Atoi(id)

	// Check if the current user had admin role.
	// if err := models.DB.Where("email = ? AND user_role_id=3", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Product can only be added to cart by user"})
	// 	return
	// }
	if !IsAuthorized(username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product can only be added to cart by user"})
		return
	}

	c.Request.ParseForm()

	if c.PostForm("title") == "" {
		ReturnParameterMissingError(c, "title")
		return
	}
	// if c.PostForm("category_id") == "" {
	// 	ReturnParameterMissingError(c, "category_id")
	// 	return
	// }

	title := template.HTMLEscapeString(c.PostForm("title"))
	//fmt.Println(name)
	// category_id := template.HTMLEscapeString(c.PostForm("category_id"))
	// price, err := strconv.Atoi(template.HTMLEscapeString(c.PostForm("price")))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "can only convert string to int",
	// 	})
	// }

	//Check if the product already exists.
	err := models.DB.Where("title = ?", title).First(&Book).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book does not exist."})
		return
	}
	fmt.Println(Book)
	cart := models.Cart{
		// BookName: Book.Title,
		// Price:    Book.Price,
		User: User,
		Book: Book,
	}

	err = models.DB.Create(&cart).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"name":    title,
	})

}

type book struct {
	Title string `json:"title"` // these json tags are used in front end to access these variables.
	Price int    `json:"price"` // Exporting variables is also required for this purpose
}

// ViewCart godoc
// @Summary ViewCart endpoint is used to list all book.
// @Description API Endpoint to view the cart items.
// @Router /api/v1/auth/cart/view [get]
// @Tags book
// @Accept json
// @Produce json
func ViewCart(c *gin.Context) {

	// allProduct := []models.Product{}
	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	var User models.User
	//var Book []models.Book
	var Cart []models.Cart
	//var exCart models.Cart

	username, _ := models.Rdb.HGet("user", "username").Result()
	if username == "" {
		fmt.Println("Redis empty....checking Database for user...")
		err := FillRedis(c)
		if err != nil {
			c.JSON(404, gin.H{
				"error": "something went wrong with redis",
			})
			return
		}
	}

	username, _ = models.Rdb.HGet("user", "username").Result()

	if err := models.DB.Where("email = ? OR mobile = ?", username, username).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userid := User.ID

	if err := models.DB.Where("user_id = ?", userid).Find(&Cart).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Cart empty",
		})
		return
	}

	//JOINING TABLE USING RAW QUERY
	// rows, err := models.DB.Raw("SELECT title,price FROM books INNER JOIN carts on books.id=carts.book_id").Rows()
	// if err != nil {
	// 	c.JSON(404, gin.H{
	// 		"error": err,
	// 	})
	// }
	// for rows.Next() {
	// 	rows.Scan(&title, &price)
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"title": title,
	// 		"price": price,
	// 	})
	// }

	//JOINING TABLE USING GORM JOINS
	rows, err := models.DB.Table("books").Select("books.title", "books.price").Joins("inner join carts on carts.book_id = books.id and carts.user_id=? ", userid).Rows()
	if err != nil {
		c.JSON(404, gin.H{
			"error": err,
		})
	}
	// if rows == nil {
	// 	c.JSON(http.StatusNoContent, gin.H{
	// 		"message": "cart is empty",
	// 	})
	// }
	var title string
	var price int
	//var onebook book
	var allbooks []book
	for rows.Next() {
		rows.Scan(&title, &price)
		onebook := book{
			Title: title,
			Price: price,
		}
		fmt.Println(onebook)
		// c.JSON(http.StatusOK, gin.H{
		// 	"title": title,
		// 	"price": price,
		// })
		allbooks = append(allbooks, onebook)
	}
	fmt.Println(allbooks)
	c.JSON(200, allbooks)

	// INDIVIDUAL QUERY SEARCHING FOR EACH CART ELEMENT
	// for _, cart := range Cart {
	// 	var tempBook models.Book
	// 	bookId := cart.BookID
	// 	//fmt.Println(bookId)
	// 	//bookid := append(bookId, cart.BookID)
	// 	if err := models.DB.Where("ID = ?", bookId).Find(&tempBook).Error; err != nil {
	// 		c.JSON(http.StatusFound, gin.H{
	// 			"message": "error searching book",
	// 		})
	// 	}
	// 	Book = append(Book, tempBook)
	// }

	// // //models.DB.Model(Cart).Find(&existingBook)
	// c.JSON(http.StatusOK, Book)

}

// DeleteFromCart godoc
// @Summary DeleteFromCart endpoint is used to delete book from cart.
// @Description DeleteFromCart endpoint is used to delete book from cart.
// @Router /api/v1/auth/cart/delete/:id/ [delete]
// @Tags book
// @Accept json
// @Produce json
func DeleteFromCart(c *gin.Context) {
	// var existingBook models.Book
	// var updateBook models.Book
	var Cart models.Cart
	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	var User models.User

	username, _ := models.Rdb.HGet("user", "username").Result()
	if username == "" {
		fmt.Println("Redis empty....checking Database for user...")
		err := FillRedis(c)
		if err != nil {
			c.JSON(404, gin.H{
				"error": "something went wrong with redis",
			})
			return
		}
	}

	username, _ = models.Rdb.HGet("user", "RoleID").Result()

	if err := models.DB.Where("email = ?", username).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userId := User.ID
	// Check if the product already exists.
	err := models.DB.Where("book_id = ? AND user_id = ?", c.Param("id"), userId).Delete(&Cart).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product not in cart."})
		return
	}
	//err = models.DB.Delete(&Cart).Error
	// models.DB.Where("book_id = ?", c.paaram)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"Success": "Book removed from cart",
		})
	}
	// if err := c.ShouldBindJSON(&updateBook); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	//models.DB.Model(&existingBook).Updates(updateBook)

}
