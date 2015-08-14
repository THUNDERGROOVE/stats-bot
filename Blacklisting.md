# Blacklisting

## Commands

- `/report <name> <additional information>`
  This would be usable by "normal" users
  - Example `/report THUNDERGROOVE for being a major scrublord`
- `/blacklist lookup <name>`
  This would print a list of all reports for the given user
  >`[ID] <reportee>: <additional information>`
  The reportee would *only be visible for slack admins(leaders)*
 

- `/blacklist clear <ID> <additional information>`
  This would only be usable by admins.  It would mark the report as resolved.
  Would keep a record of the user who marked it as resolved.

- `/blacklist leader <slackname> <outfit>`
  This would mark the given slack user as the main leader for an outfit
  *Only usable by slack admins*
 
## Leader reporting.
If a reported players outfit has a leader marked with /blacklist leader then
we will send them a message on slack via statsbot with the information.

They can then take action with the member or see what actually happened.

## Integration with !lookup&/lookup
If a user has been reported, it will display something in the output of our
lookup commands - @saithar

## Implementation
- Storage: PostgreSQL

  Using `docker run -v /host/path:/container/path db`

  Two tables:
  
  - Leaders {leadername | outfitname }
  - Reports {playername | reportee | additionalinfo | cleared}

  Use Gorm?
