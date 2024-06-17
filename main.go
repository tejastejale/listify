package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var apiURL = "https://technical-brittaney-sitrc-bdf3a6c7.koyeb.app" // Replace with your actual API URL
var configDir = getConfigDir()

func main() {
	var rootCmd = &cobra.Command{Use: "listify"}
	rootCmd.AddCommand(loginCmd, registerCmd, listTodoCmd, createTodoCmd, updateTodoCmd, deleteTodoCmd)
	rootCmd.Execute()
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate user and receive an access token",
	Run: func(cmd *cobra.Command, args []string) {
		var email, password string
		fmt.Print("Email: ")
		fmt.Scanln(&email)
		fmt.Print("Password: ")
		fmt.Scanln(&password)
		login(email, password)
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	Run: func(cmd *cobra.Command, args []string) {
		var username, email, password string
		fmt.Print("Username: ")
		fmt.Scanln(&username)
		fmt.Print("Email: ")
		fmt.Scanln(&email)
		fmt.Print("Password: ")
		fmt.Scanln(&password)
		register(username, email, password)
	},
}

var listTodoCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a list of todo tasks",
	Run: func(cmd *cobra.Command, args []string) {
		listTodos()
	},
}

var createTodoCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new todo task",
	Run: func(cmd *cobra.Command, args []string) {
		var task string
		var done bool
		fmt.Print("Task: ")
		fmt.Scanln(&task)
		fmt.Print("Done (true/false): ")
		fmt.Scanln(&done)
		createTodo(task, done)
	},
}

var updateTodoCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing todo task",
	Run: func(cmd *cobra.Command, args []string) {
		var key int
		var task string
		var done bool
		fmt.Print("Task ID: ")
		fmt.Scanln(&key)
		fmt.Print("Task: ")
		fmt.Scanln(&task)
		fmt.Print("Done (true/false): ")
		fmt.Scanln(&done)
		updateTodo(key, task, done)
	},
}

var deleteTodoCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a todo task",
	Run: func(cmd *cobra.Command, args []string) {
		var key int
		fmt.Print("Task ID: ")
		fmt.Scanln(&key)
		deleteTodo(key)
	},
}

func login(email, password string) {
	data := map[string]string{"email": email, "password": password}
	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(apiURL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		token := result["access_token"].(string)
		saveToken(token)
		fmt.Println("Login successful.")
	} else {
		fmt.Println("Invalid credentials.")
	}
}

func register(username, email, password string) {
	data := map[string]string{"username": username, "email": email, "password": password}
	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(apiURL+"/api/auth/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erroxr:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Registration successful.")
	} else {
		fmt.Println(resp.StatusCode)
	}
}

func listTodos() {
	headers := getAuthHeaders()
	if headers == nil {
		return
	}
	req, err := http.NewRequest(http.MethodGet, apiURL+"/api/todo/", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	} else {
		fmt.Println("Failed to retrieve todo tasks.")
	}
}

func createTodo(task string, done bool) {
	data := map[string]interface{}{"task": task, "done": done}
	jsonData, _ := json.Marshal(data)
	headers := getAuthHeaders()
	if headers == nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, apiURL+"/api/todo/", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Todo task created successfully.")
	} else {
		fmt.Println("Failed to create todo task.")
	}
}

func updateTodo(key int, task string, done bool) {
	data := map[string]interface{}{"task": task, "done": done}
	jsonData, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/todo/%d", apiURL, key), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	headers := getAuthHeaders()
	if headers == nil {
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Todo task updated successfully.")
	} else {
		fmt.Println("Failed to update todo task.")
	}
}

func deleteTodo(key int) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/todo/%d", apiURL, key), nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	headers := getAuthHeaders()
	if headers == nil {
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Todo task deleted successfully.")
	} else {
		fmt.Println("Failed to delete todo task.")
	}
}

func saveToken(token string) {
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, 0755)
	}
	ioutil.WriteFile(filepath.Join(configDir, "token.txt"), []byte(token), 0644)
}

func getAuthHeaders() map[string]string {
	token, err := ioutil.ReadFile(filepath.Join(configDir, "token.txt"))
	if err != nil {
		fmt.Println("You must be logged in to perform this action.")
		return nil
	}
	return map[string]string{"Authorization": "Bearer " + strings.TrimSpace(string(token))}
}

func getConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get home directory:", err)
		os.Exit(1)
	}
	return filepath.Join(homeDir, ".listify")
}
