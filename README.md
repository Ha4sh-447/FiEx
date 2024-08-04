# FiEx
A fast file explorer for windows inspired from [FileExplorer](https://github.com/conaticus/FileExplorer).

Main aim of the project is to making searching through volumes much faster as compared to windows file explorer.

## What it does?
- Cache the entire file system
- Traverse folders
- Search files/folders

## Run Locally
### Prerequisites
- [Go](https://go.dev/doc/install)
- [fyne](https://docs.fyne.io/started/)

### Steps
- Clone the repo
- Run ```go mod tidy```
- ```cd cmd```
- ```go run explorer.go```
- Build the executable
    - ```fyne package -os windows -icon icon.png```

## TODO
- [ ] Optimize Caching - currently takes really long time.
- [ ] Better file scoring
- [ ] Create and delete files/folders
