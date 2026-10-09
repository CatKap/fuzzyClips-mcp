package handlers
/*
#cgo CFLAGS: -I${SRCDIR}/../fuzzy-clips/source
#cgo LDFLAGS: -L${SRCDIR}/../fuzzy-clips/source -lfzclips -lm

#include <clips.h>
#include <strngfun.h>
#include <stdlib.h>

static void *GetpArgsValue(DATA_OBJECT_PTR arg)
{
    return GetpValue(arg);
}

static int GetpArgsType(DATA_OBJECT_PTR arg)
{
    return GetpType(arg);
}

static void SetReturnValue(
    DATA_OBJECT_PTR returnValue,
    int type,
    void *value)
{
    SetpType(returnValue, type);
    SetpValue(returnValue, value);
}


void CountDeftemplates(
    DATA_OBJECT_PTR returnValue)
{
 		DATA_OBJECT_PTR arg;
    char *deftemplateName;
    struct deftemplate *df;
    struct fact *ft;
    int count = 0;
    
    int argCount = RtnArgCount();

    if(argCount != 1){
      PrintRouter("error", "Shold be exactly one argument");
      SetEvaluationError(TRUE);
      return;
    }

    RtnUnknown(0, arg);

 	  if (!arg){
      PrintRouter("error", "Panic: no argument!");
      SetEvaluationError(TRUE);
      return;
    } 
 	  if (arg->type != STRING){
      PrintRouter("error", "Argument should be string!");
      SetEvaluationError(TRUE);
      return;
    } 

    deftemplateName = arg->value;

		printf("%s\n", deftemplateName);

    df = FindDeftemplate(deftemplateName);

    if (!df)
    {
        char retbuff[1024];

        snprintf(
            retbuff,
            sizeof(retbuff),
            "No deftemplate %s found\n",
            deftemplateName
        );

        PrintRouter("wclips", retbuff);

        SetpType(returnValue, INTEGER);
        SetpValue(returnValue, AddLong((long) 0));

        return;
    }

    ft = NULL;

    while ((ft = GetNextFact(ft)) != NULL)
    {
        if (ft->whichDeftemplate == df)
            count++;
    }

    SetpType(returnValue, INTEGER);
    SetpValue(returnValue, AddLong((long) count));
}



void RegisterUserFunctions(void)
{
    DefineFunction2(
        "count-deftemplate",
        'i',
        PTIF CountDeftemplates,
        "CountDeftemplates",
        "11"
    );
}

*/
import "C"

import (
	"context"
	"unsafe"
	"fmt"
	"sync"
	"os"
	"io"
	"bytes"
	"clips-mcp/facts"
  "github.com/mark3labs/mcp-go/mcp"
)


var ClipsMutex sync.Mutex


func InitClips(){
	C.InitializeEnvironment()
	C.RegisterUserFunctions()
}

func RunClips(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ClipsMutex.Lock()
	  C.Run(-1) 
		ClipsMutex.Unlock()
    return mcp.NewToolResultText(fmt.Sprintf("environment had runned")), nil
}

func ResetClips(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ClipsMutex.Lock()
	  C.Reset() 
		ClipsMutex.Unlock()
		fmt.Println("Clips reseted")
    return mcp.NewToolResultText(fmt.Sprintf("Reseted")), nil
}

func ClearClips(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ClipsMutex.Lock()
	  C.Clear() 
		ClipsMutex.Unlock()
		fmt.Println("Clips cleared")
    return mcp.NewToolResultText(fmt.Sprintf("Cleared")), nil
}

func EvaluateStringHandler (ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    old := os.Stdout // keep backup of the real stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    outC := make(chan string)
    // copy the output in a separate goroutine so printing can't block indefinitely
    go func() {
        var buf bytes.Buffer
        io.Copy(&buf, r)
        outC <- buf.String()
    }()

    exe, err := request.RequireString("exestring")

		var result C.DATA_OBJECT
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }

		cmd := C.CString(exe)
		defer C.free(unsafe.Pointer(cmd))
		
		ClipsMutex.Lock()
		ok := C.Eval(cmd, &result)
		ClipsMutex.Unlock()

    w.Close()
    os.Stdout = old // restoring the real stdout
    out := <-outC
		fmt.Print(out)

		if ok != 1 {
				// back to normal state
				return mcp.NewToolResultError(fmt.Sprintf("Error during evaluation: %s", out)), nil
		}	
    return mcp.NewToolResultText(fmt.Sprintf("Ok")), nil
}

func RetriveByNameHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
  pattern, err := request.RequireString("pattern")

  if err != nil {
      return mcp.NewToolResultError(err.Error()), nil
  }

	searchInFields, err := request.RequireBool("search_in_fields")

	if err != nil {
		fmt.Println("Default behavior of search_in_fields is set to ffalse due to error %s\n", err.Error())		
		searchInFields = false
	}
	
	ClipsMutex.Lock()
	allFacts := facts.Facts()	
	ClipsMutex.Unlock()
	matched, serr := facts.MatchFinder(pattern, allFacts, searchInFields)

	if serr != ""{
    return mcp.NewToolResultError(serr), nil
	}
  return mcp.NewToolResultText(facts.StringifyFacts(matched)), nil
}

func LoadFileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

    name, err := request.RequireString("filename")

    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }


    old := os.Stdout // keep backup of the real stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    outC := make(chan string)
    // copy the output in a separate goroutine so printing can't block indefinitely
    go func() {
        var buf bytes.Buffer
        io.Copy(&buf, r)
        outC <- buf.String()
    }()

		ClipsMutex.Lock()
		ret := C.Load(C.CString(name))
		ClipsMutex.Unlock()

    w.Close()
    os.Stdout = old // restoring the real stdout
    out := <-outC
		fmt.Print(out)

		if  ret != 1 {
			return mcp.NewToolResultError(fmt.Sprintf("Cannot load this file: ", out)), nil
		}
    return mcp.NewToolResultText(fmt.Sprintf("File loaded")), nil
}



func RetriveFactsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {


	ClipsMutex.Lock()
	factsInfo := facts.Facts()
	ClipsMutex.Unlock()
	fmt.Print(factsInfo)
  return mcp.NewToolResultText(facts.StringifyFacts(factsInfo)), nil
}



func LoadFolderHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    foldername, err := request.RequireString("foldername")

    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }

	  info, errors := facts.ProcessDirectory(&ClipsMutex, foldername)	

		if len(errors) != 0 {
			return mcp.NewToolResultError(fmt.Sprintf("Errors while loading diretory in files %v \n", errors)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Files sucseefully loaded, check the additional agent info:\n %s", info)), nil
}







