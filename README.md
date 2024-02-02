# BilgeninIsteitinBot

![BilgeninIsteitinBot Overview](./assets/Bot.jpg)

Dive into the world of BilgeninIsteitinBot, a multifaceted Discord bot engineered to infuse your server with a splash of fun, utility, and AI-powered savvy. From delivering precise weather forecasts to breaking down language barriers, setting up reminders, orchestrating games, and facilitating community polls, BilgeninIsteitinBot stands as your all-in-one Discord companion.

## 🌈 Features

### 🌤 Live Weather Updates

"Is it a good day for a beach picnic?" Find out with live weather updates for any corner of the globe.

![Weather Feature Screenshot](./assets/Weather.jpg)

### 🗣 Language Translation

Cross the language divide. Translate text instantly and make everyone in your community feel at home.

![Translation Feature Screenshot](./assets/Translation.jpg)

### ⏰ Reminders

Forget forgetting. Set custom reminders for tasks and events, powered by asynchronous processing to ensure your server's performance remains uninterrupted.

![Reminders Feature Screenshot](./assets/Reminder.jpg)

### 🪨📄✂️ Rock-Paper-Scissors Game

Challenge your friends to the timeless game of rock-paper-scissors and let the bot be your unbiased referee.

![RPS Game Feature Screenshot](./assets/RPS.jpg)

### 📊 Polls and Voting

Make collective decisions with ease. Create custom polls and let your community's voice be heard.

![Polls Feature Screenshot](./assets/Poll.jpg)

### 💡 Creative GPT-3.5 Interactions

Curious about what AI thinks? Engage with GPT-3.5 for creative and insightful responses to your burning questions.

![GPT Feature Screenshot](./assets/GPT.jpg)

### 🆘 Helper

Need a quick command reference? The `!help` command lays out all the functionalities BilgeninIsteitinBot has to offer.

![Helper Feature Screenshot](./assets/Helper.jpg)

## 🚀 Usage

Maximize BilgeninIsteitinBot's potential with these commands:

- `!weather <city name>` - Stay ahead of the weather.
- `!translate <text>` - Bridge the language gap.
- `!remind <time> <message>` - Keep tabs on important reminders.
- `!rps` - Settle debates with a game of rock-paper-scissors.
- `!poll "<question>" | <option1>, <option2>, ...` - Gauge community opinion.
- `!vote <option number>` - Cast your vote.
- `!result` - See what the community thinks.
- `!help` - Discover what BilgeninIsteitinBot can do for you.

## 💾 Data Persistence

Leveraging **MongoDB**, BilgeninIsteitinBot stores polls, game outcomes, and reminders efficiently. This setup enables quick data retrieval and manipulation, ensuring that asynchronous tasks like reminders fire off without a hitch and game states persist accurately over time.

## 🚀 Setup Instructions

Bringing `BilgeninIsteitinBot` to life in your Discord server involves a few simple steps. Let's embark on this journey together!

### 🛠 Prerequisites

Before we start, you'll need a few things:

- A **Discord account** and **administrative access** to a server.
- An **OpenWeatherMap API key** for fetching weather updates.
- An **OpenAI API key** for engaging GPT-3.5 interactions.
- **MongoDB** for storing the treasure trove of data like polls and reminders.

### 📦 Installation

#### 1. Clone the Treasure Map

Open your command line and run:

```bash
git clone https://github.com/yourusername/BilgeninIsteitinBot.git
cd BilgeninIsteitinBot
```

This magical incantation will clone the repository to your local machine and navigate you into the project directory.

#### 2. Gather Your Tools

Make sure you have [Go](https://golang.org/dl/) installed on your ship. Then, install the necessary dependencies to ensure your bot is well-equipped.

#### 3. Secrets of the Bot

Navigate to the secret chamber (`.env` file or your environment variables) and declare:

Replace the placeholders with your actual secrets. These are the keys to unlock the full potential of `BilgeninIsteitinBot`.

```bash
DISCORD_BOT_TOKEN=your_discord_bot_token_here
OPENWEATHERMAP_API_KEY=your_openweathermap_api_key_here
OPENAI_API_KEY=your_openai_api_key_here
MONGODB_URI=your_mongodb_uri_here
```

### 🏃‍♂️ Embark on the Journey

With your bot fully equipped and ready to explore, it's time to bring it to life.

### 📬 Invite the Bot to Your Server

1. Navigate to the [Discord Developer Portal](https://discord.com/developers/applications).
2. Find your bot and visit the **OAuth2** page.
3. In **Scopes**, select `bot`. In **Bot Permissions**, choose the permissions your bot needs.
4. Use the generated URL to invite your bot to your Discord server.

And that's it! `BilgeninIsteitinBot` is now ready to sail the digital seas of your Discord server, bringing weather updates, translations, reminders, games, and AI-powered conversations to your community.

## 📜 Acknowledgements

- **OpenWeatherMap API**: For empowering us with weather insights.
- **OpenAI API**: For enabling intelligent, AI-driven interactions.
- **DiscordGo**: For the robust library that makes bot development a breeze.
- **MongoDB**: For the flexible, scalable database that backs our bot's memory.

---

Embrace the full potential of your Discord server with **BilgeninIsteitinBot** – your digital genie for community engagement and fun!
