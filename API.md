# Quiz API — Documentación para Frontend

**Base URL:** `http://localhost:4100`  
**Base API:** `http://localhost:4100/api/v1`  
**WebSocket:** `ws://localhost:4100/ws/:key`

---

## Flujo general

```
1. Crear sala   →  POST /api/v1/room          (admin)
2. Unirse       →  POST /api/v1/room/join     (jugadores)
3. Conectar WS  →  ws://localhost:4100/ws/:key  (todos)
4. Iniciar      →  POST /api/v1/room/:id/start (solo admin)
5. Por turno    →  GET  /api/v1/game/:game_id/question  (cada jugador)
6. Responder    →  POST /api/v1/game/answer   (cada jugador)
7. Finalizar    →  POST /api/v1/game/:game_id/finish    (solo admin)
```

---

## Endpoints REST

### Categorías

#### `GET /api/v1/categories`
Obtiene todas las categorías disponibles. Se usa antes de crear una sala para mostrar las opciones al usuario.

**Request:** ninguno

**Response `200`:**
```json
{
  "data": [
    { "id": 1, "name": "GEOGRAFÍA", "created_at": "...", "updated_at": "..." },
    { "id": 2, "name": "HISTORIA UNIVERSAL", "created_at": "...", "updated_at": "..." },
    { "id": 3, "name": "CIENCIA Y TECNOLOGÍA", "created_at": "...", "updated_at": "..." }
  ]
}
```

---

### Sala (Room)

#### `POST /api/v1/room`
Crea una sala nueva. El jugador que llama este endpoint es el **admin** de la partida y el único que puede iniciarla o terminarla.

**Request body:**
```json
{
  "name": "Sala de Moisés",
  "nickname": "Moises",
  "avatar_url": "https://example.com/avatar.png",
  "category_id": 1
}
```

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `name` | string | ✅ | Nombre de la sala |
| `nickname` | string | ✅ | Apodo del jugador |
| `avatar_url` | string | ❌ | URL del avatar |
| `category_id` | number | ✅ | ID de la categoría del juego |

**Response `201`:**
```json
{
  "data": {
    "game": {
      "id": 1,
      "key": "AB3X9Z",
      "name": "Sala de Moisés",
      "status": "waiting",
      "user_id": 1,
      "category_id": 1,
      "created_at": "..."
    },
    "user": {
      "id": 1,
      "nickname": "Moises",
      "avatar_url": "https://example.com/avatar.png",
      "created_at": "..."
    },
    "departure": {
      "id": 1,
      "game_id": 1,
      "user_id": 1,
      "score": 0,
      "hits": 0,
      "total_time": 0,
      "created_at": "..."
    }
  }
}
```

> **Guarda:** `game.key` para compartir con amigos, `game.id` para iniciar/finalizar, `departure.id` para enviar respuestas, `user.id` para identificar al jugador.

---

#### `POST /api/v1/room/join`
Se une a una sala existente usando la key. Crea un nuevo usuario y una nueva partida (departure) para ese jugador.

**Request body:**
```json
{
  "key": "AB3X9Z",
  "nickname": "Carlos",
  "avatar_url": "https://example.com/avatar2.png"
}
```

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `key` | string | ✅ | Key de 6 caracteres de la sala |
| `nickname` | string | ✅ | Apodo del jugador |
| `avatar_url` | string | ❌ | URL del avatar |

**Response `201`:**
```json
{
  "data": {
    "game": {
      "id": 1,
      "key": "AB3X9Z",
      "name": "Sala de Moisés",
      "status": "waiting",
      "user_id": 1,
      "category_id": 1,
      "created_at": "..."
    },
    "user": {
      "id": 2,
      "nickname": "Carlos",
      "avatar_url": "https://example.com/avatar2.png",
      "created_at": "..."
    },
    "departure": {
      "id": 2,
      "game_id": 1,
      "user_id": 2,
      "score": 0,
      "hits": 0,
      "total_time": 0,
      "created_at": "..."
    }
  }
}
```

> Al unirse correctamente, el servidor hace un broadcast al WebSocket de la sala con el evento `player_joined`.

**Error `400` — sala no encontrada o ya iniciada:**
```json
{ "message": "sala no encontrada" }
{ "message": "la partida ya inicio o finalizo" }
```

---

#### `POST /api/v1/room/:id/start`
Inicia la partida. Solo debe llamarlo el admin (quien creó la sala). Cambia el status del juego a `started`.

