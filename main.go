package main

import (
	"fmt"
	"log"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Println("No se encontró archivo .env")
	}
	
	fmt.Println("=== QUIZ API ===")
	fmt.Println("Proyecto configurado correctamente")
	fmt.Println("Listo para empezar a codificar")
}
