1. Descubrimiento y Acceso (Lobby)
El usuario entra a la aplicación y se encuentra con la pantalla de bienvenida.

Frontend: Realiza un GET /categories para mostrar las opciones disponibles.

Interacción: El usuario selecciona una categoría y, opcionalmente, ingresa una Clave (Key) en un input para unirse a una sala específica.

Backend: Si la clave existe en la tabla Games, el backend valida que la sala esté activa y pertenezca a esa categoría.

2. Identificación Efímera (Registro)
Como mencionaste que no hay sesiones persistentes, cada partida requiere una identidad nueva.

Interacción: Se solicita al usuario su Nickname (y nombre si es necesario).

Backend: 1.  Crea un registro en la tabla Users.
2.  Crea un registro en la tabla Departure vinculando el UserID con el GameID de la sala.
3.  Inicializa el Score en 0.

3. El Ciclo de Juego (Gameplay Loop)
Esta es la parte más dinámica, donde se sirven los retos multimedia.

Backend: El sistema busca las preguntas asociadas a la CategoryID.

Optimización: El backend baraja (shuffle) las Options antes de enviarlas para que el orden sea aleatorio en cada dispositivo.

Frontend: Muestra la pregunta.

Si es image: Renderiza el componente de imagen.

Si es audio: Muestra el reproductor y permite al usuario escucharlo antes de habilitar las opciones.

Interacción: El usuario selecciona una respuesta.

4. Procesamiento de Respuesta (Validación)
Cada vez que el usuario responde, se debe dejar constancia inmediata.

Frontend: Envía un POST /responses con el QuestionID y el OptionID elegido.

Backend: 1.  Verifica si la opción es correcta en la tabla Options.
2.  Crea un registro en GameDetails con el resultado (IsCorrect).
3.  Actualización: Si la respuesta fue correcta, actualiza el Score total en la tabla Departure sumando los puntos correspondientes.

5. Finalización y Estadísticas
Una vez que se agotan las preguntas o el tiempo, se cierra la participación.

Backend: Cambia el estado del Game (si es el administrador) o simplemente marca la Departure como finalizada.

Estadísticas: El backend realiza una consulta agregada:

Cuenta los GameDetails donde is_correct = true.

Calcula el porcentaje de precisión.

Leaderboard: Trae los 5 mejores puntajes de la tabla Departure para ese GameID específico, haciendo un JOIN con Users para mostrar los nicknames.

Frontend: Muestra la pantalla de victoria/resultados con gráficas basadas en el desempeño del usuario y su posición frente a los demás competidores de la sala.

Consideraciones de Ingeniería:
Control de Estado: Como mencionaste que es para jugar entre equipos o en salas, te recomendaría que el backend use un Websocket o un sistema de Polling corto para que todos los usuarios en la misma "sala" vean las estadísticas de los demás en tiempo real al finalizar.

Seguridad Básica: Aunque los usuarios sean efímeros, el DepartureID debería viajar en el header (o un token ligero) para evitar que un usuario registre respuestas a nombre de otro.