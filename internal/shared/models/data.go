package models

var Categories = []Category{
	{ID: 1, Name: "Geografía"},
	{ID: 2, Name: "Historia Universal"},
	{ID: 3, Name: "Ciencia y Tecnología"},
	{ID: 4, Name: "Programación (Go & Backend)"},
	{ID: 5, Name: "Cine y Series"},
	{ID: 6, Name: "Música"},
	{ID: 7, Name: "Videojuegos"},
	{ID: 8, Name: "Deportes"},
	{ID: 9, Name: "Arte y Literatura"},
	{ID: 10, Name: "Gastronomía"},
	{ID: 11, Name: "Naturaleza y Animales"},
	{ID: 12, Name: "Matemáticas"},
	{ID: 13, Name: "Mitología"},
	{ID: 14, Name: "Cultura General"},
	{ID: 15, Name: "Idiomas"},
}

// esta parte es un ejemplo solamente para mostrar como se pueden cargar datos iniciales en la base de datos
var Questions = []Question{
	// Pregunta de Programación (Texto)
	{
		ID:         1,
		Statement:  "¿Cuál es el comando de Go para ejecutar las pruebas unitarias?",
		CategoryID: 4,
		MediaType:  "text",
	},
	// Pregunta de Geografía (Imagen)
	{
		ID:         2,
		Statement:  "¿A qué país pertenece esta bandera?",
		CategoryID: 1,
		MediaType:  "image",
		MediaURL:   "https://tu-bucket.s3.amazonaws.com/media/flags/mexico.png",
	},
	// Pregunta de Música (Audio)
	{
		ID:         3,
		Statement:  "Escucha el fragmento. ¿Quién es el compositor de esta obra?",
		CategoryID: 6,
		MediaType:  "audio",
		MediaURL:   "https://tu-bucket.s3.amazonaws.com/media/audio/beethoven_5th.mp3",
	},
	// Pregunta de Videojuegos (Texto)
	{
		ID:         4,
		Statement:  "¿En qué año se lanzó originalmente el juego Pac-Man?",
		CategoryID: 7,
		MediaType:  "text",
	},
	// Pregunta de Naturaleza (Imagen)
	{
		ID:         5,
		Statement:  "¿Cómo se llama científicamente esta especie?",
		CategoryID: 11,
		MediaType:  "image",
		MediaURL:   "https://tu-bucket.s3.amazonaws.com/media/nature/axolotl.jpg",
	},
	// Pregunta de Historia (Texto)
	{
		ID:         6,
		Statement:  "¿En qué año cayó el muro de Berlín?",
		CategoryID: 2,
		MediaType:  "text",
	},
}

var Options = []Option{
	// Opciones para Pregunta 1 (Programación)
	{ID: 1, Content: "go run", IsCorrect: false, QuestionID: 1},
	{ID: 2, Content: "go test", IsCorrect: true, QuestionID: 1},
	{ID: 3, Content: "go build", IsCorrect: false, QuestionID: 1},

	// Opciones para Pregunta 2 (Imagen - Bandera)
	{ID: 4, Content: "Italia", IsCorrect: false, QuestionID: 2},
	{ID: 5, Content: "México", IsCorrect: true, QuestionID: 2},
	{ID: 6, Content: "Irlanda", IsCorrect: false, QuestionID: 2},

	// Opciones para Pregunta 3 (Audio - Música)
	{ID: 7, Content: "Mozart", IsCorrect: false, QuestionID: 3},
	{ID: 8, Content: "Beethoven", IsCorrect: true, QuestionID: 3},
	{ID: 9, Content: "Bach", IsCorrect: false, QuestionID: 3},

	// Opciones para Pregunta 4 (Videojuegos)
	{ID: 10, Content: "1980", IsCorrect: true, QuestionID: 4},
	{ID: 11, Content: "1975", IsCorrect: false, QuestionID: 4},
	{ID: 12, Content: "1985", IsCorrect: false, QuestionID: 4},
}
