# Pokedex CLI

A command-line interface (CLI) application that simulates a Pokedex for exploring and catching Pokemon.

## Features

- View Pokemon locations in batches of 20 (`map` command)
- Explore specific locations to find Pokemon (`explore` command)
- Attempt to catch encountered Pokemon (`catch` command)
- View your caught Pokemon collection (`pokedex` command)
- Caching system to reduce API calls

## Usage

1. Run the program:
   ```
   go run main.go
   ```

2. Available commands:
   - `help`: Show help message
   - `exit`: Exit the Pokedex
   - `map`: Show next 20 locations
   - `mapg`: Show previous 20 locations
   - `explore <location-name>`: Explore a specific location
   - `catch <pokemon-name>`: Attempt to catch a Pokemon
   - `pokedex`: View your caught Pokemon

## API Integration

The application uses the [PokeAPI](https://pokeapi.co/) to fetch Pokemon data.

## Dependencies

- Standard Go libraries
- Custom caching package (`pokecache`)

## Project Structure

- `main.go`: Main application logic
- `pokecache/`: Custom caching implementation