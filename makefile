# please install wails first, check https://wails.io for help

init-dev:
	brew install wails

build-windows:
	wails build -webview2 embed -nsis -platform windows/amd64

build-mac:
	wails build -platform darwin/amd64

build-linux:
	wails build -platform linux/amd64
