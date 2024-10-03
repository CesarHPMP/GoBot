# GoBot

A Go-based bot that interacts with the Spotify API and Discord.

## Table of Contents

* [Introduction](#introduction)
* [Features](#features)
* [Requirements](#requirements)
* [Installation](#installation)
* [Usage](#usage)
* [API Documentation](#api-documentation)
* [Contributing](#contributing)
* [License](#license)

## Introduction

GoBot is a Go-based bot that interacts with the Spotify API and Discord. It allows users to authenticate with Spotify and retrieve their top tracks.

## Features

* Authenticates with Spotify using OAuth2
* Retrieves top tracks from Spotify
* Interacts with Discord using the Discord API

## Requirements

* Go 1.22.2 or later
* Spotify API credentials (client ID and client secret)
* Discord API credentials (bot token)

## Installation

1. Clone the repository: `git clone https://github.com/CesarHPMP/GoBot.git`
2. Install dependencies: `go mod tidy && go mod download`
3. Build the project: Make a build dirctory and run `go build ../` 

## Usage

1. Set environment variables for Spotify API credentials: `export SPOTIFY_CLIENT_ID=your_client_id` and `export SPOTIFY_CLIENT_SECRET=your_client_secret`
2. Set environment variable for Discord API credentials: `export DISCORD_BOT_TOKEN=your_bot_token`
3. Run the project by running the executable generated in the build directory.
4. Follow the prompts to authenticate with Spotify and retrieve your top tracks

## API Documentation

API documentation is available at [https://pkg.go.dev/github.com/zmb3/spotify](https://pkg.go.dev/github.com/zmb3/spotify) and [github.com/bwmarrin/discordgo](github.com/bwmarrin/discordgo)

## Contributing

Contributions are welcome! Please submit a pull request with your changes.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
