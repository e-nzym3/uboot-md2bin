package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func main() {

	// Commandline args
	//// Help handler
	help := flag.Bool("h", false, "Display help")
	helpLong := flag.Bool("help", false, "Display help")
	flag.Parse()

	if *help || *helpLong {
		fmt.Println("Usage: uboot-md2bin <fileName>")
		fmt.Println("-- Converts standard U-Boot md command output to a binary file")
		return
	}

	//// fileName arg handler
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("[!] Error: No input file specified.")
		fmt.Println("Usage: ubootdumper <fileName>")
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
	pre_re := regexp.MustCompile("^[0-9a-zA-Z].*: ")

	//// Regex pattern for removing suffixed lines
	suf_re := regexp.MustCompile("\\s{4}.*$")

	//// Regex for empty spaces
	space_re := regexp.MustCompile(" ")

	// Creating a scanner to read input file line by line
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		line = pre_re.ReplaceAllString(line, "")
		line = suf_re.ReplaceAllString(line, "")
		line = space_re.ReplaceAllString(line, "")

		// Convert hex string to bytes (hexData)
		hexData, err := hex.DecodeString(line)
		if err != nil {
			fmt.Println("[!] Error decoding hex string:", err)
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
