package controllers

import (
	"gingorm/models"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateBook godoc
// @Summary CreateBook endpoint is used by the supervisor role user to create a new book.
// @Description CreateBook endpoint is used by the supervisor role user to create a new book
// @Router /api/v1/auth/books/create [post]
// @Tags book
// @Accept json
// @Produce json
// @Param name formData string true "name of the book"
// @Param category_id formData string true "category_id of the book"
func CreateBook(c *gin.Context) {

	var existingBook models.Book
	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	//var User models.User
	var category models.Category
	// user_email, _ := Rdb.HGet("user", "email").Result()

	// // Check if the current user had admin role.
	// if err := models.DB.Where("email = ? AND user_role_id=2", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Product can only be added by supervisor user"})
	// 	return
	// }
	id, _ := models.Rdb.HGet("user", "ID").Result()

	ID, _ := strconv.Atoi(id)
	roleId, _ := models.Rdb.HGet("user", "RoleID").Result()

	if roleId != "2" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Books can only be added by supervisor"})
		return
	}

	c.Request.ParseForm()

	if c.PostForm("name") == "" {
		ReturnParameterMissingError(c, "name")
		return
	}
	if c.PostForm("category_id") == "" {
		ReturnParameterMissingError(c, "category_id")
		return
	}

	title := template.HTMLEscapeString(c.PostForm("name"))
	category_id := template.HTMLEscapeString(c.PostForm("category_id"))
	price, err := strconv.Atoi(template.HTMLEscapeString(c.PostForm("price")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can only convert string to int",
		})
	}

	//Check if the product already exists.
	err = models.DB.Where("title = ?").First(&existingBook).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product already exists."})
		return
	}

	// Check if the category exists
	err = models.DB.First(&category, category_id).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category does not exists."})
		return
	}

	book := models.Book{
		Title:      title,
		CategoryId: category.ID,
		Price:      price,
		CreatedBy:  ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = models.DB.Create(&book).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":   book.ID,
		"name": book.Title,
	})

}

// UpdateBook godoc
// @Summary UpdateBook endpoint is used by the supervisor role user to update a new book.
// @Description Updatebook endpoint is used by the supervisor role user to update a new book
// @Router /api/v1/auth/books/:id/ [PATCH]
// @Tags book
// @Accept json
// @Produce json
func UpdateBook(c *gin.Context) {
	var existingBook models.Book
	var updateBook models.Book
	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	//var User models.User
	// user_email, _ := Rdb.HGet("user", "email").Result()

	// // Check if the current user had admin role.
	// if err := models.DB.Where("email = ? AND user_role_id=2", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Product can only be updated by supervisor user"})
	// 	return
	// }

	id, _ := models.Rdb.HGet("user", "RoleID").Result()

	if id != "2" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Books can only be updated by supervisor"})
		return
	}

	// Check if the product already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingBook).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product doesnot exists."})
		return
	}

	if err := c.ShouldBindJSON(&updateBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&existingBook).Updates(updateBook)

}

type ReturnedBook struct {
	ID         int    `json:"id,string"`
	Title      string `json:"name"`
	CategoryId int    `json:"category_id"`
}

// GetBook godoc
// @Summary GetBook endpoint is used to get info of a book..
// @Description GetBook endpoint is used to get info of a book.
// @Router /api/v1/auth/books/:id/ [get]
// @Tags book
// @Accept json
// @Produce json
func GetBook(c *gin.Context) {
	var existingBook models.Book

	// Check if the product already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingBook).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product doesnot exists."})
		return
	}

	// GET FROM CACHE FIRST
	c.JSON(http.StatusOK, gin.H{"product": existingBook})
	return
}

