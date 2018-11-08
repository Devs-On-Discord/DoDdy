# Moderation

| Feature           | Supervisor(s) |
|:----------------- | -------------:|
| **``redirect``**  | - |
| **``warn``**      | - |
| **``ban``**       | @wasdennnoch  |
| ``ban-fwd``       | - |
| documentation     | @wasdennnoch<br>@ThisIsIPat |
| testing           | @wasdennnoch  |

The **Moderation** module helps node staff in their actions
by providing moderation tools with network-wide effects.

## Commands

### redirect

#### Description

_Provides an invite to a server._

Provides an invite to a specified channel of a server.
If only the server, but no channel is specified,
the bot picks the first (top > down)
 channel visible for everyone.

#### Usage:
```
redirect <server>[#<channel>]
```

#### Example:
```
Mod#0001 @ node1#channel1 > "redirect hub#welcome"
ALL      @ node1#channel1 < "https://discord.gg/WUnPvGM"
Mod#0001 @ node1#channel1 > "redirect hub"
ALL      @ node1#channel1 < "https://discord.gg/WUnPvGM"
```

<sub><sup>_Examples provide a simplified version of a theoretical conversation. Real conversations may differ heavily due to the available formatting in embed messages._</sup></sub>

#### Aliases:
* redirect

---

### warn

#### Description

_Warns the specified user._

Adds a warning to a set of accumulating warnings of a user.
Gives moderators an insight on warnings currently in that set.

#### Usage:
```
warn <user> [reason]
```

#### Example:
```
1:37PM
Mod#0001  @ node1#channel1  > "warn @User#1234 Too original username"
User#1234 @ DIRECT MESSAGES < "You have been given a warning by Mod: Too original username"
Mod#0001  @ DIRECT MESSAGES < "User#1234 was successfully warned. The user didn't have any previous warnings."

1:40PM
Mod#0002  @ node2#channel2  > "warn @User#1234 Dislike the profile picture"
User#1234 @ DIRECT MESSAGES < "You have been given a warning by Mod: Dislike the profile picture"
Mod#0002  @ DIRECT MESSAGES < "User#1234 was successfully warned. The user had the following previous warnings:"
Mod#0002  @ DIRECT MESSAGES < "(1:37PM by Mod#0001[node1#channel1]) Too original username"
```

<sub><sup>_Examples provide a simplified version of a theoretical conversation. Real conversations may differ heavily due to the available formatting in embed messages._</sup></sub>

#### Aliases:
* warn

---

### ban

#### Description

_Bans the specified user._

Bans the specified user and provides an option for other nodes to also
ban that user.

#### Usage:
```
ban <user> [reason]
```

#### Example:
```
Mod#0003  @ node3#channel3  > "ban @User#1234 We mods must really hate you"
NODE      @ node3           < BAN USER#1234
User#1234 @ DIRECT MESSAGES < "You have been banned from node3 by Mod: We must really hate you."
User#1234 @ DIRECT MESSAGES < "Please be aware that we have a network-wide ban system, you might automatically get banned from other servers as well."
Mod#0003  @ DIRECT MESSAGES < "User#1234 has been banned from node3."
ALL       @ hub#global-bans < "User#1234 (Mod#0003 @ node3#channel3): We mods must really hate you" + REACTS: [HAMMER]

Mod#0012  @ hub#global-bans > REACT[HAMMER]
Mod#0012  @ DIRECT MESSAGES < "You're a moderator on more than one server: 1: node1, 2: node2. Please reply with an index of the server you want to ban User#1234 from."
Mod#0012  @ DIRECT MESSAGES > "1"
NODE      @ node1           < BAN USER#1234
Mod#0012  @ DIRECT MESSAGES < "User#1234 has been banned from node1."
User#1234 @ DIRECT MESSAGES < "You have been banned from node1 due to the network-wide ban system."
```

<sub><sup>_Examples provide a simplified version of a theoretical conversation. Real conversations may differ heavily due to the available formatting in embed messages._</sup></sub>

#### Aliases:
* ban


## Listeners

### ban-fwd
(also: ban-forward)

#### Description

When a user of a node is banned, the bot provides
an option for other nodes to also ban that user.