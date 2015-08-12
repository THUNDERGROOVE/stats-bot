package main

// @TODO: Persistancy in population
/*
Have a map map[ID:string] struct { faction, online }
This allows us to avoid lookups and have logic to auto-purge non-offline players

Maybe not.  This could use a lot of memory.  Hmm
*/

// @TODO: Convert !lookup & !lookupeu to /lookup
// @TODO: Talk to other outfit leaders about blacklist implementation

// @FEATURE: !tkreport
// @FEATURE: !blacklist
