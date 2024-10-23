# Wikiracer

Find a path between two wikipedia articles by doing a bidirectional breadth-first search.

# Usage

```
go run racer.go -start=[first article] -end=[last article]
```

## Example

```
go run racer.go -start="https://en.wikipedia.org/wiki/Battle_of_Cr%C3%A9cy" -end="https://en.wikipedia.org/wiki/Wehrmacht"
```

### Result:

```
Path found:
https://en.wikipedia.org/wiki/Battle_of_Crécy
https://en.wikipedia.org/wiki/Battle of Sluys
https://en.wikipedia.org/wiki/Channel Islands
https://en.wikipedia.org/wiki/Wehrmacht
```
