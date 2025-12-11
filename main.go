package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	// Verifica che sia stato fornito un argomento
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Modalità batch se il primo argomento è -batch
	if os.Args[1] == "-batch" {
		if len(os.Args) < 4 {
			printUsage()
			os.Exit(1)
		}
		inputFile := os.Args[2]
		outputFile := os.Args[3]
		var dnsServer string
		if len(os.Args) >= 5 {
			dnsServer = os.Args[4]
		}
		batchMode(inputFile, outputFile, dnsServer)
		return
	}

	// Modalità singola IP
	singleMode()
}

func printUsage() {
	fmt.Println("Uso:")
	fmt.Println("  Modalità singola: dns-reverse <indirizzo-ip> [dns-server]")
	fmt.Println("  Modalità batch:   dns-reverse -batch <input.csv> <output.csv> [dns-server]")
	fmt.Println()
	fmt.Println("Esempi:")
	fmt.Println("  dns-reverse 8.8.8.8")
	fmt.Println("  dns-reverse 10.157.250.202 10.157.255.22")
	fmt.Println("  dns-reverse -batch ips.csv risultati.csv")
	fmt.Println("  dns-reverse -batch ips.csv risultati.csv 10.157.255.22")
}

func createResolver(dnsServer string) *net.Resolver {
	if dnsServer == "" {
		fmt.Println("Usando DNS server predefinito del sistema")
		return net.DefaultResolver
	}

	// Valida che il DNS server sia un IP valido
	dnsIP := net.ParseIP(dnsServer)
	if dnsIP == nil {
		fmt.Printf("Errore: '%s' non è un indirizzo IP valido per il DNS server\n", dnsServer)
		os.Exit(1)
	}

	fmt.Printf("Usando DNS server: %s\n", dnsServer)

	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 10,
			}
			return d.DialContext(ctx, "udp", dnsServer+":53")
		},
	}
}

func lookupIP(resolver *net.Resolver, ipAddress string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return resolver.LookupAddr(ctx, ipAddress)
}

func batchMode(inputFile, outputFile, dnsServer string) {
	// Apri il file di input
	inFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Errore nell'aprire il file di input '%s': %v\n", inputFile, err)
		os.Exit(1)
	}
	defer inFile.Close()

	// Crea il file di output
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Errore nella creazione del file di output '%s': %v\n", outputFile, err)
		os.Exit(1)
	}
	defer outFile.Close()

	resolver := createResolver(dnsServer)

	// Leggi il file di input
	scanner := bufio.NewScanner(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Scrivi l'intestazione del CSV di output
	writer.Write([]string{"IP", "Hostname", "Status"})

	fmt.Printf("\nElaborazione del file %s...\n", inputFile)
	lineNum := 0
	successCount := 0
	errorCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNum++

		// Salta righe vuote e commenti
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Valida l'IP
		ip := net.ParseIP(line)
		if ip == nil {
			fmt.Printf("  [%d] %s - IP non valido\n", lineNum, line)
			writer.Write([]string{line, "", "IP non valido"})
			errorCount++
			continue
		}

		// Esegui il lookup
		names, err := lookupIP(resolver, line)
		if err != nil {
			fmt.Printf("  [%d] %s - Errore: %v\n", lineNum, line, err)
			writer.Write([]string{line, "", fmt.Sprintf("Errore: %v", err)})
			errorCount++
			continue
		}

		if len(names) == 0 {
			fmt.Printf("  [%d] %s - Nessun hostname trovato\n", lineNum, line)
			writer.Write([]string{line, "", "Nessun hostname"})
			errorCount++
		} else {
			hostname := strings.Join(names, "; ")
			fmt.Printf("  [%d] %s -> %s\n", lineNum, line, hostname)
			writer.Write([]string{line, hostname, "OK"})
			successCount++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Errore nella lettura del file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ Completato!\n")
	fmt.Printf("  Totale IP processati: %d\n", lineNum)
	fmt.Printf("  Successo: %d\n", successCount)
	fmt.Printf("  Errori: %d\n", errorCount)
	fmt.Printf("  Risultati salvati in: %s\n", outputFile)
}

func singleMode() {
	ipAddress := os.Args[1]

	// Valida che l'input sia un IP valido
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		fmt.Printf("Errore: '%s' non è un indirizzo IP valido\n", ipAddress)
		os.Exit(1)
	}

	var dnsServer string
	if len(os.Args) >= 3 {
		dnsServer = os.Args[2]
	}

	resolver := createResolver(dnsServer)

	// Esegui il reverse DNS lookup
	fmt.Printf("Ricerca reverse DNS per %s...\n", ipAddress)

	names, err := lookupIP(resolver, ipAddress)
	if err != nil {
		fmt.Printf("Errore durante la ricerca: %v\n", err)
		os.Exit(1)
	}

	// Mostra i risultati
	if len(names) == 0 {
		fmt.Println("Nessun nome trovato per questo indirizzo IP")
	} else {
		fmt.Printf("\nNomi trovati per %s:\n", ipAddress)
		for i, name := range names {
			fmt.Printf("  %d. %s\n", i+1, name)
		}
	}
}
