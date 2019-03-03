# Sourceseer

[![Go Report Card](https://goreportcard.com/badge/github.com/LacledesLAN/sourceseer)](https://goreportcard.com/report/github.com/LacledesLAN/sourceseer)

[Sourceseer](https://github.com/LacledesLAN/sourceseer) is an ["wrapper"](https://en.wikipedia.org/wiki/Wrapper_library) around an executing [source dedicated server](https://developer.valvesoftware.com/wiki/Source_Dedicated_Server) instance allowing us to monitor and manipulating its state.

**We consider this project an unproven experiment at this time.**

## Why

We run all of our [game servers in Docker containers](https://github.com/LacledesLAN/README.1ST/blob/master/GameServers/DockerAndGameServers.md). By using golang we can compile sourceseer to a single binary (for either linux or windows) and include it inside of our Docker images.

To understand the reason for the efforts behind this new project let's review the two most-popular tools in widespread usage for managing csgo tournament servers:

### WarMod [BFG]

[WarMod [BFG]](https://forums.alliedmods.net/showthread.php?t=225474) is "*designed to be used for competitive matches, and provides the flexibility to suit any form of competition, including large tournaments down to clan matches.*"  It is an addon for [sourcemod](https://www.sourcemod.net/), which is itself an addon for [Metamod:Source](https://wiki.alliedmods.net/Metamod:Source), which is an interceptor that sits between the Source engine and a subsequent game providing APIs.

Unfortunately when a game receives a large update its compatibility with MetaMod:Source and/or Sourcemod can be broken. While the community is quick to roll out new builds (often within 2 weeks) a poorly timed CSGO update prior to one of our charity LANs could risk our CSGO tournament's viability. While the probability of this kind of unfortunate timing may may be low the consequences are great enough to pursue this effort.  Additionally the internal state of SourceMod addons reset when the engine reloads for events such as `changelevel` limiting our ability to automate a tournament across multiple maps (such as a best-of-three).

### eBot

[eBot](https://github.com/deStrO/eBot-CSGO) is a is a full managed server-bot written in PHP and nodeJS. eBot features easy match creation and tons of player and match statistics.

Our reservations for using eBot stem from its orchestration requirements. While not overly technical or challenging they add points of failure and potentials for discrepancies. Under normal circumstances these risks are more minimal, but given the nature of our events we want our servers to be 100% independent of as many external dependencies. We want our servers to continue operating under as many infrastructure issues as possible; it's one thing to push logs to a remote node as a value-add but entirely different for the integrity of all game server's to depend on streaming/receiving real-time data to/from a single, remote node.

## How

## Definition

Match = map
Set = set of matches
