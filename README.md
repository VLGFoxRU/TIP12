# ЭФМО-01-25 Буров М.А. ПР12

# Описание проекта
Подключение Swagger/OpenAPI. Автоматическая генерация документации

# Требования к проекту
* Go 1.25+
* Git

# Версия Go
<img width="317" height="55" alt="image" src="https://github.com/user-attachments/assets/43f9087b-95b9-4c7d-86e9-746258c45c63" />

# Цели:
1.	Освоить основы спецификации OpenAPI (Swagger) для REST API.
2.	Подключить автогенерацию документации к проекту из ПЗ 11 (notes-api).
3.	Научиться публиковать интерактивную документацию (Swagger UI / ReDoc) на эндпоинте GET /docs.
4.	Синхронизировать код и спецификацию (комментарии-аннотации → генерация) и/или «schema-first» (генерация кода из openapi.yaml).
5.	Подготовить процесс обновления документации (Makefile/скрипт).

# Краткое описание
В работе использован code-first подход, при котором аннотации добавляются непосредственно к исходному коду обработчиков (handlers). Инструмент swaggo/swag парсит эти комментарии и автоматически генерирует спецификацию OpenAPI в формате YAML и JSON, а также код инициализации для встраивания Swagger UI в сервер.

Этот подход выбран потому что:
- Быстрое внедрение в уже готовый проект
- Минимальные ручные действия с YAML
- Аннотации остаются рядом с кодом, что упрощает их поддержку

# Структура проекта
Дерево структуры проекта: 
```
pz11-notes-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── core/
│   │   ├── note.go
│   │   └── service/
│   │       └── note_service.go
│   ├── http/
│   │   ├── router.go
│   │   └── handlers/
│   │       └── notes.go
│   └── repo/
│       └── note_mem.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
└── go.sum
```

# Фрагменты аннотаций

```
// CreateNote godoc
// @Summary Создать новую заметку
// @Description Создаёт и сохраняет новую заметку в системе
// @Tags notes
// @Accept json
// @Produce json
// @Param input body core.CreateNoteRequest true "Данные новой заметки"
// @Success 201 {object} core.Note "Заметка успешно создана"
// @Failure 400 {object} map[string]string "Некорректные данные (пустые поля)"
// @Failure 500 {object} map[string]string "Ошибка при создании заметки"
// @Router /notes [post]
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request)
```

```
// GetAllNotes godoc
// @Summary Получить список всех заметок
// @Description Возвращает список всех заметок с поддержкой пагинации и фильтра по названию
// @Tags notes
// @Produce json
// @Param page query int false "Номер страницы (по умолчанию 1)" default(1)
// @Param limit query int false "Размер страницы (по умолчанию 10)" default(10)
// @Param q query string false "Поиск по названию (title)"
// @Success 200 {array} core.Note "Список заметок"
// @Header 200 {integer} X-Total-Count "Общее количество заметок в системе"
// @Failure 400 {object} map[string]string "Некорректные параметры пагинации"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /notes [get]
func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request)
```

```
// UpdateNote godoc
// @Summary Обновить заметку (частичное обновление)
// @Description Обновляет один или несколько полей заметки. Поля, не указанные в запросе, не изменяются
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор заметки для обновления"
// @Param input body core.UpdateNoteRequest true "Поля для обновления (все опциональны)"
// @Success 200 {object} core.Note "Заметка успешно обновлена"
// @Failure 400 {object} map[string]string "Некорректный ID или тело запроса"
// @Failure 404 {object} map[string]string "Заметка не найдена"
// @Failure 500 {object} map[string]string "Ошибка при обновлении заметки"
// @Router /notes/{id} [patch]
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request)
```

# Скриншоты

Работающая страница Swagger UI:

<img width="1466" height="918" alt="image" src="https://github.com/user-attachments/assets/6b0a6641-1e83-4b3e-967c-fb827d218153" />

# Команды генерации документации 

В корне проекта запустите:
```
swag init -g cmd/api/main.go -o docs
```
-g — главный файл, откуда начнётся парсинг аннотаций;
-o — каталог для вывода (docs/).
После выполнения появятся docs/docs.go, docs/swagger.json/swagger.yaml.

# Краткие выводы

Что удалось сделать:
- Успешная интеграция инструмента swaggo/swag в проект notes-api
- Полная документация всех 5 методов (CreateNote, GetNote, GetAllNotes, UpdateNote, DeleteNote) с описаниями, параметрами, примерами и кодами ответов
- Интерактивная Swagger UI доступна по адресу http://localhost:8080/docs/index.html

# Ответы на контрольные вопросы

