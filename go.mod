module github.com/wkalt/chatbot

go 1.21.0

require (
	github.com/joho/godotenv v1.5.1
	github.com/slack-go/slack v0.12.3
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/lib/pq v1.10.9
)

replace github.com/wkalt/migrate => ../../projects/migrate
