package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/spf13/cobra"
)

var (
	ShowAll     bool
	ShowID      bool
	ShowImage   bool
	ShowCmd     bool
	ShowCreated bool

	purple    = lipgloss.Color("99")
	gray      = lipgloss.Color("245")
	lightGray = lipgloss.Color("241")

	HeaderStyle = lipgloss.NewStyle().
			Foreground(purple).
			Bold(true)

	CellStyle = lipgloss.NewStyle().
			Padding(0, 1)

	OddRowStyle  = CellStyle.Foreground(gray)
	EvenRowStyle = CellStyle.Foreground(lightGray)
)

type container struct {
	ID      string
	Image   string
	Command string
	Created string
	Status  string
	Ports   string
	Names   string
}

func RunPs(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	containers, err := getContainers(ctx, ShowAll)
	if err != nil {
		return err
	}

	if len(containers) == 0 {
		cmd.Println("No containers found")
		return nil
	}

	renderTable(containers)
	return nil
}

func getContainers(ctx context.Context, all bool) ([]container, error) {
	dockerArgs := []string{"ps", "--format", "{{.ID}}|{{.Image}}|{{.Command}}|{{.CreatedAt}}|{{.Status}}|{{.Ports}}|{{.Names}}"}
	if all {
		dockerArgs = []string{"ps", "-a", "--format", "{{.ID}}|{{.Image}}|{{.Command}}|{{.CreatedAt}}|{{.Status}}|{{.Ports}}|{{.Names}}"}
	}

	c := exec.CommandContext(ctx, "docker", dockerArgs...)
	output, err := c.Output()
	if err != nil {
		return nil, err
	}

	var containers []container
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 7 {
			containers = append(containers, container{
				ID:      parts[0],
				Image:   parts[1],
				Command: parts[2],
				Created: parts[3],
				Status:  parts[4],
				Ports:   parts[5],
				Names:   parts[6],
			})
		}
	}

	return containers, nil
}

func renderTable(containers []container) {
	width := 80

	if envWidth := os.Getenv("COLUMNS"); envWidth != "" {
		var w int
		if _, err := fmt.Sscanf(envWidth, "%d", &w); err == nil && w > 0 {
			width = w
		}
	}

	headers := []string{"STATUS", "PORTS", "NAMES"}
	rows := make([][]string, len(containers))

	if ShowID {
		headers = append([]string{"CONTAINER ID"}, headers...)
	}
	if ShowImage {
		headers = append(headers, "IMAGE")
	}
	if ShowCmd {
		headers = append(headers, "COMMAND")
	}
	if ShowCreated {
		headers = append(headers, "CREATED")
	}

	for i, c := range containers {
		var row []string

		if ShowID {
			row = append(row, c.ID[:12])
		}
		row = append(row, c.Status)
		if ShowImage {
			row = append(row, Truncate(c.Image, width/4))
		}
		row = append(row, FormatPorts(c.Ports))
		row = append(row, c.Names)
		if ShowCmd {
			row = append(row, c.Command)
		}
		if ShowCreated {
			row = append(row, c.Created)
		}

		rows[i] = row
	}

	t := table.New().
		Border(lipgloss.Border{}).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return HeaderStyle
			case row%2 == 0:
				return EvenRowStyle
			default:
				return OddRowStyle
			}
		}).
		Headers(headers...).
		Rows(rows...)

	_, _ = lipgloss.Println(t)
}

func FormatPorts(ports string) string {
	if ports == "" {
		return ""
	}
	portList := strings.Split(ports, ", ")
	return strings.Join(portList, "\n")
}

func Truncate(s string, maxLen int) string {
	if maxLen <= 0 || len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