1.	Чем отличается OpenAPI от Swagger?

OpenAPI — это открытый стандарт спецификации (формальное описание REST API в YAML/JSON формате), разработанный OpenAPI Initiative.

Swagger — это набор инструментов (Swagger UI, Swagger Editor, Swagger Codegen) для работы со спецификациями OpenAPI. Первоначально Swagger развивался как собственный проект, но позднее был передан в OpenAPI Initiative.

Простое объяснение: OpenAPI — это "что" (формат описания), Swagger — это "чем" (инструменты для работы с этим форматом). 

2.	В чём различие подходов code-first и schema-first? Плюсы/минусы.

Существует два распространённых подхода:
1) Code-First (аннотационный)
Программист пишет код API (хэндлеры, маршруты) и добавляет над функциями специальные комментарии-аннотации, из которых инструменты вроде swaggo/swag автоматически генерируют OpenAPI-файл.
Преимущества:
-	быстрое подключение к уже готовому проекту;
-	минимальная ручная работа с YAML;
-	простая поддержка актуальности при изменениях кода.
Недостатки:
-	ограниченная гибкость описания сложных сценариев;
-	возможны расхождения, если аннотации забыли обновить.
2) Schema-First (design-first)
Сначала проектируется OpenAPI-файл (openapi.yaml), а из него с помощью генераторов (oapi-codegen, openapi-generator, swagger-codegen) создаются шаблонные контроллеры, структуры данных и клиенты.
Преимущества:
-	строгий контроль над контрактом API;
-	удобен для крупных команд и публичных API;
-	возможность использовать единую спецификацию для разных языков.
Недостатки:
-	требует знаний YAML и структуры OAS;
-	сложнее поддерживать при быстрых итерациях разработки.

3.	Какие обязательные разделы содержит спецификация OpenAPI?

openapi — версия спецификации
info.title — название
info.version — версия
paths — описание маршрутов (может быть пусто, но должна быть секция)

4.	Для чего нужны components.schemas и как их переиспользовать в responses?

components.schemas — это раздел для определения переиспользуемых моделей данных (структур объектов), которые используются в запросах и ответах.

Переиспользование в responses (через $ref):
```
text
paths:
  /notes/{id}:
    get:
      responses:
        '200':
          description: Заметка найдена
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Note'  # Ссылка на компонент
        '404':
          description: Не найдена
```

5.	Что описывают аннотации @Param, @Success, @Failure, @Router, @Security?

- @Param	— Описание параметров запроса (path, query, body)	@Param id path int true "ID заметки"
- @Success	— Успешный ответ с кодом и типом	@Success 200 {object} core.Note
- @Failure	— Ошибка с кодом и типом	@Failure 404 {object} map[string]string
- @Router	— Маршрут и метод HTTP	@Router /notes/{id} [get]
- @Security	— Требуемая схема аутентификации	@Security BearerAuth

6.	Как опубликовать Swagger UI на отдельном префиксе (/docs) и ограничить к нему доступ?

Подключение на префиксе /docs:

```
import httpSwagger "github.com/swaggo/http-swagger"
import "example.com/pz11-notes-api/docs"  // Импорт пакета docs

// В main():
r := chi.NewRouter()

// Подключаем Swagger UI на /docs
r.Get("/docs/*", httpSwagger.WrapHandler)
```

Ограничение доступа (middleware):

```
// Middleware для проверки авторизации
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// Применяем middleware только к /docs
r.Route("/docs", func(r chi.Router) {
    r.Use(authMiddleware)
    r.Get("/*", httpSwagger.WrapHandler)
})
```

7.	Как поддерживать актуальность документации при изменениях кода?

1. Вручную при каждом изменении:

```
# После изменения аннотаций или обработчиков
swag init -g cmd/api/main.go -o docs
```

2. Через Makefile:

```
run: swagger
	go run ./cmd/api
```
  
8.	Как подключить Bearer-аутентификацию в спецификации и что изменится в UI?

Шаг 1: Определить SecurityScheme в cmd/api/main.go

```
// Package main Notes API server.
//
// @title Notes API
// @version 1.0
// @description REST API для управления заметками
// @BasePath /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer <token>
package main
```

Шаг 2: Добавить @Security над защищёнными методами

```
// GetNote godoc
// @Summary Получить заметку
// @Tags notes
// @Param id path int true "ID"
// @Success 200 {object} core.Note
// @Failure 404 {object} map[string]string
// @Router /notes/{id} [get]
// @Security BearerAuth          // Добавляем эту строку
func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) { ... }
```

Шаг 3: Сгенерировать документацию

```
swag init -g cmd/api/main.go -o docs
```
