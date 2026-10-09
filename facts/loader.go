package facts 

import (
	"strings"
	"path/filepath"
	"fmt"
	"os"
	"sync"
	"sort"

)

/*
#cgo CFLAGS: -I/home/kapitan/Downloads/gits/FuzzyCLIPS/source
#cgo LDFLAGS: -L/home/kapitan/Downloads/gits/FuzzyCLIPS/source -lfzclips -lm

#include <clips.h>
#include <strngfun.h>
#include <stdlib.h>

void UserFunctions(void){

}

*/
import "C"


// naturalLess compares two strings the way humans expect:
//   "00_", "01_", ..., "09_", "10_", "11_", ...
func naturalLess(a, b string) bool {
    // Strip any leading path components – we only want to compare the base name.
    // If you need to sort full paths, you can omit this step.
    fa := filepath.Base(a)
    fb := filepath.Base(b)

    // Simple case‑insensitive numeric prefix detection.
    // This works for patterns like "00_*", "01_*", etc.
    // For a more robust natural sort you could use a library (e.g. github.com/iamacarpet/go-natural).
    i, j := 0, 0
    for i < len(fa) && j < len(fb) {
        // Skip non‑digits
        if !isDigit(fa[i]) && !isDigit(fb[i]) {
            if fa[i] != fb[i] {
                return fa[i] < fb[i]
            }
            i++
            j++
            continue
        }
        // Both are digits – consume the whole number
        if isDigit(fa[i]) && isDigit(fb[i]) {
            na := consumeNumber(fa[i:])
            nb := consumeNumber(fb[j:])
            if na != nb {
                return na < nb
            }
            i += len(fmt.Sprintf("%d", na))
            j += len(fmt.Sprintf("%d", nb))
            continue
        }
        // One side is digit, the other isn't – digits come first
        if isDigit(fa[i]) {
            return true
        }
        return false
    }
    // If we reach here, one string is a prefix of the other.
    return len(fa) < len(fb)
}

// isDigit reports whether r is an ASCII digit.
func isDigit(b byte) bool { return b >= '0' && b <= '9' }

// consumeNumber reads a leading integer from s and returns its value.
func consumeNumber(s string) int {
    var n int
    for i := 0; i < len(s) && isDigit(s[i]); i++ {
        n = n*10 + int(s[i]-'0')
    }
    return n
}

// sortedWalk returns a slice of all files (or dirs+files) under root,

func sortedWalk(root string) (entries []string, agentInfo string, err_ret error) {
		agentInfo = ""

    err := filepath.WalkDir(root, func(path_r string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }
        // You can decide here whether you want only files, only dirs, or both.
        // Example: only *.clp files:
				name := d.Name()
				lowerName := strings.ToLower(name)
				fmt.Printf("Walking %s\n", name)

				ext := strings.ToLower(filepath.Ext(name))
				if ext == ".clp" || ext == ".clips" || ext == ".cls" {
					 entries = append(entries, path_r)	
				}

				
				if lowerName == "agent.md" || lowerName == "claude.md" {
					data, readErr := os.ReadFile(path_r)
					if readErr != nil {
						return nil
					}
					agentInfo = fmt.Sprintf("=== %s ===\n%s\n\n", name, string(data))
					return nil
				}

        //if d.IsDir() {
        //	entries = append(entries, path)
        //}
        return nil
    })

    if err != nil {
        return nil, agentInfo, err
    }

    sort.Slice(entries, func(i, j int) bool { return naturalLess(entries[i], entries[j]) })
    return entries, agentInfo, nil
}

// ---------------------------------------------------------------------
// Example usage with the fuzzy‑clips‑mcp toolset (pseudo‑code)
// ---------------------------------------------------------------------
func ProcessDirectory(clipsMutex *sync.Mutex, dir string) (agentContent string, errors []string) {
    files, agentInfo, err := sortedWalk(dir)
    if err != nil {
        errors = append(errors, fmt.Sprintf("error walking %s: %v\n", dir, err))
        return
    }

		clipsMutex.Lock()
    for _, f := range files {
				if C.Load(C.CString(f)) != 1 {
            fmt.Fprintf(os.Stderr, "failed to load %s\n", f)
						errors = append(errors, fmt.Sprintf("Cannot load this file: %s", f))
            continue
        }
        fmt.Printf("loaded: %s\n", f)
    }
		clipsMutex.Unlock()
	return agentInfo, errors
}


