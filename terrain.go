package main

type Terrain int

const (
	Forest Terrain = iota
	Farm
	City
	Laboratory
	Edge
)

var terrainTypes = [5]Terrain{Forest, Farm, City, Laboratory, Edge}

func (t Terrain) toString() string {
	return []string{"Forest", "Farm", "City", "Laboratory", "Edge"}[t]
}