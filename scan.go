package main

import (
	"database/sql"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type HostInfo struct {
	IP      string
	Puertos string
	OS      string
	Conex   string
	Fecha   string
}

func main() {
	fmt.Println(" Iniciando escaneo NAC con Nmap...")

	db, err := sql.Open("sqlite3", "./inventario.db")
	if err != nil {
		fmt.Println("Error abriendo DB:", err)
		return
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		fmt.Println("Error creando tablas:", err)
		return
	}

	ips := discoverHosts()
	for _, ip := range ips {
		fmt.Printf(" Procesando %s\n", ip)

		puertos := scanPorts(ip)
		osDetected := detectOS(ip)
		conexiones := checkActiveConnectionsLocal() // de momento lo hacemos local

		info := HostInfo{
			IP:      ip,
			Puertos: puertos,
			OS:      osDetected,
			Conex:   conexiones,
			Fecha:   time.Now().Format("2006-01-02 15:04:05"),
		}

		insertHistorico(db, info)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println(" Escaneo finalizado.")
}

func initDB(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS historico (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ip TEXT,
			fecha TEXT,
			puertos TEXT,
			os TEXT,
			conex TEXT
		);
	`)
	return err
}

func discoverHosts() []string {
	fmt.Println(" Descubriendo hosts...")

	nmapPath := `C:\Program Files (x86)\Nmap\nmap.exe`
	cmd := exec.Command(nmapPath, "-sn", "192.168.40.0/24") // Cambia tu red aquí

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error ejecutando nmap:", err)
		return nil
	}

	re := regexp.MustCompile(`Nmap scan report for ([\d\.]+)`)
	matches := re.FindAllStringSubmatch(string(out), -1)

	var ips []string
	for _, match := range matches {
		ips = append(ips, match[1])
	}
	return ips
}

func scanPorts(ip string) string {
	fmt.Printf(" Escaneando puertos en %s...\n", ip)

	nmapPath := `C:\Program Files (x86)\Nmap\nmap.exe`
	cmd := exec.Command(nmapPath, "-p-", "--min-rate", "1000", ip)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error en escaneo de puertos:", err)
		return ""
	}

	re := regexp.MustCompile(`(\d+)/tcp\s+open`)
	matches := re.FindAllStringSubmatch(string(out), -1)

	var puertos []string
	for _, match := range matches {
		puertos = append(puertos, match[1])
	}
	return strings.Join(puertos, ",")
}

func detectOS(ip string) string {
	fmt.Printf("🖥 Detectando sistema operativo en %s...\n", ip)

	nmapPath := `C:\Program Files (x86)\Nmap\nmap.exe`
	cmd := exec.Command(nmapPath, "-O", ip)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error detectando OS:", err)
		return "Desconocido"
	}

	re := regexp.MustCompile(`OS details: (.+)`)
	matches := re.FindStringSubmatch(string(out))
	if len(matches) > 1 {
		return matches[1]
	}

	return "Desconocido"
}

// Nueva función para revisar conexiones activas (simulado local por ahora)
func checkActiveConnectionsLocal() string {
	fmt.Println(" Revisando conexiones activas (simulado local)...")

	cmd := exec.Command("netstat", "-an")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error ejecutando netstat:", err)
		return "Error"
	}

	var conexiones []string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "ESTABLISHED") {
			conexiones = append(conexiones, line)
		}
	}

	// Solo devolvemos los destinos únicos como string resumido
	return strings.Join(conexiones, " | ")
}

func insertHistorico(db *sql.DB, info HostInfo) {
	_, err := db.Exec(`
		INSERT INTO historico (ip, fecha, puertos, os, conex)
		VALUES (?, ?, ?, ?, ?)`,
		info.IP, info.Fecha, info.Puertos, info.OS, info.Conex)
	if err != nil {
		fmt.Println("Error insertando histórico:", err)
	} else {
		fmt.Printf(" Insertado: %s [%s] OS: %s\n", info.IP, info.Puertos, info.OS)
	}
}