**URL param:** `:id` → `game.id` retornado al crear la sala

**Request body:** ninguno

**Response `200`:**
```json
{ "message": "game started" }
```

> Al iniciarse, el frontend debe escuchar el evento `game_started` en el WebSocket para que todos los jugadores comiencen a pedir preguntas.

**Error `400`:**
```json
{ "message": "el juego no esta en espera" }
```

---

### Juego (Game)

#### `GET /api/v1/game/:game_id/question`
Obtiene la siguiente pregunta disponible para ese juego. Las preguntas son únicas por partida (no se repiten entre jugadores).

**URL param:** `:game_id` → `game.id`

**Request:** ninguno

**Response `200`:**
```json
{
  "data": {
    "id": 5,
    "text": "¿Cuál es la capital de Francia?",
    "media_type": "text",
    "category_id": 1,
    "created_at": "...",
    "options": [
      { "id": 10, "content": "París", "is_correct": true, "question_id": 5 },
      { "id": 11, "content": "Londres", "is_correct": false, "question_id": 5 },
      { "id": 12, "content": "Madrid", "is_correct": false, "question_id": 5 },
      { "id": 13, "content": "Roma", "is_correct": false, "question_id": 5 }
    ]
  }
}
```

**Error `404` — no hay más preguntas:**
```json
{ "message": "no hay mas preguntas disponibles" }
```

---

#### `POST /api/v1/game/answer`
Envía la respuesta de un jugador. El servidor calcula si es correcta, actualiza score, hits y tiempo total. Hace broadcast del scoreboard actualizado al WebSocket.

**Request body:**
```json
{
  "departure_id": 2,
  "question_id": 5,
  "answer_id": 10,
  "response_time": 8,
  "game_key": "AB3X9Z"
}
```

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `departure_id` | number | ✅ | ID de la partida del jugador (`departure.id`) |
| `question_id` | number | ✅ | ID de la pregunta que se respondió |
| `answer_id` | number | ✅ | ID de la opción seleccionada |
| `response_time` | number | ✅ | Tiempo de respuesta en segundos |
| `game_key` | string | ✅ | Key de la sala (para hacer broadcast WS) |

**Response `200`:**
```json
{
  "data": {
    "detail": {
      "id": 1,
      "departure_id": 2,
      "question_id": 5,
      "answer_id": 10,
      "is_correct": true,
      "response_time": 8
    },
    "departure": {
      "id": 2,
      "game_id": 1,
      "user_id": 2,
      "score": 10,
      "hits": 1,
      "total_time": 8
    }
  }
}
```

> Cada respuesta correcta suma **10 puntos** al score del jugador.  
> Tras registrar la respuesta, el servidor envía un `score_update` por WebSocket a todos en la sala.

---

#### `GET /api/v1/game/:game_id/scoreboard`
Obtiene el marcador actual ordenado por score de mayor a menor.

**URL param:** `:game_id` → `game.id`

**Response `200`:**
```json
{
  "data": [
    {
      "id": 1,
      "game_id": 1,
      "user_id": 1,
      "score": 30,
      "hits": 3,
      "total_time": 22,
      "user": { "id": 1, "nickname": "Moises", "avatar_url": "..." }
    },
    {
      "id": 2,
      "game_id": 1,
      "user_id": 2,
      "score": 20,
      "hits": 2,
      "total_time": 18,
      "user": { "id": 2, "nickname": "Carlos", "avatar_url": "..." }
    }
  ]
}
```

---

#### `POST /api/v1/game/:game_id/finish`
Finaliza la partida. Solo debe llamarlo el admin. Cambia el status a `finished`, hace broadcast del resultado final al WebSocket y retorna el scoreboard.

**URL param:** `:game_id` → `game.id`  
**Query param:** `?key=AB3X9Z` (requerido para el broadcast WS)

**Ejemplo:** `POST /api/v1/game/1/finish?key=AB3X9Z`

**Request body:** ninguno

**Response `200`:**
```json
{
  "message": "game finished",
  "data": [
    {
      "id": 1,
      "score": 30,
      "hits": 3,
      "total_time": 22,
      "user": { "id": 1, "nickname": "Moises", "avatar_url": "..." }
    }
  ]
}
```

---

## WebSocket

### Conexión

**URL:** `ws://localhost:4100/ws/:key`  
**Param:** `:key` → key de 6 caracteres de la sala (`game.key`)

