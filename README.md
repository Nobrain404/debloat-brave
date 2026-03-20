# Brave Debloater & Speedup

Script em Go para otimização do navegador Brave através da edição do arquivo de preferências. O foco é remover componentes desnecessários, desativar telemetria e melhorar a performance.

> **Nota:** Testado apenas no **Fedora** com Brave instalado via **Flatpak**.

## Funcionalidades

- Desativa Brave VPN, Rewards e News.
- Remove Assistente Leo (IA) e Brave Talk.
- Desativa Telemetria e coleta de dados.
- Habilita Aceleração de Hardware.
- Configura o idioma para Português (Brasil).
- Desativa módulos de Wallet e IPFS para economizar RAM.
- Cria backup automático (`Preferences.bak`) antes de qualquer alteração.

## Requisitos

- Go (Golang) instalado para rodar ou compilar.
- O navegador Brave deve estar fechado durante a execução.

## Instalação e Uso

### Rodar sem compilar

```Bash
# Executa o script diretamente
go run main.go


```
### Como Compilar (Linux)
``` Bash
# Gera o binário com o nome do projeto
go build -o debloat-brave main.go
chmod +x debloat-brave
./debloat-brave

```

#### Como Compilar (Windows)
``` PowerShell
# Gera o executável para Windows
go build -o debloat-brave.exe main.go
.\debloat-brave.exe

``` 
