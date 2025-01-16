package core

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type ProcessOrchestrator interface {
	StartApp(path, edition, urlParam string) (int, string, error)
	StopApp(pid int) error
	DeleteProcess(pid int) error
	GetStatusFromPID(pid int) (string, error)
	GetAllProcesses() []*Process
	createConfig(edition, urlParam, additionalFolder string) error
	CleanCash(edition, additionalFolder string) error
}

type ProcessInfo struct {
	Path             string
	AdditionalFolder string
	Edition          string // "default", "chain"
	UrlParam         string
}

type Process struct {
	PID    int
	Status string // "running", "stopped"
	Info   *ProcessInfo
}

type ProcessManager struct {
	processes map[int]*Process
	mutex     sync.Mutex
}

// NewProcessManager Конструктор менеджера процессов
func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		processes: make(map[int]*Process),
		mutex:     sync.Mutex{},
	}
}

// StartApp запуск редакции. Возвращает PID, папку кеша, ошибку
func (pm *ProcessManager) StartApp(path, edition, urlParam string) (int, string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return 0, "", fmt.Errorf("file does not exist: %s", path)
	}
	randomNumber := rand.Int()
	additionalFolder := fmt.Sprintf("iiko%d", randomNumber)

	cmd := exec.Command(path, fmt.Sprintf("/AdditionalTmpFolder=%s", additionalFolder))

	err := pm.createConfig(edition, urlParam, additionalFolder)
	if err != nil {
		return 0, "", err
	}

	err = cmd.Start()
	if err != nil {
		return 0, "", fmt.Errorf("error starting app: %v", err)
	}

	pid := cmd.Process.Pid

	pm.mutex.Lock()
	pm.processes[pid] = &Process{
		PID:    pid,
		Status: "running",
		Info: &ProcessInfo{
			Path:             path,
			AdditionalFolder: additionalFolder,
			Edition:          edition,
			UrlParam:         urlParam,
		},
	}
	pm.mutex.Unlock()

	go func() {
		err = cmd.Wait()
		if err != nil {
			return
		}

		pm.mutex.Lock()
		pm.processes[pid].Status = "stopped"
		pm.mutex.Unlock()

		pm.CleanCash(edition, additionalFolder)
	}()

	return pid, additionalFolder, nil
}

func (pm *ProcessManager) StopApp(pid int) error {
	pm.mutex.Lock()
	_, ok := pm.processes[pid]
	pm.mutex.Unlock()

	if !ok {
		return fmt.Errorf("process not found")
	}

	err := exec.Command("taskkill", "/PID", strconv.Itoa(pid)).Run()
	if err != nil {
		return fmt.Errorf("error stopping app: %v", err)
	}

	return nil
}

func (pm *ProcessManager) DeleteProcess(pid int) error {
	running, err := pm.IsProcessRunning(pid)
	if err != nil {
		return err
	}
	if running {
		return fmt.Errorf("process not stopped")
	}

	pm.mutex.Lock()
	delete(pm.processes, pid)
	pm.mutex.Unlock()

	return nil
}

func (pm *ProcessManager) GetStatusFromPID(pid int) (string, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	process, ok := pm.processes[pid]
	if !ok {
		return "", fmt.Errorf("process not found")
	}

	return process.Status, nil
}

func (pm *ProcessManager) GetAllProcesses() []*Process {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	processList := make([]*Process, 0, len(pm.processes))
	for _, process := range pm.processes {
		processList = append(processList, process)
	}
	return processList
}

func (pm *ProcessManager) IsProcessRunning(pid int) (bool, error) {
	cmd := exec.Command("tasklist", "/FI", "PID eq "+strconv.Itoa(pid))
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("checking process %d error: %v", pid, err)
	}

	return strings.Contains(string(output), strconv.Itoa(pid)), nil
}

func (pm *ProcessManager) createConfig(edition, urlParam, additionalFolder string) error {
	var cashPath string

	switch edition {
	case "default":
		cashPath = filepath.Join(os.Getenv("APPDATA"), "iiko", "Rms", additionalFolder)
	case "chain":
		cashPath = filepath.Join(os.Getenv("APPDATA"), "iiko", "Chain", additionalFolder)
	default:
		return fmt.Errorf("unknown edition: %s", edition)
	}

	configDir := filepath.Join(cashPath, "config")
	configPath := filepath.Join(configDir, "backclient.config.xml")

	// Создаем папку config, если она не существует
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err = os.MkdirAll(configDir, os.ModePerm); err != nil {
			return fmt.Errorf("error creating config directory: %v", err)
		}
	}

	// Парсим URL
	parsedURL, err := url.Parse(urlParam)
	if err != nil {
		return fmt.Errorf("error parsing URL: %v", err)
	}

	// Определяем порт
	port := parsedURL.Port()
	if port == "" {
		if parsedURL.Scheme == "https" {
			port = "443"
		} else {
			port = "8080"
		}
	}

	// Создаем структуру XML
	type Server struct {
		ServerName   string `xml:"ServerName"`
		Version      string `xml:"Version"`
		ComputerName string `xml:"ComputerName"`
		Protocol     string `xml:"Protocol"`
		ServerAddr   string `xml:"ServerAddr"`
		ServerSubUrl string `xml:"ServerSubUrl"`
		Port         string `xml:"Port"`
		IsPresent    string `xml:"IsPresent"`
	}

	type Config struct {
		XMLName xml.Name `xml:"config"`
		Servers []Server `xml:"ServersList"`
		Login   string   `xml:"Login"`
	}

	// Читаем существующий конфиг или создаем новый
	var (
		cfg  Config
		file []byte
	)
	if _, err = os.Stat(configPath); err == nil {
		file, err = os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("error reading config file: %v", err)
		}

		if err = xml.Unmarshal(file, &cfg); err != nil {
			return fmt.Errorf("error parsing config file: %v", err)
		}
	}

	// Удаляем старые ServerList
	cfg.Servers = nil

	// Добавляем новый ServerList
	newServer := Server{
		ServerName:   "",
		Version:      "",
		ComputerName: "",
		Protocol:     parsedURL.Scheme,
		ServerAddr:   parsedURL.Hostname(),
		ServerSubUrl: "/resto",
		Port:         port,
		IsPresent:    "false",
	}
	cfg.Servers = append(cfg.Servers, newServer)

	// Обновляем Login
	cfg.Login = "iikoUser"

	// Сериализуем структуру в XML
	file, err = xml.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling XML: %v", err)
	}

	// Добавляем XML-заголовок
	file = []byte(xml.Header + string(file))

	// Записываем XML в файл
	if err = os.WriteFile(configPath, file, os.ModePerm); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

func (pm *ProcessManager) CleanCash(edition, additionalFolder string) error {
	var cashPath string

	switch edition {
	case "default":
		cashPath = filepath.Join(os.Getenv("APPDATA"), "iiko", "Rms", additionalFolder)
	case "chain":
		cashPath = filepath.Join(os.Getenv("APPDATA"), "iiko", "Chain", additionalFolder)
	default:
		return fmt.Errorf("unknown edition: %s", edition)
	}

	if err := os.RemoveAll(cashPath); err != nil {
		return fmt.Errorf("error deleting %s: %v", cashPath, err)
	}

	return nil
}
