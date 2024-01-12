# APIs Swagger (demo-Shop-t-shirt) 🛍️

Esta es una API simple en Go (Golang) para gestionar productos en una tienda ficticia (demo-Shop-t-shirt). Utiliza MongoDB como base de datos para almacenar la información de los productos.

## Especificaciones del Proyecto

### Tecnologías Utilizadas

- Go (Golang)
- MongoDB
- Gin (Framework web para Go)

### Configuración de MongoDB

Asegúrate de tener una instancia de MongoDB en ejecución y actualiza las credenciales en el código o en variables de entorno según sea necesario.

### Instalación de Dependencias

```bash
go get -u github.com/gin-gonic/gin
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/mongo/options
```

### Ejecución del Proyecto
```bash
go run main.go
```
La aplicación se ejecutará en http://localhost:8080 por defecto.

### Endpoints de la API

## Obtener todos los productos

```bash
curl -X GET http://localhost:8080/productos
```

## Crear un nuevo producto

```bash
curl -X POST http://localhost:8080/productos -d '{
  "nombre": "Camiseta",
  "foto": "url_de_la_foto",
  "descripcion": "Descripción de la camiseta",
  "precio": 19.99,
  "color": "Rojo",
  "talla": "M"
}'
```

## Obtener un producto por ID

```bash
curl -X GET http://localhost:8080/productos/<ID_DEL_PRODUCTO>
```

## Actualizar un producto por ID

```bash
curl -X PUT http://localhost:8080/productos/<ID_DEL_PRODUCTO> -d '{
  "nombre": "Camiseta Actualizada",
  "foto": "url_actualizada",
  "descripcion": "Descripción actualizada",
  "precio": 24.99,
  "color": "Azul",
  "talla": "L"
}'
```

## Actualizar parcialmente un producto por ID

```bash
curl -X PATCH http://localhost:8080/productos/<ID_DEL_PRODUCTO> -d '{
  "precio": 29.99
}'
```

## Eliminar un producto por ID

```bash
curl -X DELETE http://localhost:8080/productos/<ID_DEL_PRODUCTO>
```

Nota: Reemplaza <ID_DEL_PRODUCTO> con el ID real de un producto existente si estás realizando una operación en un producto específico.

¡Disfruta de tu experiencia de compra en nuestra Tienda Online ficticia!

**Firma:** Rodrigo Campos Tapia [@DonTester]

**Sígueme en mis redes sociales:**

[<img src="https://simpleicons.org/icons/instagram.svg" alt="Instagram" width="24"/>](https://www.instagram.com/dontester_/) **Instagram** &nbsp; &nbsp;
[<img src="https://simpleicons.org/icons/twitter.svg" alt="Twitter" width="24"/>](https://twitter.com/DonTester_) **Twitter** &nbsp; &nbsp;
[<img src="https://simpleicons.org/icons/linkedin.svg" alt="LinkedIn" width="24"/>](https://www.linkedin.com/in/rcampostapia) **LinkedIn** &nbsp; &nbsp;
[<img src="https://simpleicons.org/icons/github.svg" alt="GitHub" width="24"/>](https://github.com/rcampos09) **GitHub** &nbsp; &nbsp;
[<img src="https://simpleicons.org/icons/youtube.svg" alt="YouTube" width="24"/>](https://www.youtube.com/@dontester) **YouTube** &nbsp; &nbsp;
[<img src="https://simpleicons.org/icons/medium.svg" alt="Medium" width="24"/>](https://medium.com/@rcampos.tapia) **Medium**
