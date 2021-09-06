# DiscordMediaDownloader
Downloads all .png, .jpg, .jpeg, .gif, .webm, and .mp4 files from a discord channel
Great for archiving memes

## Media downloader script usage
- Use [DiscordChatExporter](https://github.com/Tyrrrz/DiscordChatExporter) to export a .json file of the desired text channel(s).
- Move the .json file into the same directory as main.go and rename the .json file "messages.json".
- Run the go command with "go run main.go" in the directory.
- All files will be saved to the downloads folder found in the repository directory. Duplicate files will have _1, _2, _3, etc added to the name.

## Bot usage (optional)
- Use the discord developer portal to create a new application, and then add a bot to it.
- Generate a OAuth2 URL in the developer portal with permissions to read messages, manage messages, view channels, and view message history.
- Make sure you have discord.py installed and then run the bot with "python bot.py"
- In a given channel, type "-remove " followed by a list of message ids seperated by any sort of character.
