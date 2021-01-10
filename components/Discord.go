package components

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goinventory/models"
	"net/http"
)

type Discord struct {
	Webhook string
}

func (discord *Discord) SendNotification(message string) {
	if discord.Webhook == "" {
		return
	}
	messageout := models.DiscordMessage{Username: "inventoryBot", Content: message}

	bytesout, _ := json.Marshal(messageout)
	req, err := http.NewRequest("POST", discord.Webhook, bytes.NewBuffer(bytesout))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("problem sending web hook:", err)
	}
	defer resp.Body.Close()
}
