package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Sink struct {
	Name string
	Desc string
}

func listSinks() ([]Sink, error) {
	cmd := exec.Command("pactl", "list", "sinks")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var sinks []Sink
	lines := strings.Split(string(out), "\n")
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "Sink #") {
			sink := Sink{}
			for ; i < len(lines); i++ {
				if strings.HasPrefix(lines[i], "\tName:") {
					sink.Name = strings.TrimSpace(strings.Split(lines[i], ":")[1])
				}
				if strings.HasPrefix(lines[i], "\tDescription:") {
					sink.Desc = strings.TrimSpace(strings.Split(lines[i], ":")[1])
				}
				if sink.Name != "" && sink.Desc != "" {
					sinks = append(sinks, sink)
					break
				}
			}
		}
	}
	return sinks, nil
}

func combineSinks(newSinkName string, sinks []Sink) error {
	slaveSinks := ""
	for _, sink := range sinks {
		if sink.Name != newSinkName {
			slaveSinks += sink.Name + ","
		}
	}
	cmd := exec.Command("pactl", "load-module", "module-combine-sink", "sink_name="+newSinkName, "sink_properties=slaves="+slaveSinks)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	sinks, err := listSinks()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Available Sinks:\n")
	for i, sink := range sinks {
		fmt.Printf("Sink %d:\n", i+1)
		fmt.Printf("Name: %s\n", sink.Name)
		fmt.Printf("Description: %s\n\n", sink.Desc)
	}
	fmt.Printf("Enter the sink number to combine with (comma-separated): ")
	reader := bufio.NewReader(os.Stdin)
	inp, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	inp = strings.TrimSpace(inp)
	slaveSinks := make([]Sink, 0, len(sinks))
	for _, slave := range strings.Split(inp, ",") {
		sinkInd, err := strconv.Atoi(strings.TrimSpace(slave))
		if err != nil {
			panic(err)
		}
		if sinkInd > len(sinks) || sinkInd < 1 {
			panic("Invalid sink number")
		}
		slaveSinks = append(slaveSinks, sinks[sinkInd-1])
	}
	fmt.Printf("Enter the new combined sink's name: ")
	inp, err = reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	inp = strings.TrimSpace(inp)
	err = combineSinks(inp, slaveSinks)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sinks combined successfully\n")
}
