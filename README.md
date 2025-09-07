# playfair-cipher

## About
This project was written as a university homework. It shows how text ciphering/deciphering occurs in real time. An application works fully in terminal and allows to read and save ciphered/deciphered text to files. As the application had been written on Golang, that allowed to use one source code to create executable files for Linux and Windows.

You can see in [`images/`](https://github.com/akaspb/playfair-cipher/tree/main/images), how it works in terminal.


## Instruction
To run application in Linux, use `sudo ./bin/playfair`.

To run it in windows, use `.\bin\playfair.exe`.

Application works identically on both OS's.

1. First! You need to write key and save it. To do that, you need to fill `Key` field in `Settings` tab.
In application you can move between fields (up and down) with `↑` and `↓` buttons pressing.
Key must contain characters from `Alphabet` field excluding `Separator character`.
After that, press `ctrl + s` to save settings.
2. Now you can move between tabs using `tab` and `shift + tab`.
3. To do cipher/decipher text on `Cipher`/`Decipher` tabs write text in input.
You can use `ctrl + v` to put text from clipboard, `ctrl + s` to load processed text to clipboard.
Also, you can use `ctrl + d` to delete text from input.
During text processing application might show additional information about decipher errors and so on.
4. To save processed text to file, first, you should write file name to the field `file name with extension` and then press `ctrl + w`.
To load text from file, you should also write file name and press `ctrl + r`.
5. To close the application, use `esc` button. 
