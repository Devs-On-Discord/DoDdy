# Info
    
| Feature          | Supervisor(s) |
|:---------------- | -------------:|
| **``polls``**    | - |
| ``polls``        | - |
| ``announcement`` | - |
| ``welcome``      | - |
| documentation    | @ThisIsIPat   |
| testing          | - |


The **Info** module provides information globally.

## Commands

### poll

#### Description

Opens a poll for all **node**s to vote on.
The votes return to the **hub** anonymously.


## Listeners

### announcement

#### Description

Forwards an announcement written in the dedicated text channel of the **hub** on Discord to all **node** announcement channels.


### welcome

#### Description

Generates a welcome message shown in all **node** welcome channels containing:

- Rules
- Channel categories
- Roles

defined in the **hub**.

```
TODO("Define how exactly that message should be generated in the wiki (f.e. custom roles would need to be specified extra...)")
```


### polls

#### Description

Updates results of started polls to be shown in the hub.