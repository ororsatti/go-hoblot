package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ororsatti/go-searchdex/search"
)

func clearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {
	// Sample documents
	docs := []search.Document{
		{Id: "doc1", Content: "The first document contains information about artificial intelligence and machine learning."},
		{Id: "doc2", Content: "Quantum computing explores new ways to process information using quantum-mechanical phenomena."},
		{Id: "doc3", Content: "Blockchain technology provides a decentralized and secure ledger for transactions."},
		{Id: "doc4", Content: "The history of the internet is a fascinating journey from ARPANET to the World Wide Web."},
		{Id: "doc5", Content: "Renewable energy sources like solar and wind power are crucial for a sustainable future."},
		{Id: "doc6", Content: "Genetic engineering allows for precise modification of an organism's DNA."},
		{Id: "doc7", Content: "The human brain is an incredibly complex organ responsible for thought, emotion, and memory."},
		{Id: "doc8", Content: "Space exploration continues to push the boundaries of human knowledge and discovery."},
		{Id: "doc9", Content: "Cybersecurity is essential to protect digital systems and data from theft and damage."},
		{Id: "doc10", Content: "The development of vaccines has revolutionized public health and eradicated many diseases."},
		{Id: "doc11", Content: "Nanotechnology deals with materials and devices on an atomic and molecular scale."},
		{Id: "doc12", Content: "Virtual reality and augmented reality are transforming how we interact with digital content."},
		{Id: "doc13", Content: "The study of exoplanets helps us understand the potential for life beyond Earth."},
		{Id: "doc14", Content: "Robotics and automation are increasingly integrated into manufacturing and daily life."},
		{Id: "doc15", Content: "Climate change is a global challenge requiring urgent action and international cooperation."},
		{Id: "doc16", Content: "The evolution of programming languages reflects the changing needs of software development."},
		{Id: "doc17", Content: "Neuroscience investigates the structure and function of the nervous system."},
		{Id: "doc18", Content: "Big data analytics extracts valuable insights from vast and complex datasets."},
		{Id: "doc19", Content: "The principles of astrophysics explain the phenomena of the universe, from stars to galaxies."},
		{Id: "doc20", Content: "Sustainable agriculture practices aim to produce food in an environmentally friendly way."},
	}

	// Create a new search index
	index := search.New(docs)

	// Interactive search prompt
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Search Index Ready. Enter a query or type 'exit' to quit.")

	for {
		fmt.Print("> ")
		query, _ := reader.ReadString('\n')
		query = strings.TrimSpace(query)

		if query == "exit" {
			break
		}

		results := index.Search(query, 2)

		if len(results) == 0 {
			fmt.Println("No results found.")
		} else {
			fmt.Println("Search results:")
			for _, docID := range results {
				for _, doc := range docs {
					if doc.Id == docID {
						fmt.Printf("- %s: %s\n", doc.Id, doc.Content)
						break
					}
				}
			}
		}
		fmt.Println("Click enter to reset")
		reader.ReadString('\n') // Wait for user to press Enter
		clearScreen()
	}
}
