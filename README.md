# QuickRadio

QuickRadio streams live **NHL radio games** using **FFmpeg** for transcoding and **Beep** for audio playback. The app is built using a **multithreaded producer-consumer model**, where each game runs on its own dedicated thread. The **Model-View-Controller (MVC)** architecture ensures separation of concerns and a clean, modular structure. The responsive, data-driven UI dynamically updates in real-time, displaying live game stats, player information (on ice and penalty box), and providing seamless audio streaming through an efficient audio queue.

![Screenshot 2025-04-07 194608](https://github.com/user-attachments/assets/34afea88-2943-4245-86b4-a5fb5d6ea929)


## Features
- **Live NHL Radio Streaming:** Stream live NHL radio broadcasts with audio transcoded via **FFmpeg** for format compatibility.
- **Multithreaded Per-Game Architecture:** Each game is handled by its own dedicated thread for isolated management of game data, stats, and audio streaming.
- **Producer-Consumer Model:** Uses a producer-consumer architecture to manage separate threads for fetching, transcoding, and streaming game data/audio.
- **Dynamic Labels for Real-Time Stats:** The UI updates dynamically with real-time stats such as:
  - Game score, period, and time remaining.
  - Players currently on the ice and in the penalty box.
- **Audio Queue for Seamless Playback:** Ensures smooth, uninterrupted audio streaming using an audio queue for managing live audio samples.
- **Responsive Data-Driven UI:** Built with **Qt**, the user interface responds to game events and updates game stats and player info in real-time.
- **MVC Architecture:**
  - **Model:** Manages game data, audio stream fetching, transcoding, and processing of updates.
  - **View:** Displays the game stats, player information, and allows interaction with the audio controls.
  - **Controller:** Coordinates the interaction between the model and view, ensuring UI updates are non-blocking and seamless.
- **FFmpeg for Audio Transcoding:** Uses **FFmpeg** to transcode the live NHL radio stream into a format suitable for playback.
- **Beep for Audio Playback:** Leverages **Beep** for synchronized and smooth audio playback during streaming.

## How It Works
- **Per-Game Multithreaded Model:**  
  Each game runs on its own dedicated thread. The threads:
  - Fetch live audio data and game statistics.
  - Use **FFmpeg** to transcode the audio.
  - Stream the audio through a queue, ensuring smooth playback.
  
- **Producer-Consumer Model:**  
  The system runs separate threads for:
  - **Producer:** Fetching live data and transcoding audio.
  - **Consumer:** Managing the audio queue and updating the UI dynamically with live stats.
  
- **Responsive UI with Dynamic Labels:**  
  The user interface updates real-time game stats such as score, time, players on ice, and penalties using dynamic labels, while keeping the audio stream uninterrupted.

- **MVC Architecture:**  
  - **Model:** Fetches and processes game data, including real-time stats, players, and audio streaming.
  - **View:** Displays the game stats and player information, including dynamic labels that update in real time.
  - **Controller:** Manages communication between the model and view, ensuring smooth and non-blocking UI updates.

## Requirements
- Go 1.18+
- Qt (with Go bindings)
- FFmpeg (for transcoding the audio stream)
- Beep library for audio playback

## Setup
1. Clone the repository.
2. Install Go dependencies.
3. Set up Qt bindings for Go.
4. Install **FFmpeg** for audio stream transcoding.
5. Build and run the application to begin streaming live NHL games with real-time stats, player information, and audio playback.

## Contributing
Contributions are welcome! If you have suggestions, bug reports, or would like to improve the project, please feel free to:
- Open an issue for feature requests or bugs.
- Fork the repository, create a feature branch, and submit a pull request.
