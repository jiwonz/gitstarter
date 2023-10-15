package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func deleteDir(dirPath string) {
	err := os.RemoveAll(dirPath)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Directory deleted:", dirPath)
	}
}

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

func initGitWithPrompt(dirPath *string) (int, string) {
	var workingDirPath string
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return 0, ""
	}
	if dirPath == nil {
		absolutePath, err := filepath.Abs(".git")
		if err != nil {
			fmt.Println("Error:", err)
			return 0, ""
		}
		workingDirPath = absolutePath
		err = runCommand("git", "init")
		if err != nil {
			fmt.Println("Error:", err)
			return 0, workingDirPath
		}
	} else {
		workingDirPath = filepath.Dir(currentDir)
		err := runCommand("git", "init", *dirPath)
		if err != nil {
			fmt.Println("Error:", err)
			return 0, workingDirPath
		}
		err = os.Chdir(*dirPath)
		if err != nil {
			fmt.Println("Error changing directory:", err)
			return 0, workingDirPath
		}
	}
	var branchName string
	fmt.Print("branch name: (main) ")
	fmt.Scanln(&branchName)
	if branchName == "" {
		branchName = "main"
	}

	currentBranch, err := getCurrentBranch()
	if err != nil {
		fmt.Println("Error:", err)
		return 0, workingDirPath
	}

	err = runCommand("git", "branch", "-m", currentBranch, branchName)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, workingDirPath
	}

	var remoteURL string
	fmt.Print("remote url: (press enter to skip [setting remote]) ")
	fmt.Scanln(&remoteURL)
	fmt.Println()
	if remoteURL == "" {
		return 1, workingDirPath
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
				return 0, workingDirPath
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
				return 0, workingDirPath
			}
			err = runCommand("git", "commit", "-m", commitMessage)
			if err != nil {
				fmt.Println("Error:", err)
				return 0, workingDirPath
			}
			err = runCommand("git", "push", "--set-upstream", commitMessage)
			if err != nil {
				fmt.Println("Error:", err)
				return 0, workingDirPath
			}
		}
	}
	return 1, workingDirPath
}

func main() {
	// Access command-line arguments
	args := os.Args

	// Get the name of the executable (the first argument, index 0)
	executableName := filepath.Base(args[0])

	// Ensure that at least one command-line argument is provided
	if len(args) < 2 {
		fmt.Printf("Usage: %s <command>\n\n", executableName)
		fmt.Println("init <directory:optional>    ", "Init git with prompt")
		return
	}

	// Access the second argument (index 1), which is "[somecommand]" in your example
	command := args[1]
	//fmt.Printf("You provided the command: %s\n", command)

	// You can perform actions based on the command here
	if command == "init" {
		var status int
		var deleteAfterError string
		if len(args) < 3 {
			status, deleteAfterError = initGitWithPrompt(nil)
		} else {
			status, deleteAfterError = initGitWithPrompt(&args[2])
		}
		if status == 0 {
			if deleteAfterError != "" {
				deleteDir(deleteAfterError)
			}
		} else {
			fmt.Println("gitstarter: successfully finished. enjoy your git!")
		}
	}
}
