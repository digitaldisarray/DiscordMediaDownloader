import discord
import asyncio

discord_api_key = "INSERT YOUR BOT API KEY HERE"

client = discord.Client()

@client.event
async def on_ready():
    print('Username: ' + client.user.name)
    print('User ID: ' + str(client.user.id))
    print('=====================')

@client.event
async def on_message(message):
    if message.content.lower().startswith('-remove '):
        ids = message.content[8:] # Get characters after "-remove "
        
        async for msg in message.channel.history(limit = 100): # Only will check last 100 messages
            print("Checking ID: " + str(msg.id))
            if str(msg.id) in ids:
                await msg.delete()
                print('Removed Message: ' + str(msg.id))
        
        print('Done!')
        
        await message.delete()
    
client.run(discord_api_key)
