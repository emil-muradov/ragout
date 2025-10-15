# AI - Assistant

An example of implementing RAG (Retrieval Augmented Generation)  with Go, YandexGPT as LLM, LangChain and Qdrant as vector database. As a domain knowledge base I used a plain text file.
The domain itself is real estate in Saint Petersburg city (Russia). It is best to provide a request in russian language.

---
# ИИ Ассистент

Пример реализации RAG (Retrieval Augmented Generation) с использованием Go, YandexGPT в качестве LLM, LangChain и Qdrant в качестве векторной базы данных. В качестве базы знаний предметной области я использовал обычный текстовый файл.
Сама предметная область — недвижимость в городе Санкт-Петербург (Россия). Лучше всего подавать запрос на русском языке.

Make sure you are in a project root <br>
Убедись, что ты в корне проекта

To start the server use cli: <br>
Для запуска сервера введи следующие команды:
```
docker build -f cmd/ai-assistant/Dockerfile .
docker compose -f cmd/ai-assistant/docker-compose.yaml build --no-cache
docker compose -f cmd/ai-assistant/docker-compose.yaml up --watch
```

To test how it works change `msg` prop in query below <br>
Пример запроса:
```
curl -X POST -d '{ "msg": "Нужна вместительная квартира в СПб"}' http://localhost:8080/ask
```
