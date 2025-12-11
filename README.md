# DNS Reverse Lookup

Programma Go per eseguire reverse DNS lookup (trovare il nome di dominio partendo da un indirizzo IP).

## Installazione

```bash
go build
```

## Uso

Il programma supporta due modalità:

### Modalità Singola

Interroga un singolo indirizzo IP:

```bash
dns-reverse.exe <indirizzo-ip> [dns-server]
```

### Modalità Batch

Processa un file CSV con un elenco di indirizzi IP:

```bash
dns-reverse.exe -batch <input.csv> <output.csv> [dns-server]
```

### Parametri

- `<indirizzo-ip>`: L'indirizzo IP da interrogare (modalità singola)
- `<input.csv>`: File CSV di input con lista di indirizzi IP (uno per riga)
- `<output.csv>`: File CSV di output con risultati (IP, Hostname, Status)
- `[dns-server]`: Il server DNS da usare per la query (opzionale)
  - Se non specificato, usa il DNS predefinito del sistema

## Esempi

### Modalità Singola

```bash
# Lookup usando il DNS predefinito del sistema
dns-reverse.exe 8.8.8.8

# Lookup usando un DNS specifico
dns-reverse.exe 192.168.1.1 8.8.4.4
```

### Modalità Batch

```bash
# Processa un file CSV usando il DNS predefinito
dns-reverse.exe -batch ips.csv risultati.csv

# Processa un file CSV usando un DNS specifico
dns-reverse.exe -batch ips.csv risultati.csv 8.8.4.4
```

#### Formato File Input (ips.csv)

Il file di input deve contenere un indirizzo IP per riga:

```
8.8.8.8
1.1.1.1
192.168.1.1
```

Puoi anche inserire commenti (righe che iniziano con `#`) e righe vuote che saranno ignorate:

```
# Server pubblici
8.8.8.8
1.1.1.1

# Server locali
192.168.1.1
```

#### Formato File Output (risultati.csv)

Il file di output è un CSV con tre colonne:

```csv
IP,Hostname,Status
8.8.8.8,dns.google.,OK
1.1.1.1,one.one.one.one.,OK
192.168.1.1,router.local.,OK
```

## Come funziona

Il programma:
1. Prende un indirizzo IP come argomento
2. Valida che l'IP sia formattato correttamente
3. Esegue una query DNS inversa usando `net.LookupAddr()`
4. Mostra tutti i nomi di dominio associati all'IP

## Requisiti

- Go 1.21 o superiore
