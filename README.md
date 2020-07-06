# Skipgraph
Skipgraph is one of the technologies of structured overay network.  
It supports range queries unlike DHT (Distributed Hash Table).

## Topology
Each node has key and membership vector (MV)
- key: real number
- MV: random string generated from the set of alphabet, and whose length is infinite

Nodes are always ordered and every nodes except for both ends has 2 neighbors at each level  
Neighbors of each node at the level are decided by prefix of MV

↓ image example

<img width="550" src="https://user-images.githubusercontent.com/65460975/86628262-9d51b000-c004-11ea-849f-503a80049963.png">


## Details
https://en.wikipedia.org/wiki/Skip_graph
