# Reputation

| Feature              | Supervisor(s)   |
|:-------------------- | ---------------:|
| **``reputation``**   | @ThisIsIPat     |
| ``question-stalker`` | - |
| ``role-stalker``     | @FabianTerhorst |
| documentation        | @ThisIsIPat     |
| testing              | - |

The **Reputation** module describes the network-wide reputation and its ecosystem.

## Commands

### reputation

#### Description

_Shows a users reputation._

Returns ones own reputation unless another user is specified.

#### Usage:
```
reputation [user]
```

#### Example:
```
Mod#0001 @ node1#channel1 > "reputation @User#1234"
ALL      @ node1#channel1 < "12"
Mod#0001 @ node1#channel1 > "reputation"
ALL      @ node1#channel1 < "1.337"
```

<sub><sup>_Examples provide a simplified version of a theoretical conversation. Real conversations may differ heavily due to the available formatting in embed messages._</sup></sub>

#### Aliases:
* reputation
* rep
* r


## Listeners

### question-stalker

#### Description

Watches over dedicated _#ask_ channels within categories defining the topic to manage questions.

F.e., there might be a "Kotlin" category containing an _#ask_ channel.
If a question is asked in that channel, the bot creates another text channel
in that category.
In this channel, it creates an Embed version of the message
and reacts to it (so that others can as well.)

After this, the bot continues to watch over the questions.
Once they are answered, it takes care of adding the reputation.


#### Example:
```
TODO("Example scenario")
```

<sub><sup>_Examples provide a simplified version of a theoretical conversation. Real conversations may differ heavily due to the available formatting in embed messages._</sup></sub>

---

### role-stalker

#### Description

Watches over the reputation of users constantly to ensure they get the role defined by their reputation.