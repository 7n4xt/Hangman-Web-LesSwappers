# 🏎️HangEngine: A Car-Themed Hangman Web Application

Welcome to **Hang Engine**, a car-themed web-based version of the classic Hangman game. This is a small university project developed using Go, HTML, and CSS.

## Installation☑️

To get started, ensure that Go is installed on your machine. You can download and install Go by following this link: [https://go.dev/doc/install](https://go.dev/doc/install).

Next, clone the repository:

```
git clone https://github.com/7n4xt/Hangman-Web-LesSwappers.git
```

Open the terminal in the cloned repository and run the following command to start the application:

```
go run main.go
```

After running the command, open your web browser and navigate to [http://localhost:8080/index](http://localhost:8080/index) to access the game.

## Application Routes📶

### Aesthetic Routes (Visible to Players)

1. **Index (`/index`):**
   - The home page and main entry point of the game
   - Provides the rules of the game and a button to start a new game session

2. **Engine (`/engine`):**
   - Serves as a mini-encyclopedia for car engine information
   - Allows players to explore and learn about different engine components

3. **Scoreboard (`/scoreboard`):**
   - Displays the leaderboard and player rankings

4. **Game (`/game`):**
   - The main game interface where players guess the car-related words

5. **Result (`/result`):**
   - Shows the game results and provides options to play again or view the scoreboard

6. **Choose (`/choose`):**
   - Allows players to enter their name and select the game difficulty

### Backend Routes (Data Initialization)

1. **Start-Game (`/start-game`):**
   - Saves the player's name and difficulty level from the `/choose` route

2. **Guess (`/guess`):**
   - Handles the core game logic, processing player guesses and updating the game state

## Development Team 👨‍💻

This project was a collaboration between two university students:

- **ESUGHI Abdulmalek:**
  - Responsible for the frontend development, including CSS and HTML
  - Contributed to the Go-based HTML integration

- **ROSSI Enzo:**
  - Integrated the existing Hangman project into the current application
  - Implemented the core game logic and backend functionality

## Contributing

Contributions to improve the Hang Engine project are welcome. Please feel free to submit issues and pull requests.
