# Random Item From Discogs

This is a simple script for picking a random item from my Discogs collection. 
The script requires an env file with the following, in the same directory.

     DISCOGS_TOKEN = "<YOUR TOKEN>"

You can then run make, and use the script with the following:
     Usage: ./random-discogs-item [who] [options]

     Positional arguments:
     who   Who's folder to get the item from (choices: alice, james, both)
     Options:
     -debug
          Enable debug mode
     -force-update
          Whether to force update the cache
     -not-shared
          Whether to exclude shared items
     -singles
          Whether to include singles in the selection
