package main

import (
	"clips-mcp/handlers"
	"context"
	"fmt"
	"os"

	//"clips-mcp/facts"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)
func getenv(key, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}

func main() {
    // Create a new MCP server
    s := server.NewMCPServer(
        "Fuzzy CLIPS MCP 🚀",
        "6.10",
        server.WithToolCapabilities(false),
    )

		handlers.InitClips()
		//facts.ProcessDirectory(&handlers.ClipsMutex, "/home/kapitan/Cooooodiiiing/FinTech/InvestmentKnowledgeDatabase/for-agent/")

		// 2. Загрузка правил

		//var result C.DATA_OBJECT
		//C.Eval(C.CString("(assert (person Bob))"), &result)
		//
		//f := facts.Facts()
		//fmt.Print("Finded\n")
		//fmt.Print(facts.MatchFinder("Bob", f, true))
		//C.Reset()


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
        mcp.WithDescription("Set deffacts, delete assertion facts and inherited facts in CLIPS environment"),
    )

    clearEnv := mcp.NewTool("clear_env",
        mcp.WithDescription("Tottaly delete all in CLIPS environment"),
    )

    run := mcp.NewTool("run",
        mcp.WithDescription("Run the CLIPS environment"),
    )

    description := mcp.NewTool("description",
        mcp.WithDescription("Get basic syntax rules description"),
    )


    retriveFacts := mcp.NewTool("retrive_all_facts",
        mcp.WithDescription("Return all facts"),
    )

		retriveByName := mcp.NewTool("find_facts_by_name",
        mcp.WithDescription("Find facts by given string, can search through the fields if parameter specifyed. Allows regexp"),
				mcp.WithString("pattern",
					mcp.Required(),
					mcp.Description("Name or regexp pattern to search"),
				),

				mcp.WithBoolean("search_in_fields",
					mcp.DefaultBool(false),
					mcp.Description("Name or regexp pattern to search"),
				),
		)

    evaluateString := mcp.NewTool("evaluate_string",
        mcp.WithDescription("Directry evalute command in clips syntax (supports fuzzy CF)"),
        mcp.WithString("exestring",
            mcp.Required(),
            mcp.Description("CLIPS-syntax string to execute"),
        ),
    )

    //retriveFact := mcp.NewTool("retrive_fact",
    //    mcp.WithDescription("Run the CLIPS environment"),
    //    mcp.WithString("factname",
    //        mcp.Required(),
    //        mcp.Description("The name of the fact or template"),
    //    ),
    //)

    s.AddTool(loadFile,   		handlers.LoadFileHandler)
		s.AddTool(resetEnv,   		handlers.ResetClips)
		s.AddTool(clearEnv,   		handlers.ClearClips)
		s.AddTool(run, 			  		handlers.RunClips)
		s.AddTool(loadFolder, 		handlers.LoadFolderHandler)
		s.AddTool(retriveFacts, 	handlers.RetriveFactsHandler)
		s.AddTool(evaluateString, handlers.EvaluateStringHandler)
		s.AddTool(retriveByName,  handlers.RetriveByNameHandler)

		s.AddTool(description,    DescriptionHandler)

		httpOpts := []server.StreamableHTTPOption{}
		httpServer := server.NewStreamableHTTPServer(s, httpOpts...)
    // Start the stdio server
		
		if err := httpServer.Start(getenv("FUZZY_MCP_ADDRES", "localhost:9000")); err != nil {
        fmt.Printf("Server error: %v\n", err)
    
		}

}

func DescriptionHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    return mcp.NewToolResultText("This is FuzzyCLIPS mcp server. To use CF in facts declaration you should write (assert (<fact data>) CF %2.f). To use CF in rules use (defrule <rule-name> (declare CF %2.f) <<remain-rule-part>). To access cf of a fact use function (get-cf fact)"), nil
}

