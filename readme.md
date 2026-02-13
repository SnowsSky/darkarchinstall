# DarckArchInstall | An easy way to install DarkArchLinux
 # 1 - Downloading 🛜.
  - Clone the repo
  `git clone https://github.com/SnowsSky/darkarchinstall.git`
 # 2 - Compiling ⬇.
   ## 2.1 - Requirements`
   - You need to have go>=1.25.5 installed to compile it.
   `sudo pacman -Syu go`
   ## 2.2 - Build
   - Just run : 
   `go build -o darkarchinstall`
 # 3 - Install it
  - To install it, you just need to move the binary file to `/usr/bin/`
  - Or install it with : 
  `install -Dm755 darkarchinstall "/usr/bin/darkarchinstall""`