// ListAllBook godoc
// @Summary ListAllBook endpoint is used to list all book.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/books/ [get]
// @Tags book
// @Accept json
// @Produce json
func ListAllBook(c *gin.Context) {

	// allProduct := []models.Product{}
	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	var User models.User
	var Book []models.Book
	var existingBook []ReturnedBook
	user_email, _ := models.Rdb.HGet("user", "email").Result()

	if err := models.DB.Where("email = ?", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	models.DB.Model(Book).Find(&existingBook)
	c.JSON(http.StatusOK, existingBook)
	return
}

// DeleteBook godoc
// @Summary DeleteBook endpoint is used to delete a book.
// @Description DeleteBook endpoint is used to delete a book.
// @Router /api/v1/auth/books/delete/:id/ [delete]
// @Tags book
// @Accept json
// @Produce json
func DeleteBook(c *gin.Context) {
	var existingBook models.Book
	// var User models.User
	// user_email, _ := Rdb.HGet("user", "email").Result()

	// if err := models.DB.Where("email = ? AND user_role_id=2", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Product can only be updated by supervisor user"})
	// 	return
	// }
	id, _ := models.Rdb.HGet("user", "RoleID").Result()

	if id != "2" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Books can only be deleted by supervisor"})
		return
	}
	// Check if the product already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingBook).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product doesnot exists."})
		return
	}
	models.DB.Where("id = ?", c.Param("id")).Delete(&existingBook)
	// GET FROM CACHE FIRST
	c.JSON(http.StatusOK, gin.H{"Success": "Book deleted"})
	return
}

type UploadedFile struct {
	Status   bool
	BookId   int
	Filename string
	Path     string
	Err      string
}

func generateFilePath(id string, extension string) string {
	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := uuid.New().String() + extension

	projectFolder, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	localS3Folder := projectFolder + "/locals3/"
	imageFolder := localS3Folder + id + "/"

	if _, err := os.Stat(imageFolder); os.IsNotExist(err) {
		os.Mkdir(imageFolder, os.ModeDir)
	}

	imagePath := imageFolder + newFileName
	return imagePath
}

func SaveToBucket(c *gin.Context, f *multipart.FileHeader, extension string, filename string) UploadedFile {
	/*
		whitelist doctionary for extensions
		golang doesnot support "for i in x" construct like python,
		Iterating the list would be expensive, thus we need to use a struct to prevent for loop.
	*/
	acceptedExtensions := map[string]bool{
		".png":  true,
		".jpg":  true,
		".JPEG": true,
		".PNG":  true,
	}
	id, _ := strconv.Atoi(c.Param("id"))

	if !acceptedExtensions[extension] {
		return UploadedFile{Status: false, BookId: id, Filename: filename, Err: "Invalid Extension"}
	}

	filePath := generateFilePath(c.Param("id"), extension)
	err := c.SaveUploadedFile(f, filePath)

	if err == nil {
		return UploadedFile{
			Status:   true,
			BookId:   id,
			Filename: filename,
			Path:     filePath,
			Err:      "",
		}
	}
	return UploadedFile{Status: false, BookId: id, Filename: filename, Err: ""}
}

// UploadProductImages godoc
// @Summary UploadProductImages endpoint is used to add images to product.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/product/:id/image/upload [post]
// @Tags product
// @Accept json
// @Produce json
func UploadBookImages(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if !IsSupervisor(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Product Image can only be added by supervisor"})
		return
	}

	if !DoesProductExist(id) {
		c.JSON(http.StatusNotFound, "Product does not exist")
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["file"]

	var SuccessfullyUploadedFiles []UploadedFile
	var UnSuccessfullyUploadedFiles []UploadedFile
	var BookImages []models.BookImage

	for _, f := range files {
		//save the file to specific dst
		extension := filepath.Ext(f.Filename)
		uploaded_file := SaveToBucket(c, f, extension, f.Filename)
		if uploaded_file.Status {
			SuccessfullyUploadedFiles = append(SuccessfullyUploadedFiles, uploaded_file)
			BookImages = append(BookImages, models.BookImage{
				URL:       uploaded_file.Path,
				BookId:    uploaded_file.BookId,
				CreatedAt: time.Now(),
			})

		} else {
			UnSuccessfullyUploadedFiles = append(UnSuccessfullyUploadedFiles, uploaded_file)
		}
	}
	models.DB.Create(&BookImages)

	c.JSON(http.StatusOK, gin.H{
		"successful": SuccessfullyUploadedFiles, "unsuccessful": UnSuccessfullyUploadedFiles,
	})

}
