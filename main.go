package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	// Redirect the command's output to the current process's standard output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func initGitWithPrompt(dirPath *string) {
	if dirPath == nil {
		err := runCommand("git", "init")
		if err != nil {
			fmt.Println("Error:", err)
		}
	} else {
		err := runCommand("git", "init", *dirPath)
		if err != nil {
			fmt.Println("Error:", err)
		}
		err = os.Chdir(*dirPath)
		if err != nil {
			fmt.Println("Error changing directory:", err)
		}
	}
	var branchName string
	fmt.Print("branch name: (main) ")
	fmt.Scanln(&branchName)

	currentBranch, err := getCurrentBranch()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = runCommand("git", "branch", "-m", currentBranch, branchName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var remoteURL string
	fmt.Print("remote url: (press enter to skip [setting remote]) ")
	fmt.Scanln(&remoteURL)
	fmt.Println()
	if remoteURL == "" {
		fmt.Println("gitstarter finish. enjoy your git!")
		return
	} else {
		var remoteName string
		fmt.Print("remote name: (origin) ")
		fmt.Scanln(&remoteName)
		fmt.Println()
		if remoteName == "" {
			remoteName = "origin"
			err = runCommand("git", "remote", "add", remoteName, remoteURL)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
	}

	var doAddAndPush string
	fmt.Print("add & push whole files? (y/n): (press enter to skip [add & push]) ")
	fmt.Scanln(&doAddAndPush)
	fmt.Println()
	if doAddAndPush == "y" {
		var commitMessage string
		fmt.Print("commit message: (press enter to cancel [add & push]) ")
		fmt.Scanln(&commitMessage)
		fmt.Println()
		if commitMessage != "" {
			err = runCommand("git", "add", ".")
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			err = runCommand("git", "commit", "-m", commitMessage)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			err = runCommand("git", "push", "--set-upstream", commitMessage)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
	}
}

func main() {
	// Access command-line arguments
	args := os.Args

	// Get the name of the executable (the first argument, index 0)
	executableName := filepath.Base(args[0])

	// Ensure that at least one command-line argument is provided
	if len(args) < 2 {
		fmt.Printf("Usage: %s <command>\n\n", executableName)
		fmt.Printf("%-12s%s\n", "init <directory:optional>", "Init git with prompt")
		return
	}

	// Access the second argument (index 1), which is "[somecommand]" in your example
	command := args[1]
	//fmt.Printf("You provided the command: %s\n", command)

	// You can perform actions based on the command here
	if command == "init" {
		if len(args) < 3 {
			initGitWithPrompt(nil)
			return
		}
		initGitWithPrompt(&args[2])
	}
}
