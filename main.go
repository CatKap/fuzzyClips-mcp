package main
/*
#cgo CFLAGS: -I/home/kapitan/Downloads/gits/FuzzyCLIPS/source
#cgo LDFLAGS: -L/home/kapitan/Downloads/gits/FuzzyCLIPS/source -lfzclips -lm

#include <clips.h>
#include <stdlib.h>
*/
import "C"

import "clips-mcp/facts"

import (

  "context"
  "fmt"
	"os"
	"unsafe"
  "github.com/mark3labs/mcp-go/mcp"
  "github.com/mark3labs/mcp-go/server"
)

const FACT_BUFF_SIZE = 4096

func main() {
    // Create a new MCP server
    s := server.NewMCPServer(
        "Fuzzy CLIPS MCP 🚀",
        "6.10",
        server.WithToolCapabilities(false),
    )

		C.InitializeEnvironment()

		// 2. Загрузка правил
		//file := C.CString("rules.clp")
		//

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

    retriveFact := mcp.NewTool("retrive_fact",
        mcp.WithDescription("Run the CLIPS environment"),
        mcp.WithString("factname",
            mcp.Required(),
            mcp.Description("The name of the fact or template"),
        ),
    )

    s.AddTool(loadFile, loadFileHandler)
		s.AddTool(resetEnv, clipsReset)
		s.AddTool(run, runClips)
		s.AddTool(loadFolder, loadFolderHandler)

		httpOpts := []server.StreamableHTTPOption{}
		httpServer := server.NewStreamableHTTPServer(s, httpOpts...)
    // Start the stdio server
		
		if err := httpServer.Start(":9000"); err != nil {
        fmt.Printf("Server error: %v\n", err)
    
		}

}


func runClips(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	  C.Run(-1) 
    return mcp.NewToolResultText(fmt.Sprintf("environment had runned")), nil
}


func clipsReset(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	  C.Reset() 
		fmt.Println("Clips reseted")
    return mcp.NewToolResultText(fmt.Sprintf("Reseted")), nil
}

func loadFileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    name, err := request.RequireString("filename")

    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }


		/********************************************************/
		/* Load: C access routine for the load command. Returns */
		/*   0 if the file couldn't be opened, -1 if the file   */
		/*   was opened but an error occurred while loading     */
		/*   constructs, and 1 if the file was opened and no    */
		/*   errors occured while loading.                      */
		/********************************************************/

		if C.Load(C.CString(name)) != 1 {
        return mcp.NewToolResultError("Cannot load this file"), nil
		}


    return mcp.NewToolResultText(fmt.Sprintf("File loaded")), nil
}



func retriveFactHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
  factname, err := request.RequireString("factname")

  if err != nil {
      return mcp.NewToolResultError(err.Error()), nil
  }

	//var factBuff [FACT_BUFF_SIZE]C.char
	// Get fact ptr
			// Fuction name is just for error chicking. The argument must be the symbol
			// from struct expr
//   static struct fact *getFactPtr(struct expr *theArgument, char *functionName);


	//C.GetFactPPForm(factBuff, FACT_BUFF_SIZE, )

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
