## Task Manager Discord Bot
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
        <li><a href="#database">Database</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#trying-out-commands">Trying out commands</a></li>
    <li><a href="#contributing">Contributing</a></li>
  </ol>
</details>


## About The Project
This project is written for the purpose of practical acquaintance with the Golang programming language. Task Manager Bot is a Discord Bot that helps you to organize your tasks.

### Built With

This discord bot is built with:

* [Go version 1.19.4](https://golang.org/)
* [DiscordGo](https://github.com/bwmarrin/discordgo)
* [MariaDB]()

## Getting Started

To get a local copy up and running follow these simple example steps.

### Prerequisites

In order to run this project you will need the following:

* Go 1.19.4 installed
* Discord account
* Configured MariaDB database 

### Installation

1. Go to the [Discord developer portal](https://discord.com/developers)
2. Create a new application
3. Add a bot user to the application
4. Get the token for the bot
5. Clone the repository

   ```sh
   git clone https://github.com/oxanahr/discord-bot.git
   ```

6. Install dependencies

   ```sh
   go mod download all
   ```

7. Create the environment variables file `.env` in the root folder and add the following:

    ```dotenv
    DISCORD_TOKEN="Your discord token"
    SERVER_GENERAL_CHANNEL_ID="Discord server general channel id"
    DB_USER="Your database username"
    DB_PASSWORD="Your database password"
    DB_SCHEMA="Your database schema"
    DB_HOST="Your database host"
    DB_PORT="Your database port"
    ```

## Usage

To run the discord bot from root directory, execute the following command
```sh
go run main.go
```

### Trying out commands
The default prefix is `/` (slash)
 ```
    /add-task <task-name> <task-description> <priority> <assignee> <deadline>
    /assign-task <task-id> <assignee>
    /start-task <task-id>
    /complete-task <task-id>
    /my-tasks <order-by-deadline-or-priority> <deadline-current-week>
    /all-tasks <unassigned> <order> <soon>
    /comment <task-id> <comment-text>
```


### Contributing
If you want to contribute to this project, feel free to do so.