**Ejemplo:**
```
ws://localhost:4100/ws/AB3X9Z
```

Todos los jugadores (incluido el admin) deben conectarse a este canal **inmediatamente después** de crear/unirse a la sala, antes de que comience el juego.

### Código de conexión (JavaScript)

```js
const key = "AB3X9Z"; // key recibida al crear/unirse
const ws = new WebSocket(`ws://localhost:4100/ws/${key}`);

ws.onopen = () => {
  console.log("Conectado a la sala:", key);
};

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  console.log("Evento recibido:", msg.event, msg.data);
};

ws.onclose = () => {
  console.log("Desconectado de la sala");
};

ws.onerror = (err) => {
  console.error("Error WS:", err);
};
```

---

### Eventos que el servidor envía al cliente

Todos los mensajes recibidos tienen el formato:
```json
{
  "event": "nombre_del_evento",
  "data": { ... }
}
```

---

#### Evento: `player_joined`
Se dispara cuando un nuevo jugador se une a la sala (via `POST /api/v1/room/join`).

```json
{
  "event": "player_joined",
  "data": {
    "game": { "id": 1, "key": "AB3X9Z", "status": "waiting", ... },
    "user": { "id": 2, "nickname": "Carlos", "avatar_url": "..." },
    "departure": { "id": 2, "game_id": 1, "user_id": 2, "score": 0, ... }
  }
}
```

**¿Cuándo usarlo?** Actualizar la lista de jugadores en la sala de espera.

---

#### Evento: `score_update`
Se dispara cuando cualquier jugador envía una respuesta (via `POST /api/v1/game/answer`). Contiene el scoreboard completo actualizado.

```json
{
  "event": "score_update",
  "data": [
    {
      "id": 1,
      "game_id": 1,
      "user_id": 1,
      "score": 20,
      "hits": 2,
      "total_time": 14,
      "user": { "id": 1, "nickname": "Moises", "avatar_url": "..." }
    },
    {
      "id": 2,
      "game_id": 1,
      "user_id": 2,
      "score": 10,
      "hits": 1,
      "total_time": 8,
      "user": { "id": 2, "nickname": "Carlos", "avatar_url": "..." }
    }
  ]
}
```

**¿Cuándo usarlo?** Actualizar la tabla de puntajes en tiempo real.

---

#### Evento: `game_finished`
Se dispara cuando el admin finaliza la partida (via `POST /api/v1/game/:id/finish`). Contiene el scoreboard final ordenado.

```json
{
  "event": "game_finished",
  "data": [
    {
      "id": 1,
      "score": 30,
      "hits": 3,
      "total_time": 22,
      "user": { "id": 1, "nickname": "Moises", "avatar_url": "..." }
    }
  ]
}
```

**¿Cuándo usarlo?** Mostrar la pantalla de resultados finales a todos los jugadores.

---

### Manejo de eventos (ejemplo completo)

```js
ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);

  switch (msg.event) {
    case "player_joined":
      agregarJugador(msg.data.user);
      break;

    case "score_update":
      actualizarScoreboard(msg.data);
      break;

    case "game_finished":
      mostrarResultadosFinales(msg.data);
      ws.close();
      break;
  }
};
```

---

## Estados del juego (`game.status`)

| Status | Descripción |
|---|---|
| `waiting` | Sala creada, esperando jugadores. Solo en este estado se puede unir alguien. |
| `started` | Partida en curso. El admin presionó "Comenzar". |
| `finished` | Partida finalizada. El admin presionó "Terminar". |

---

## Referencia rápida de endpoints

| Método | URL | Descripción |
|---|---|---|
| `GET` | `/api/v1/categories` | Listar categorías |
| `POST` | `/api/v1/room` | Crear sala (admin) |
| `POST` | `/api/v1/room/join` | Unirse a sala por key |
| `POST` | `/api/v1/room/:id/start` | Iniciar partida (admin) |
| `GET` | `/api/v1/game/:game_id/question` | Obtener siguiente pregunta |
| `POST` | `/api/v1/game/answer` | Enviar respuesta |
| `GET` | `/api/v1/game/:game_id/scoreboard` | Ver marcador |
| `POST` | `/api/v1/game/:game_id/finish?key=KEY` | Finalizar partida (admin) |
| `WS` | `/ws/:key` | WebSocket de la sala |
