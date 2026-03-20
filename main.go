package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func main() {
	header()

	if !ensureBraveIsClosed() {
		fmt.Println("Operacao cancelada. O arquivo Preferences nao pode ser editado com o Brave aberto.")
		return
	}

	path, err := getPrefsPath()
	if err != nil {
		fmt.Printf("Erro: %v\n", err)
		return
	}

	fmt.Println("1. Aplicar Otimizacoes (Escolher o que configurar)")
	fmt.Println("2. Restaurar ao Padrao (Usar backup .bak)")
	fmt.Println("3. Sair")
	fmt.Print("\nEscolha uma opcao: ")

	var op int
	fmt.Scanln(&op)

	switch op {
	case 1:
		runOptimizer(path)
		fmt.Print("\nDeseja restaurar ao padrao agora? (s/n): ")
		var res string
		fmt.Scanln(&res)
		if strings.ToLower(res) == "s" {
			restoreBackup(path)
		} else {
			fmt.Println("Finalizado. Saindo...")
		}
	case 2:
		restoreBackup(path)
	default:
		fmt.Println("Saindo...")
	}
}

func runOptimizer(path string) {
	selectedTweaks := make(map[string]any)

	fmt.Println("\n--- Selecione as configuracoes (s/n) ---")

	if ask("Desativar Brave VPN (Botao e servicos)?") {
		selectedTweaks["brave.vpn.enabled"] = false
	}

	if ask("Habilitar Aceleracao de Hardware?") {
		selectedTweaks["hardware_acceleration_mode_previous"] = true
	}

	if ask("Configurar idioma para Portugues (Brasil)?") {
		selectedTweaks["intl.accept_languages"] = "pt-BR,pt"
		selectedTweaks["spellcheck.dictionaries"] = []any{"pt-BR"}
		selectedTweaks["translate.enabled"] = true
	}

	if ask("Remover Recompensas (Brave Rewards)?") {
		selectedTweaks["brave.rewards.enabled"] = false
		selectedTweaks["brave.rewards.hide_button"] = true
	}

	if ask("Remover Brave News, Talk e Assistente Leo (IA)?") {
		selectedTweaks["brave.brave_news.enabled"] = false
		selectedTweaks["brave.talk.enabled"] = false
		selectedTweaks["brave.ai_chat.enabled"] = false
	}

	if ask("Desativar Telemetria e Diagnosticos?") {
		selectedTweaks["metrics.reporting_enabled"] = false
	}

	if ask("Desativar IPFS e Carteira (Wallet)?") {
		selectedTweaks["brave.ipfs.enabled"] = false
		selectedTweaks["brave.wallet.enabled"] = false
	}


	selectedTweaks["session.restore_on_startup"] = 5

	fmt.Println("\nAplicando alteracoes no perfil...")

	backupFile(path)

	data, _ := os.ReadFile(path)
	var prefs map[string]any
	json.Unmarshal(data, &prefs)

	for k, v := range selectedTweaks {
		applyTweak(prefs, k, v)
	}

	newData, _ := json.MarshalIndent(prefs, "", "  ")
	os.WriteFile(path, newData, 0644)

	fmt.Println("\nCONCLUIDO! Voce ja pode abrir o Brave.")
}

func ensureBraveIsClosed() bool {
	if isProcessRunning("brave") {
		fmt.Print("O Brave esta aberto. Deseja fecha-lo automaticamente para aplicar as mudancas? (s/n): ")
		var res string
		fmt.Scanln(&res)

		if strings.ToLower(res) == "s" {
			killProcess("brave")
			fmt.Print("Aguardando fechamento total")
			for range 5 {
				time.Sleep(1 * time.Second)
				fmt.Print(".")
				if !isProcessRunning("brave") {
					fmt.Println("\nNavegador fechado.")
					return true
				}
			}
		} else {
			return false
		}
	}
	return true
}

func killProcess(name string) {
	if runtime.GOOS == "windows" {
		exec.Command("taskkill", "/F", "/IM", name+"*", "/T").Run()
	} else {
		exec.Command("pkill", "-9", "-f", name).Run()
	}
}

func isProcessRunning(name string) bool {
	var out []byte
	if runtime.GOOS == "windows" {
		out, _ = exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s.exe", name)).Output()
	} else {
		out, _ = exec.Command("pgrep", "-f", name).Output()
	}
	return strings.Contains(strings.ToLower(string(out)), strings.ToLower(name))
}

func restoreBackup(path string) {
	bakPath := path + ".bak"
	if _, err := os.Stat(bakPath); os.IsNotExist(err) {
		fmt.Println("Arquivo de backup (.bak) nao encontrado.")
		return
	}

	os.Remove(path)
	err := os.Rename(bakPath, path)
	if err != nil {
		fmt.Printf("Erro ao restaurar: %v\n", err)
		return
	}
	fmt.Println("Restaurado com sucesso!")
}

func ask(question string) bool {
	fmt.Printf("%s (s/n): ", question)
	var res string
	fmt.Scanln(&res)
	return strings.ToLower(res) == "s"
}

func header() {
	fmt.Println("--------------------------------------")
	fmt.Println("    BRAVE DEBLOATER & SPEEDUP")
	fmt.Println("--------------------------------------")
}
func getPrefsPath() (string, error) {
	home, _ := os.UserHomeDir()
	var paths []string

	if runtime.GOOS == "windows" {
		paths = append(paths, filepath.Join(home, "AppData/Local/BraveSoftware/Brave-Browser/User Data/Default/Preferences"))
	} else {
		paths = append(paths,
			filepath.Join(home, ".config/BraveSoftware/Brave-Browser/Default/Preferences"),                           
			filepath.Join(home, ".var/app/com.brave.Browser/config/BraveSoftware/Brave-Browser/Default/Preferences"), 
			filepath.Join(home, "snap/brave/current/.config/BraveSoftware/Brave-Browser/Default/Preferences"),        
		)
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("perfil do Brave nao encontrado. Verifique se o navegador esta instalado")
}
func backupFile(src string) {
	dst := src + ".bak"
	if _, err := os.Stat(dst); err == nil {
		return
	}
	source, _ := os.Open(src)
	defer source.Close()
	destination, _ := os.Create(dst)
	defer destination.Close()
	io.Copy(destination, source)
	fmt.Println("Backup de seguranca criado (.bak)")
}

func applyTweak(m map[string]any, path string, value any) {
	parts := strings.Split(path, ".")
	curr := m
	for i, part := range parts {
		if i == len(parts)-1 {
			curr[part] = value
			fmt.Printf("  [OK] %s\n", path)
			return
		}
		if next, ok := curr[part].(map[string]any); ok {
			curr = next
		} else {
			newMap := make(map[string]any)
			curr[part] = newMap
			curr = newMap
		}
	}
}
