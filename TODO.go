package main

// @TODO: Persistancy in population
/*
Have a map map[ID:string] struct { faction, online }
This allows us to avoid lookups and have logic to auto-purge non-offline players

Maybe not.  This could use a lot of memory.  Hmm
*/

// @FEATURE: Add alerts to pop page
// @FEATURE: !myreports
// @FEATURE: report status in !lookup

/* @saithar
Would be great on statsbot to have a line that says:
Consiousness: Clean/suspicious/bad
to indicate report status, clean = no reports, suspicious got reports but not verified yet and bad got verified reports
Were incoorporating this stuff into our interview process for the outfit
*/

// @FEATURE: !blacklist
