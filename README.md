# Quiz API

## Requisitos minimos
- Go 1.25.6
- Postgres

## Clonar el repositorio

```
git clone https://github.com/Moisessantos45/quiz-api.git
cd quiz-api
```

## Comandos basicos de Go

```
go mod tidy
go run main.go
```

## Variables de entorno (.env)
Crea un archivo .env en la raiz con lo minimo:

```
HOST=localhost
PORT=5432
DBUSER=postgres
PASSWORD=postgres
DBNAME=quiz
HOST_URL_DEV=http://localhost:3000
HOST_URL_PROD=https://tu-dominio.com
HOST_URL_PROD_WWW=https://www.tu-dominio.com
```

## Estructura del proyecto

```
config/
  db/
    connection.go
    initialize.go
    seed.go
internal/
  features/
    user/
      domain.go
      repository.go
      usecase.go
      handler.go
  middleware/
    auth.go
    rate.go
  routes/
    category.go
    user.go
  shared/
    models/
      data.go
      models.go
  utils/
    brcrypt.go
    generate_jwt.go
    getIp.go
    validator/
      validator.go
    validators.go
main.go
logical_flow.md
user_journey.md
```

## Que hace cada carpeta
- config: configuracion de base de datos y migraciones.
- internal/features: cada feature con su dominio, repositorio, casos de uso y handler.
- internal/routes: registro de rutas por feature.
- internal/shared/models: modelos GORM y datos semilla.
- internal/middleware: middlewares globales.
- internal/utils: utilidades comunes.

## Arquitectura actual
Seguimos una separacion por capas dentro de cada feature:
- domain.go: contratos (interfaces) y constructores de entidades o validaciones.
- repository.go: acceso a datos (GORM, queries, joins, transacciones).
- usecase.go: reglas de negocio y coordinacion entre repositorios.
- handler.go: entrada HTTP, validacion basica, respuestas.

## Como pensar al crear metodos
- No crees una feature por cada modelo. La feature representa una responsabilidad.
- Un metodo de Option pertenece a Question si su uso depende de la pregunta (ejemplo: opciones de una pregunta).
- Un metodo de User pertenece a User si es identidad efimera y registro en partida.
- Si un metodo es transversal, ponlo en una feature clara o en shared si es util general.

## Relaciones entre tablas
- Define la relacion en models.go y evita logica de negocio en el modelo.
- Para consultas con relaciones, coloca el join o preload en el repositorio.
- Si una operacion requiere varias tablas, coordina en el usecase.
- No expongas la DB directa desde handlers.

## Contenido minimo por feature
- domain.go: interfaces y constructores (validaciones de entrada).
- repository.go: implementacion Postgres con GORM.
- usecase.go: logica principal y orquestacion.
- handler.go: endpoints HTTP.
- routes: el archivo en internal/routes para registrar endpoints.

## Flujo de trabajo con ramas (Git)

### Escenario A: Creando una nueva rama (inicio)
No necesitas merge, solo el pull para preparar main.

```
git switch main
git pull origin main        # Clave: trae los cambios mas recientes
git switch -c feature/nueva-tarea
git push -u origin feature/nueva-tarea
```

### Escenario B: Trayendo cambios posteriores a tu rama
Usas merge cuando se hayan hecho cambios en main despues de que creaste tu rama.

```
git switch main
git pull origin main        # Trae los cambios nuevos
git switch feature/tu-rama
git merge main              # Aplica esos cambios nuevos a tu rama de trabajo
```

### Guardar cambios temporalmente (stash)
Si tienes cambios sin terminar y necesitas moverte a otra rama:

```
# Guardar cambios
git stash push -m "Mensaje descriptivo para estos cambios"

# Recuperar los cambios guardados
git stash pop
```

### Subir cambios desde tu rama

```
git add .           # Puedes dejar el punto para subir todo, o especificar un archivo
git commit -m "mensaje"
git push
```
