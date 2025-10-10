import discogs_client
import random
import os
from dotenv import load_dotenv
import argparse

load_dotenv()

TOKEN = os.getenv("TOKEN")

parser = argparse.ArgumentParser(description="Get a random item from a Discogs collection folder.")
parser.add_argument("who", type=str, help="Who's folder to get the item from", choices=["alice","james","both"])
parser.add_argument('-s', '--singles', action='store_true', help="Include singles")
args = parser.parse_args()

d = discogs_client.Client('ExampleApplication/0.1', user_token=TOKEN)
user = d.identity()
print("Folders:")
for folder in user.collection_folders:
    print(f"{folder.id}: {folder.name} ({folder.count} items)")