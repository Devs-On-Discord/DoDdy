# Support Module

| Feature         | Supervisor(s) |
|:--------------- | -------------:|
| **``help``**    | @ThisIsIPat   |
| documentation   | @ThisIsIPat   |
| testing         | - |


The **Support** module advises users on how to use the bot.

## Commands

### help

#### Description
_Displays help for commands._

Displays the help strings associated with every command.
If a command is provided as an argument, the help string for that command is displayed instead.

#### Usage:
```
help [command]
```

#### Example:
```
Mod#0001 @ node1#channel1  > "help"
ALL      @ node1#channel1  < "The requested information has been sent."
Mod#0001 @ DIRECT MESSAGES < "Available commands: help, ban, warn" [...]

Mod#0001 @ node1#channel1  > "help ban"
ALL      @ node1#channel1  < "The requested information has been sent."
Mod#0001 @ DIRECT MESSAGES < "ban <user> {reason]"
Mod#0001 @ DIRECT MESSAGES < "Bans the specified user."
Mod#0001 @ DIRECT MESSAGES < "Bans the specified user and provides an option for other nodes to also ban that user."
Mod#0001 @ DIRECT MESSAGES < "Aliases: ban"
```

<sub><sup>_Examples provide a simplified version of a theoretical conversation. Real conversations may differ heavily due to the available formatting in embed messages._</sup></sub>

#### Aliases:
* help
* h