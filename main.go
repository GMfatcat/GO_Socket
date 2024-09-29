package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coder/websocket"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var conn *websocket.Conn

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme()) // 設定為暗黑主題
	w := a.NewWindow("GO WebSocket UI")

	// 創建 IP 和 Port 輸入框
	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("輸入伺服器 IP")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("輸入伺服器 Port")

	statusLabel := widget.NewLabel("未連接")

	// 創建連接按鈕
	connectButton := widget.NewButton("連接", func() {
		ip := ipEntry.Text
		port := portEntry.Text
		if ip != "" && port != "" {
			err := connectWebSocket(ip, port, statusLabel)
			if err == nil {
				showChatUI(w, statusLabel)
			} else {
				statusLabel.SetText(fmt.Sprintf("連接失敗: %v", err))
			}
		} else {
			statusLabel.SetText("請輸入正確的 IP 和 Port")
		}
	})

	// 首頁佈局：IP、Port 輸入框和連接按鈕
	startLayout := container.NewVBox(
		widget.NewLabel("請輸入伺服器 IP 和 Port"),
		ipEntry,
		portEntry,
		connectButton,
		statusLabel,
	)

	w.SetContent(startLayout)
	w.ShowAndRun()
}

// 連接 WebSocket 伺服器
func connectWebSocket(ip, port string, statusLabel *widget.Label) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	url := fmt.Sprintf("ws://%s:%s/ws", ip, port)
	var err error
	conn, _, err = websocket.Dial(ctx, url, nil)
	if err != nil {
		return err
	}

	statusLabel.SetText("已成功連接到伺服器")
	return nil
}

// 顯示聊天 UI
func showChatUI(w fyne.Window, statusLabel *widget.Label) {
	// 聊天區域：顯示聊天記錄的 Label
	chatDisplay := widget.NewLabel("聊天記錄:\n")

	// 輸入消息的文本框
	input := widget.NewEntry()
	input.SetPlaceHolder("輸入消息...")

	// 發送消息按鈕
	sendButton := widget.NewButton("發送", func() {
		msg := input.Text
		if msg != "" {
			sendMessage(msg, chatDisplay)
			input.SetText("")
		}
	})

	// 結束連接按鈕
	closeButton := widget.NewButton("結束連接", func() {
		closeConnection()
		statusLabel.SetText("已關閉連接")
		chatDisplay.SetText("聊天記錄:\n")       // 清空聊天記錄
		showStartPage(w, statusLabel, input) // 回到開始頁面
	})

	// 聊天介面佈局
	chatLayout := container.NewVBox(
		statusLabel,
		chatDisplay,
		input,
		sendButton,
		closeButton, // 將結束連接按鈕加入佈局
	)

	// 更新窗口內容為聊天介面
	w.SetContent(chatLayout)
}

// 顯示開始頁面
func showStartPage(w fyne.Window, statusLabel *widget.Label, input *widget.Entry) {
	// 清空消息輸入框
	input.SetText("")

	// 重新創建 IP 和 Port 輸入框
	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("輸入伺服器 IP")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("輸入伺服器 Port")

	// 重新創建連接按鈕
	connectButton := widget.NewButton("連接", func() {
		ip := ipEntry.Text
		port := portEntry.Text
		if ip != "" && port != "" {
			err := connectWebSocket(ip, port, statusLabel)
			if err == nil {
				showChatUI(w, statusLabel)
			} else {
				statusLabel.SetText(fmt.Sprintf("連接失敗: %v", err))
			}
		} else {
			statusLabel.SetText("請輸入正確的 IP 和 Port")
		}
	})

	// 重新顯示最初的頁面
	w.SetContent(container.NewVBox(
		widget.NewLabel("請輸入伺服器 IP 和 Port"),
		ipEntry,
		portEntry,
		connectButton,
		statusLabel,
	))
}

// 發送消息並顯示在聊天區域
func sendMessage(msg string, chatDisplay *widget.Label) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := conn.Write(ctx, websocket.MessageText, []byte(msg))
	if err != nil {
		chatDisplay.SetText(fmt.Sprintf("%s\n發送消息錯誤: %v", chatDisplay.Text, err))
		return
	}

	_, receivedMsg, err := conn.Read(ctx)
	if err != nil {
		chatDisplay.SetText(fmt.Sprintf("%s\n接收消息錯誤: %v", chatDisplay.Text, err))
		return
	}

	chatDisplay.SetText(fmt.Sprintf("%s\n收到消息: %s", chatDisplay.Text, receivedMsg))
}

// 當需要結束時再關閉連接，例如在應用程式結束時
func closeConnection() {
	err := conn.Close(websocket.StatusNormalClosure, "結束連接")
	if err != nil {
		// 處理關閉錯誤
		fmt.Println(err)
	}
}
