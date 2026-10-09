# Fuzzy CLIPS MCP in go


Fuzzy clips it's a extension for CLIPS (rules engine) which allows to use certanty factors.

This is my implementation for model context protocol to allow the LLM interation with the FuzzyClips.


To use you should run the server. Use `FUZZY_MCP_ADDRES` to specify listen addres. 

The `localhost:9000` is default.


# Why?

I personally feels that's the knowdge sysetm are widely exteds capabilities of LLM's, because you are not saving the facts itself, but you can save the reflection itself. 

Fuzzy facts are really good support for making logic desigions in uncertain environments. 


# Build 

First, build and install fuzzyCLIPS lib.

```bash
cd ./fuzzy-clips/source/
make lib
make install-lib
```


Than, build go executable or run through the `go run main.go`

```bash
go build
./clips-mcp
```

Connect via mcp protocol to specifyed port

