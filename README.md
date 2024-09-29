# GO WEB SOCKET CONNECT UI

# GO Compile First then fyne Compile

go build -o ChatApp.exe main.go

fyne package --exe ChatApp.exe --os windows --icon lick.png --name GO_WebSocket_UI

# [ref](https://www.cnblogs.com/holychan/p/17299042.html)
