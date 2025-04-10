package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func main() {
	// Commandline args
	//// Help handler
	help := flag.Bool("h", false, "Display help")
	helpLong := flag.Bool("help", false, "Display help")
	fix := flag.Bool("fix", false, "Fix non-sequential memory addresses")
	flag.Parse()

	if *help || *helpLong {
		fmt.Println("Usage: uboot-md2bin [options] <fileName>")
		fmt.Println("Options:")
		fmt.Println("  --fix    Fix non-sequential memory addresses")
		fmt.Println("  -h, --help    Display this help message")
		return
	}

	//// fileName arg handler
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("[!] Error: No input file specified.")
		fmt.Println("Usage: uboot-md2bin [options] <fileName>")
		return
	}

	fileName := args[0]

	// Begin Execution
	fmt.Println("\n╻━━    ┳┳  ┳┓            ┓┏┓┓ •      ━━╻")
	fmt.Println("┃━━━━  ┃┃━━┣┫┏┓┏┓╋   ┏┳┓┏┫┏┛┣┓┓┏┓  ━━━━┃")
	fmt.Println("╹━━    ┗┛  ┻┛┗┛┗┛┗   ┛┗┗┗┻┗━┗┛┗┛┗    ━━╹")
	fmt.Println("   - By enzym3 (github.com/e-nzym3) -\n")

	fmt.Println("[*] Reading " + fileName + " file for processing...")
	inputFile, outputFile, err := fileHandler(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Close both files once function exits
	defer inputFile.Close()
	defer outputFile.Close()

	// Regex Handling
	//// Regex pattern for removing prefix to all data lines. Matches on strings such as "8012abcd: "
	pre_re := regexp.MustCompile("^[0-9a-zA-Z]+: ")

	//// Regex pattern for removing suffixed lines
	suf_re := regexp.MustCompile("\\s{4}.*$")

	//// Regex for empty spaces
	space_re := regexp.MustCompile(" ")

	// Creating a scanner to read input file line by line
	scanner := bufio.NewScanner(inputFile)

	// Store all lines for potential fixing
	var originalLines []string
	x := 1 // Start Line Counter
	badData := false // Defining variable used to check if bad data exists in the supplied dump

	// Regex pattern for validating the complete line format
	lineFormatRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}: [0-9a-fA-F]{2}( [0-9a-fA-F]{2}){15}    [\x20-\x7E]{16}$`)

	fmt.Println("[*] Validating contents...")
	for scanner.Scan() {
		line := scanner.Text()
		originalLines = append(originalLines, line)

		if !lineFormatRegex.MatchString(line) {
			fmt.Printf("[!] Error: Invalid line format at line %d of %s\n", x, fileName)
			fmt.Printf("    Expected format: 00000000: ff ff ff ff ff ff ff ff ff ff ff ff ff ff ff ff    ................\n")
			fmt.Printf("    Got: %s\n", line)
			badData = true
		}
		x++
	}

	// If bad data existed in the supplied dump, exit.
	if badData {
		fmt.Println("\n[!] Corrupted data detected in supplied dump, exiting... Fix dump and re-run.")
		return
	}

	// If --fix flag is set, fix the memory addresses and create a new file
	if *fix {
		fmt.Println("[*] Fixing non-sequential memory addresses...")
		fixedLines := fixMemoryAddresses(originalLines)
		
		// Create fixed file
		baseName := fileName[:len(fileName)-len(filepath.Ext(fileName))]
		fixedFileName := baseName + "_fixed" + filepath.Ext(fileName)
		fixedFile, err := os.Create(fixedFileName)
		if err != nil {
			fmt.Printf("[!] Error creating fixed file: %v\n", err)
			return
		}
		defer fixedFile.Close()

		// Write fixed lines to new file
		for _, line := range fixedLines {
			fmt.Fprintln(fixedFile, line)
		}
		fmt.Printf("[*] Created fixed version: %s\n", fixedFileName)

		// Use the fixed lines for processing
		originalLines = fixedLines
	}

	// Process the lines (either original or fixed)
	for _, line := range originalLines {
		line = pre_re.ReplaceAllString(line, "")
		line = suf_re.ReplaceAllString(line, "")
		line = space_re.ReplaceAllString(line, "")

		// Convert hex string to bytes (hexData)
		hexData, err := hex.DecodeString(line)
		if err != nil {
			fmt.Printf("[!] Error decoding hex string: %v\n", err)
			return
		}

		// Writing bytes to file
		_, err = outputFile.Write(hexData)
		if err != nil {
			fmt.Println("[!] Error writing to file:", err)
			return
		}
	}

	baseName := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	fmt.Println("[*] Processing complete! Resulting file is: " + baseName + "_output.bin\n")
}

// fixMemoryAddresses takes a slice of lines and returns a new slice with sequential memory addresses
func fixMemoryAddresses(lines []string) []string {
	var fixedLines []string
	baseAddr := uint64(0)

	// Regex pattern for memory address (8 hex digits followed by colon)
	addrRegex := regexp.MustCompile(`^([0-9a-fA-F]{8}):`)

	// Extract the first memory address to use as base
	if len(lines) > 0 {
		matches := addrRegex.FindStringSubmatch(lines[0])
		if len(matches) > 1 {
			addr, err := strconv.ParseUint(matches[1], 16, 64)
			if err == nil {
				baseAddr = addr
			}
		}
	}

	// Process each line and fix the address
	for i, line := range lines {
		// Find the memory address portion
		matches := addrRegex.FindStringSubmatch(line)
		if len(matches) < 2 {
			continue
		}

		// Calculate new address
		newAddr := baseAddr + uint64(i*16)
		newAddrStr := fmt.Sprintf("%08x", newAddr)

		// Replace the old address with the new one
		fixedLine := addrRegex.ReplaceAllString(line, newAddrStr+":")
		fixedLines = append(fixedLines, fixedLine)
	}

	return fixedLines
}

func fileHandler(fileName string) (*os.File, *os.File, error) {
	// Strip extension from supplied file
	baseName := fileName[:len(fileName)-len(filepath.Ext(fileName))]

	// Reading file
	data, err := os.Open(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("[!] Error reading file: %w", err)
	}

	// Write file
	outputFile, err := os.Create(baseName + "_output.bin")
	if err != nil {
		return nil, nil, fmt.Errorf("[!] Error writing file: %w", err)
	}

	return data, outputFile, nil
}
