# Random Item From Discogs

This is a simple script for picking a random item from my Discogs collection. 
The script requires an env file with the following, in the same directory.

     TOKEN = "<YOUR TOKEN>"

You can then run make, and use the script with the following:

    usage: random-discogs-item [-h] [-s] [--notShared] {alice,james,both}

    Get a random item from a Discogs collection folder.

    positional arguments:
    {alice,james,both}  Who's folder to get the item from

    options:
    -h, --help          show this help message and exit
    -s, --singles       Include singles
    --notShared         Exclude shared folder

