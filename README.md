# Terms lookup
For peer discovery
- DHT 
- PEX 
- magnet links

# Resources 
- [RoadMapSh](https://roadmap.sh/guides/torrent-client)

# Steps 
## Parsing a .torrent file

Encoded in BENCODE

Strings have lenght prefix

    EX : 4:spam

Integers : between start and end markers 
    EX : 7 -> i7e

Dict and lists : 
    EX : ['spam',7] -> l4:spami7ee
    EX : {"spam" : 7} -> d4:spami7ee

Extract [infohash] and pieceshashes

## Retrieving peers from tracker
-> GEt request to the announce url to retrieve a list of peers 
With params 
- info_hash -> file willing to download
- peer_id -> generated 20byte name id to identify iurselves on the network


## PArsing tracker response

## Download from peers 
### Start a TCP connection

### Complete the handshake 
-> Bittorrent handshake 

