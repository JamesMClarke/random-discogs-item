import discogs_client
import random
import os
from dotenv import load_dotenv
import argparse
import requests_cache
import sys
from pathlib import Path

APP_NAME = "discogs_random_item"

def default_cache_dir():
    """Return a platform-appropriate cache directory for the app (Path).

    Respects environment override DISCOGS_CACHE_DIR if set.
    """
    env_override = os.environ.get("DISCOGS_CACHE_DIR")
    if env_override:
        return Path(env_override).expanduser()

    if sys.platform == "darwin":
        return Path.home() / "Library" / "Caches" / APP_NAME
    if os.name == "nt":
        local = os.environ.get("LOCALAPPDATA") or (Path.home() / "AppData" / "Local")
        return Path(local) / APP_NAME / "Cache"
    # default: XDG spec or fallback to ~/.cache
    xdg = os.environ.get("XDG_CACHE_HOME")
    if xdg:
        return Path(xdg) / APP_NAME
    return Path.home() / ".cache" / APP_NAME


def ensure_cache_dir(path: Path):
    path.mkdir(parents=True, exist_ok=True)
    try:
        # make it reasonably private
        path.chmod(0o700)
    except Exception:
        # chmod may fail on some filesystems (e.g., Windows), ignore
        pass


def main():
    load_dotenv()

    file_dir = os.path.dirname(os.path.abspath(__file__))
    os.chdir(file_dir)

    parser = argparse.ArgumentParser(description="Get a random item from a Discogs collection folder.")
    parser.add_argument("who", type=str, help="Who's folder to get the item from", choices=["alice","james","both"])
    parser.add_argument('-s', '--singles', action='store_true', help="Include singles")
    parser.add_argument('--notShared', action='store_false', help="Exclude shared folder")
    parser.add_argument('--debug', action='store_true', help="Enable debug mode")
    parser.add_argument('--cache-dir', type=str, help="Override cache directory (overrides DISCOGS_CACHE_DIR env var)")
    args = parser.parse_args()

    # Read token from env (loaded via dotenv if present)
    TOKEN = os.environ.get("DISCOGS_TOKEN") or os.environ.get("DISCOGS_USER_TOKEN") or os.environ.get("TOKEN")
    if not TOKEN:
        print("Error: Discogs token not found. Set DISCOGS_TOKEN or TOKEN in environment.")
        sys.exit(2)

    # determine cache dir (CLI arg > env > default)
    if args.cache_dir:
        cache_dir = Path(args.cache_dir).expanduser()
    else:
        cache_dir = default_cache_dir()
    ensure_cache_dir(cache_dir)

    # install requests_cache using sqlite file in cache dir
    cache_path = str(cache_dir / "discogs_cache")
    # requests_cache will append extension for sqlite backend if needed
    requests_cache.install_cache(cache_path, backend='sqlite', expire_after=86400)

    d = discogs_client.Client('ExampleApplication/0.1', user_token=TOKEN)

    user = d.identity()

    if args.debug:
        print(f"User: {user.username}")
        print("Folders:")
        for folder in user.collection_folders:
            print(f"{folder.id}: {folder.name} ({folder.count} items)")


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