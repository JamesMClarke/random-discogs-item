import discogs_client
import random
import os
from dotenv import load_dotenv
import argparse

def main():
    load_dotenv()

    TOKEN = os.getenv("TOKEN")

    parser = argparse.ArgumentParser(description="Get a random item from a Discogs collection folder.")
    parser.add_argument("who", type=str, help="Who's folder to get the item from", choices=["alice","james","both"])
    parser.add_argument('-s', '--singles', action='store_true', help="Include singles")
    parser.add_argument('--notShared', action='store_false', help="Exclude shared folder")
    args = parser.parse_args()

    d = discogs_client.Client('ExampleApplication/0.1', user_token=TOKEN)
    user = d.identity()
    # print("Folders:")
    # for folder in user.collection_folders:
    #     print(f"{folder.id}: {folder.name} ({folder.count} items)")

    if args.who == "alice" or args.who == "both":
        alice_records = [folder.releases for folder in user.collection_folders if "Alice LPs" in folder.name]
        if args.singles:
            alice_singles = [folder.releases for folder in user.collection_folders if "Alice Singles" in folder.name]
    if args.who == "james" or args.who == "both":
        james_records = []
        james_records = [folder.releases for folder in user.collection_folders if "James LPs" in folder.name]
        if args.singles:
            james_singles = [folder.releases for folder in user.collection_folders if "James Singles" in folder.name]
    shared_records = [folder.releases for folder in user.collection_folders if "Shared LPs" in folder.name]

    if args.who == "alice":
        # make a concrete list so we can safely concatenate multiple PaginatedList objects
        pool = list(alice_records[0])

        if args.notShared:
            pool = list(pool) + list(shared_records[0])

        if args.singles:
            pool = list(pool) + list(alice_singles[0])
    elif args.who == "james":
        pool = list(james_records[0])
        if args.notShared:
            pool = list(pool) + list(shared_records[0])

        if args.singles:
            pool = list(pool) + list(james_singles[0])
    else:
        # both: create lists from the PaginatedList objects before combining
        pool = list(alice_records[0]) + list(james_records[0])
        if args.notShared:
            pool = list(pool) + list(shared_records[0])

        if args.singles:
            pool = list(pool) + list(alice_singles[0]) + list(james_singles[0])

    print(f"Picking a random item from {args.who}'s folder of {len(pool)} items...")
    item = random.choice(pool)
    release = d.release(item.id)
    print(f"Random item from {args.who}'s folder: {release.title} by {', '.join([artist.name for artist in release.artists])}, {release.year}")
    print(f"Tracklist: {', '.join([track.title for track in release.tracklist])}")

if __name__ == "__main__":
    main()