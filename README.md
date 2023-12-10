# Chatbot

This is an extensible slack chatbot.

Instructions for use:
1. Clone the repo.
2. Copy `.env-example` and `manifest-example.yaml` to `.env` and `manifest.yaml`.
3. Obtain slack bot and API keys for the env file and create yourself an
   application in slack. The manifest.yaml file can be used as a guide.
4. Run the app with `go run main.go` and look at `commands.go` in the
   `external` directory. Extend this file with new commands that implement your
   bot logic.
