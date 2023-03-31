package main

import (
	ctx "context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	wp "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var client *whatsmeow.Client

func loadenv() *envs {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("Loaded environment variables from .env file")
	// Create struct containing environment variables
	env := envs{
		OpenAIKey:   os.Getenv("OPENAIKEY"),
		AiImgKey:    os.Getenv("IMGAIKEY"),
		AiImgSecret: os.Getenv("IMGAISECRET"),
	}

	return &env
}

var env = loadenv()

// eventHandler handles incoming events and dispatches them to the appropriate functions
func eventHandler(evt interface{}) {
	// Check if event is a message
	if evt, ok := evt.(*events.Message); ok {
		// Handle image messages
		if evt.Message.ImageMessage != nil {
			go handleImageMessage(evt)
		}
		// Handle text messages
		if evt.Message.GetConversation() != "" {
			go handleConversation(evt)
		}
	}
}

// handleImageMessage handles incoming image messages by generating an AI-generated image response and sending it back
func handleImageMessage(evt *events.Message) {
	// Download image data from the message
	imageData, _ := client.Download(evt.Message.GetImageMessage())
	// Generate an AI-generated image using the downloaded data
	imageBytes, err := GenAIimg(imageData)
	if err != nil {
		// Send an error message if there was an error generating the image
		client.SendMessage(ctx.Background(), evt.Info.Chat, &wp.Message{
			Conversation: proto.String(err.Error()),
		})
		return
	}
	// Upload the generated image to WhatsApp
	uploadedImage, err := client.Upload(ctx.Background(), imageBytes, whatsmeow.MediaImage)
	if err != nil {
		// Send an error message if there was an error uploading the image
		client.SendMessage(ctx.Background(), evt.Info.Chat, &wp.Message{
			Conversation: proto.String(err.Error()),
		})
		return
	}
	// Create a message to send back containing the uploaded image
	msgToSend := &wp.Message{
		ImageMessage: &wp.ImageMessage{
			Caption: 		 proto.String("AI-generated image"),
			Url:               proto.String(uploadedImage.URL),
			DirectPath:        proto.String(uploadedImage.DirectPath),
			MediaKey:          uploadedImage.MediaKey,
			MediaKeyTimestamp: proto.Int64(time.Now().Unix()),
			Mimetype:          proto.String(http.DetectContentType(imageBytes)),
			FileEncSha256:     uploadedImage.FileEncSHA256,
			FileSha256:        uploadedImage.FileSHA256,
			FileLength:        proto.Uint64(uint64(len(imageBytes))),
			Height:            proto.Uint32(uint32(evt.Message.GetImageMessage().GetHeight())),
			Width:             proto.Uint32(uint32(evt.Message.GetImageMessage().GetWidth())),
		},
	}
	// Send the message back to the chat
	_, err = client.SendMessage(ctx.Background(), evt.Info.Chat, msgToSend)
	if err != nil {
		// Send an error message if there was an error sending the message
		client.SendMessage(ctx.Background(), evt.Info.Chat, &wp.Message{
			Conversation: proto.String(err.Error()),
		})
	}
}

// handleConversation handles incoming text messages by generating an AI-generated text response and sending it back
func handleConversation(evt *events.Message) {
	// Generate an AI-generated text response using the message text
	msg, err := GetAiTextResponse(evt.Message.GetConversation())
	if err != nil {
		// Send an error message if there was an error generating the text response
		client.SendMessage(ctx.Background(), evt.Info.Chat, &wp.Message{
			Conversation: proto.String(err.Error()),
		})
		return
	}
	// Create a message to send back containing the generated text
	client.SendMessage(ctx.Background(), evt.Info.Chat, &wp.Message{
		ExtendedTextMessage: &wp.ExtendedTextMessage{
			Text: proto.String(msg),
			ContextInfo: &wp.ContextInfo{
				QuotedMessage: &wp.Message{
					Conversation: proto.String(evt.Message.GetConversation()),
				},
				StanzaId:    proto.String(evt.Info.ID),
				Participant: proto.String(evt.Info.Sender.ToNonAD().String()),
			},
		},
	})
}

func main() {
	// Set up logging for the database
	dbLog := waLog.Stdout("Database", "INFO", true)
	// Create a new SQL store container
	container, err := sqlstore.New("sqlite3", "file:whatgpt.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	// Get the first device from the container
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	// Set up logging for the client
	clientLog := waLog.Stdout("Client", "INFO", true)

	client = whatsmeow.NewClient(deviceStore, clientLog)

	client.AddEventHandler(eventHandler)

	// Check if client is already logged in
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(ctx.Background())
		// Connect to WhatsApp
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		// Print the QR code to the console
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Connect to WhatsApp if already logged in
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
