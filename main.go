package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"
)

type DiskStruct struct {
	Free uint64 `json:"free"`
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	disk := DiskFreeSpace(os.Getenv("DISK"))
	size := int64(disk.Free) / int64(1024*1024)
	limit, err := strconv.ParseInt(os.Getenv("LIMIT"), 10, 64)
	if err != nil {
		fmt.Println("ERR LIMIT:", err)
	}

	channelId, err := strconv.Atoi(os.Getenv("CHANNEL_ID"))
	if err != nil {
		fmt.Println("Channel parse:", err)
	}

	if size < limit {
		sendMessageTelegram(channelId, os.Getenv("MESSAGE"), size)
	}

}

func DiskFreeSpace(path string) (disk DiskStruct) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	return
}

func sendMessageTelegram(chatId int, text string, size int64) (string, error) {

	log.Printf("Sending %s to chat_id: %d", text, chatId)
	var telegramApi string = "https://api.telegram.org/bot" + os.Getenv("TOKEN") + "/sendMessage"
	//msg := strings.Join([]string{"err:", err.Error(), "message:", message}, "")

	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}
