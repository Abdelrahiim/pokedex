package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"pokedexcli/pokecache"
	"strings"
	"time"
)

type Config struct {
	Next     string
	Previous string
	Caught   map[string]bool
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, cache *pokecache.Cache, name string) error
}

type LocationResults struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ? Kepp in mind to Decode json You have to add the tags
type LocationResponse struct {
	Count    int               `json:"count"`
	Next     string            `json:"next"`
	Previous string            `json:"previous"`
	Results  []LocationResults `json:"results"`
}

type PokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon"`
}

type LocationAreaResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

func getLocations(url string, cache *pokecache.Cache) (LocationResponse, error) {
	value, ok := cache.Get(url)
	if ok {
		var locationResponse LocationResponse
		err := json.Unmarshal(value, &locationResponse)
		if err != nil {
			return LocationResponse{}, err
		}
		return locationResponse, nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationResponse{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return LocationResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationResponse{}, err
	}

	cache.Add(url, body)

	var locationResponse LocationResponse
	err = json.Unmarshal(body, &locationResponse)
	if err != nil {
		return LocationResponse{}, fmt.Errorf("failed to decode response: %v", err)
	}
	return locationResponse, nil
}

func getLocationAreaByName(url string, cache *pokecache.Cache) (LocationAreaResponse, error) {
	value, ok := cache.Get(url)
	if ok {
		var locationAreaResponse LocationAreaResponse
		err := json.Unmarshal(value, &locationAreaResponse)
		if err != nil {
			return LocationAreaResponse{}, fmt.Errorf("failed to unmarshal cached data: %v", err)
		}
		return locationAreaResponse, nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("failed to create request: %v", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	cache.Add(url, body)

	var locationAreaResponse LocationAreaResponse
	err = json.Unmarshal(body, &locationAreaResponse)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("failed to decode API response: %v", err)
	}
	return locationAreaResponse, nil
}

type Ability struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type AbilityInfo struct {
	Ability  Ability `json:"ability"`
	IsHidden bool    `json:"is_hidden"`
	Slot     int     `json:"slot"`
}

type Cries struct {
	Latest string `json:"latest"`
	Legacy string `json:"legacy"`
}

type Form struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Pokemon struct {
	Abilities      []AbilityInfo `json:"abilities"`
	BaseExperience int           `json:"base_experience"`
	Cries          Cries         `json:"cries"`
	Forms          []Form        `json:"forms"`
}

func getPokemon(url string, cache *pokecache.Cache) (Pokemon, error) {
	value, ok := cache.Get(url)
	if ok {
		var pokemon Pokemon
		err := json.Unmarshal(value, &pokemon)
		if err != nil {
			return Pokemon{}, fmt.Errorf("failed to unmarshal cached data: %v", err)
		}
		return pokemon, nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, fmt.Errorf("failed to create request: %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return Pokemon{}, fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Pokemon{}, fmt.Errorf("API returned non-200 status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, fmt.Errorf("failed to read response body: %v", err)
	}

	cache.Add(url, body)

	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return Pokemon{}, fmt.Errorf("failed to decode API response: %v", err)
	}

	return pokemon, nil
}

func commandExit(config *Config, cache *pokecache.Cache, name string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// ------------------------------------------------------------------
// ------------------------- COMMANDS -------------------------------
// ------------------------------------------------------------------
func commandHelp(config *Config, cache *pokecache.Cache, name string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Show next 20 locations")
	fmt.Println("mapg: Show previous 20 locations")
	fmt.Println("explore <location-name>: Explore a specific location")
	fmt.Println("catch <pokemon-name>: Attempt to catch a Pokemon")
	fmt.Println("pokedex: View your caught Pokemon")
	return nil
}
func commandMap(config *Config, cache *pokecache.Cache, name string) error {
	if config.Next == "" {
		return fmt.Errorf("no next URL provided")
	}

	locationResponse, err := getLocations(config.Next, cache)
	if err != nil {
		return err
	}
	config.Next = locationResponse.Next
	config.Previous = locationResponse.Previous
	for _, location := range locationResponse.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapg(config *Config, cache *pokecache.Cache, name string) error {
	if config.Previous == "" {
		return fmt.Errorf("no next URL provided")
	}

	locationResponse, err := getLocations(config.Previous, cache)
	if err != nil {
		return err
	}
	config.Next = locationResponse.Next
	config.Previous = locationResponse.Previous
	for _, location := range locationResponse.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(config *Config, cache *pokecache.Cache, name string) error {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", name)

	data, err := getLocationAreaByName(url, cache)
	if err != nil {
		return err
	}
	fmt.Println("Exploring " + name + "...")
	fmt.Println("Found Pokemon: ")
	for _, value := range data.PokemonEncounters {
		fmt.Println(value.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *Config, cache *pokecache.Cache, name string) error {
	fmt.Println("Throwing a Pokeball at " + name + " ...")
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	data, err := getPokemon(url, cache)
	if err != nil {
		return err
	}

	catchChance := 100 - data.BaseExperience/4
	catchChance = max(catchChance, 10)

	rand := rand.Intn(100)
	if rand >= catchChance {
		fmt.Println(name + " escaped!")
		return nil
	}

	fmt.Println(name + " was caught!")

	if config.Caught == nil {
		config.Caught = make(map[string]bool)
	}
	config.Caught[name] = true
	return nil
}

func commandPokedex(config *Config, cache *pokecache.Cache, name string) error {
	fmt.Println("Your Pokedex:")
	for key, _ := range config.Caught {
		pokemon := fmt.Sprintf("\t - %s", key)
		fmt.Println(pokemon)
	}
	return nil
}

func main() {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Show help",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Show the next 20 locations",
			callback:    commandMap,
		},
		"mapg": {
			name:        "mapg",
			description: "Show the previous 20 locations",
			callback:    commandMapg,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon",
			callback:    commandCatch,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Show the pokedex",
			callback:    commandPokedex,
		},
	}
	// fmt.Println(fetchTasks("https://api.boot.dev/v1/courses_rest_api/learn-http/issues", "Low"))
	message := bufio.NewScanner(os.Stdin)
	config := Config{
		Next:     "https://pokeapi.co/api/v2/location-area",
		Previous: "",
	}
	cache := pokecache.NewCache(time.Minute)
	for {
		fmt.Print("Pokedex > ")
		message.Scan()
		input := message.Text()
		cleanedInput := cleanInput(input)

		command, ok := commands[cleanedInput[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if len(cleanedInput) == 0 {
			fmt.Println("Please provide a command")
			continue
		}
		name := ""
		if len(cleanedInput) > 1 {
			name = cleanedInput[1]
		}
		if (command.name == "explore" || command.name == "catch") && name == "" {
			fmt.Printf("Error: %s command requires a name parameter\nUsage: %s [name]\n", command.name, command.name)
			continue
		}

		err := command.callback(&config, &cache, name)
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println("You may now inspect it with the inspect command.")

	}
}

func cleanInput(input string) []string {
	// ?
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)
	cleanedInput := strings.Fields(input)

	return cleanedInput
}
