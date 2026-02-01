# AI-Assistant

An example of implementing RAG (Retrieval Augmented Generation) with Go, YandexGPT as LLM, LangChain and Qdrant as vector database. As a domain knowledge base I used a plain text file.
The domain itself is real estate in Saint Petersburg city (Russia). It is best to provide a request in russian language.

Setup:
1. Go to yandex cloud console and create a new service account, get API key for it
2. Create .env file and copy content from .env.example and add your yandex cloud credentials
3. Makue sure to have docker installed and running

Start server:
```
cd ai-assistant
docker compose -f compose.dev.yml up
```


Example of query:
```
curl -X POST -d '{ "question": "Нужна вместительная квартира в СПб"}' http://localhost:8080/ask
```

---
# ИИ-Ассистент

Пример реализации RAG (Retrieval Augmented Generation) с использованием Go, YandexGPT в качестве LLM, LangChain и Qdrant в качестве векторной базы данных. В качестве базы знаний предметной области я использовал обычный текстовый файл.
Сама предметная область — недвижимость в городе Санкт-Петербург (Россия). Лучше всего подавать запрос на русском языке.

Предварительные шаги:
1. Перейдите в консоль Yandex Cloud и создайте новый сервисный аккаунт, получите API ключ для него
2. Создайте файл .env и скопируйте содержимое из .env.example и добавьте ваши данные из Yandex Cloud
3. Убедитесь, что у вас установлен и запущен Docker

Для запуска сервера введи следующие команды:
```
cd ai-assistant
docker compose -f compose.dev.yml up
```


Пример запроса:
```
curl -X POST -d '{ "question": "Нужна вместительная квартира в СПб"}' http://localhost:8080/ask
```
