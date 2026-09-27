package main
/*
#cgo CFLAGS: -I/home/kapitan/Downloads/gits/FuzzyCLIPS/source
#cgo LDFLAGS: -L/home/kapitan/Downloads/gits/FuzzyCLIPS/source -lfzclips -lm

#include <clips.h>
#include <stdlib.h>
*/
import "C"

import (
  "context"
  "fmt"
	"os"
	"unsafe"
  "github.com/mark3labs/mcp-go/mcp"
  "github.com/mark3labs/mcp-go/server"
)
//func main() {
//	// 1. Инициализация (старый API)
//
//
//	// 3. Reset
//	C.Reset()
//
//	// 4. Assert факта
//	fact := C.CString(`(obligation (id "OFZ-26238") (yield_at_close 7.85) (simple_yield 0.12) (calculated yes))`)
//	defer C.free(unsafe.Pointer(fact))
//
//	C.AssertString(fact)
//
//	// 5. Запуск
//	C.Run(-1)
//
//	fmt.Println("CLIPS (Fuzzy) отработал")
//}


func main() {
    // Create a new MCP server
    s := server.NewMCPServer(
        "Fuzzy CLIPS MCP 🚀",
        "6.10",
        server.WithToolCapabilities(false),
    )

		C.InitializeEnvironment()

		// 2. Загрузка правил
		file := C.CString("rules.clp")
		defer C.free(unsafe.Pointer(file))

    // Add tool
    loadFile := mcp.NewTool("load_file",
        mcp.WithDescription("Load file into the CLIPS"),
        mcp.WithString("filename",
            mcp.Required(),
            mcp.Description("Fileame with appropriate CLIPS syntax"),
        ),
    )

    loadFolder := mcp.NewTool("load_folder",
        mcp.WithDescription("Load files with .clp or .clips extension to the CLIPS"),
        mcp.WithString("foldername",
            mcp.Required(),
            mcp.Description("Fileame with appropriate CLIPS syntax"),
        ),
    )

    resetEnv := mcp.NewTool("reset_env",
        mcp.WithDescription("Reset CLIPS environment"),
    )

    run := mcp.NewTool("run",
        mcp.WithDescription("Run the CLIPS environment"),
    )

    //retriveFact := mcp.NewTool("retrive_fact",
    //    mcp.WithDescription("Run the CLIPS environment"),
    //    mcp.WithString("factname",
    //        mcp.Required(),
    //        mcp.Description("The name of the fact or template"),
    //    ),
    //)

    s.AddTool(loadFile, loadFileHandler)
		s.AddTool(resetEnv, clipsReset)
		s.AddTool(run, runClips)


    // Start the stdio server
    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("Server error: %v\n", err)
    }
}


func runClips(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	  C.Run(-1) 
    return mcp.NewToolResultText(fmt.Sprintf("environment had runned")), nil
}


func clipsReset(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	  C.Reset() 
    return mcp.NewToolResultText(fmt.Sprintf("Reseted")), nil
}

func loadFileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    name, err := request.RequireString("filename")

    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }

		if C.Load(C.CString(name)) != 1 {
        return mcp.NewToolResultError("Cannot load this file"), nil
		}

    return mcp.NewToolResultText(fmt.Sprintf("File loaded")), nil
}


func loadFolderHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    name, err := request.RequireString("foldername")

    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }

		info, err := os.Stat(name)
			if os.IsNotExist(err) {
        return mcp.NewToolResultError("Folder does not exist"), nil
			}
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Error checking path: %s", err)), nil
			}
		
			if ! info.IsDir() {
				return mcp.NewToolResultError("Path exists, but it is a file, not a folder"), nil
			}

		if C.Load(C.CString(name)) != 1 {
		}

    return mcp.NewToolResultText(fmt.Sprintf("File loaded")), nil
}




