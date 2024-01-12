package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "swagger-demo-Shop-t-shirt/docs"
)

var collection *mongo.Collection
var secretKey = "token" // Cambia esto y mantenlo seguro

// Respuesta struct para manejar respuestas con error o éxito
type Respuesta struct {
	Error string      `json:"error,omitempty"`
	Datos interface{} `json:"datos,omitempty"`
}

// Función para generar un nuevo token JWT
func generateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix() // Token válido por 15 min
	return token.SignedString([]byte(secretKey))
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := strings.TrimSpace(c.GetHeader("Authorization"))
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, Respuesta{Error: "Token no proporcionado"})
			c.Abort()
			return
		}

		// Verifica si el encabezado comienza con 'Bearer '
		const bearerPrefix = "Bearer "
		if strings.HasPrefix(tokenString, bearerPrefix) {
			// Elimina el prefijo 'Bearer ' del token
			tokenString = tokenString[len(bearerPrefix):]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			fmt.Println("Error al validar el token:", err)
			c.JSON(http.StatusUnauthorized, Respuesta{Error: "Token inválido"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func init() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	username := "demo-Shop-t-shirt-user"
	password := "PEfVJEOkQR75o4Ey"
	encodedUsername := url.QueryEscape(username)
	encodedPassword := url.QueryEscape(password)
	opts := options.Client().ApplyURI("mongodb+srv://" + encodedUsername + ":" + encodedPassword + "@cluster.wiczv8n.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("demo-Shop-t-shirt").Collection("productos")
}

// Producto struct ...
// @swagger:response Producto
type Producto struct {
	ID               string                 `json:"id" bson:"_id"`
	Nombre           string                 `json:"nombre" bson:"nombre"`
	Foto             string                 `json:"foto" bson:"foto"`
	Descripcion      string                 `json:"descripcion" bson:"descripcion"`
	Precio           float64                `json:"precio" bson:"precio"`
	Color            string                 `json:"color" bson:"color"`
	Talla            string                 `json:"talla" bson:"talla"`
	OtrasPropiedades map[string]interface{} `json:"otras_propiedades" bson:"otras_propiedades"`
}

// @title APIs Swagger-Shop-t-Shirt
// @version 1.0
// @description API para gestionar productos (Shop-t-shirt)
// @description @dontesterlabs
// @host localhost:8080
func main() {
	r := gin.Default()

	// Grupo /v1
	v1 := r.Group("/v1")
	v1.GET("/token", generateTokenHandler)
	v1.Use(authMiddleware()) // Aplica el middleware de autenticación a todas las rutas en /v1

	// Rutas del CRUD dentro del grupo /v1
	v1.GET("/productos", getProductos)
	v1.POST("/productos", createProducto)
	v1.GET("/productos/:id", getProductoByID)
	v1.PUT("/productos/:id", updateProducto)
	v1.PATCH("/productos/:id", partialUpdateProducto)
	v1.DELETE("/productos/:id", deleteProducto)

	// Ruta para la documentación Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}

// Función para generar y devolver un nuevo token
func generateTokenHandler(c *gin.Context) {
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Respuesta{Error: "Error al generar el token"})
		return
	}

	fmt.Println("Token generado:", token)

	c.JSON(http.StatusOK, Respuesta{Datos: gin.H{"token": token}})
}

// @Summary Obtener todos los productos
// @Description Obtener todos los productos
// @Tags productos
// @Accept json
// @Produce json
// @Param Authorization header string true "Token de autorización" default(Bearer )
// @Success 200 {array} Respuesta
// @Router /v1/productos [get]
func getProductos(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	fmt.Println("Token recibido:", tokenString)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, Respuesta{Error: "Error al obtener productos"})
		return
	}
	defer cur.Close(ctx)

	var productos []Producto
	if err := cur.All(ctx, &productos); err != nil {
		c.JSON(http.StatusInternalServerError, Respuesta{Error: "Error al decodificar productos"})
		return
	}

	c.JSON(http.StatusOK, Respuesta{Datos: productos})
}

// @Summary Obtener un producto por ID
// @Description Obtener un producto por ID
// @Tags productos
// @Accept json
// @Produce json
// @Param Authorization header string true "Token de autorización" default(Bearer ")
// @Param id path string true "ID del producto"
// @Success 200 {object} Respuesta
// @Failure 404 {object} Respuesta
// @Router /v1/productos/{id} [get]
func getProductoByID(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var producto Producto
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&producto)
	if err != nil {
		c.JSON(http.StatusNotFound, Respuesta{Error: "Producto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, Respuesta{Datos: producto})
}

// @Summary Crear un nuevo producto
// @Description Crear un nuevo producto
// @Tags productos
// @Accept json
// @Produce json
// @Param Authorization header string true "Token de autorización" default(Bearer ")
// @Param producto body Producto true "Datos del producto a crear"
// @Success 201 {object} Respuesta
// @Failure 400 {object} Respuesta
// @Router /v1/productos [post]
func createProducto(c *gin.Context) {
	var producto Producto

	if err := c.ShouldBindJSON(&producto); err != nil {
		c.JSON(http.StatusBadRequest, Respuesta{Error: err.Error()})
		return
	}

	producto.ID = fmt.Sprintf("%s", primitive.NewObjectID())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, producto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Respuesta{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, Respuesta{Datos: producto})
}

// @Summary Actualizar un producto
// @Description Actualizar un producto por ID
// @Tags productos
// @Accept json
// @Produce json
// @Param Authorization header string true "Token de autorización" default(Bearer ")
// @Param id path string true "ID del producto"
// @Param producto body Producto true "Datos del producto a actualizar"
// @Success 200 {object} Respuesta
// @Failure 400 {object} Respuesta
// @Router /v1/productos/{id} [put]
func updateProducto(c *gin.Context) {
	id := c.Param("id")
	var producto Producto

	if err := c.ShouldBindJSON(&producto); err != nil {
		c.JSON(http.StatusBadRequest, Respuesta{Error: err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.ReplaceOne(ctx, bson.M{"_id": id}, producto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Respuesta{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Respuesta{Datos: producto})
}

// @Summary Actualizar parcialmente un producto
// @Description Actualizar parcialmente un producto por ID
// @Tags productos
// @Accept json
// @Produce json
// @Param Authorization header string true "Token de autorización" default(Bearer ")
// @Param id path string true "ID del producto"
// @Param updates body map[string]interface{} true "Datos del producto a actualizar"
// @Success 200 {object} Respuesta
// @Failure 400 {object} Respuesta
// @Router /v1/productos/{id} [patch]
func partialUpdateProducto(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, Respuesta{Error: err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": updates}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Respuesta{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Respuesta{Datos: gin.H{"message": "Producto actualizado exitosamente"}})
}

// @Summary Eliminar un producto
// @Description Eliminar un producto por ID
// @Tags productos
// @Accept json
// @Produce json
// @Param Authorization header string true "Token de autorización" default(Bearer ")
// @Param id path string true "ID del producto"
// @Success 200 {object} Respuesta
// @Failure 400 {object} Respuesta
// @Router /v1/productos/{id} [delete]
func deleteProducto(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, Respuesta{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Respuesta{Datos: gin.H{"message": "Producto eliminado exitosamente"}})
}
