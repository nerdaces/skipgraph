# Skipgraph
Repository for Skipgraph in Golang.  
***This code needs refactoring***

Skipgraph is one of the technologies of structured overay network.  
It supports range queries unlike DHT (Distributed Hash Table).

## Topology
Each node has key and membership vector (MV)
- key: real number
- MV: random string generated from the set of alphabet, and whose length is infinite

Nodes are always ordered and every nodes except for both ends has 2 neighbors at each level  
Neighbors of each node at the level are decided by prefix of MV


## Details
https://en.wikipedia.org/wiki/Skip_graph
