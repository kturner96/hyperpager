package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Workspace struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type ActiveWs struct {
	ActiveId int `json:"id"`
}

func GetActiveWorkspace() (int, error) {
	out, err := exec.Command("hyprctl", "activeworkspace", "-j").Output()
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	var active ActiveWs


	err = json.Unmarshal(out, &active)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	return active.ActiveId, nil
}

func GetWorkspaces() ([]Workspace, error) {
	ws, err := exec.Command("hyprctl", "workspaces", "-j").Output()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	workspaces := []Workspace{}

	err = json.Unmarshal(ws, &workspaces)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return workspaces, nil
}

func Render() {

	workspaces, err := GetWorkspaces()
	if err != nil {
		fmt.Println(err)
		return
	}
	activeWs, err := GetActiveWorkspace()
	if err != nil {
		fmt.Println(err)
		return
	}

	sort.Slice(workspaces, func(f, l int) bool {
		return workspaces[f].Id < workspaces[l].Id
	})

	allWs := ""

	for _, ws := range workspaces {
		if activeWs == ws.Id {
    	allWs += fmt.Sprintf("● ")
		} else {
    		allWs += fmt.Sprintf("○ ")
}
	}

	fmt.Println(allWs)
}

func main() {
	socketPath := fmt.Sprintf("%s/hypr/%s/.socket2.sock", os.Getenv("XDG_RUNTIME_DIR"), os.Getenv("HYPRLAND_INSTANCE_SIGNATURE"))
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatal(err)
	}

	Render()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "workspace>>") {
			Render()
		}
	}

}