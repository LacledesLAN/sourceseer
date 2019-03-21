# Sourceseer

[![Go Report Card](https://goreportcard.com/badge/github.com/LacledesLAN/sourceseer)](https://goreportcard.com/report/github.com/LacledesLAN/sourceseer)

**This project is still an alpha-version and not yet recommended for production scenarios**

[Sourceseer](https://github.com/LacledesLAN/sourceseer) is a process ["wrapper"](https://en.wikipedia.org/wiki/Wrapper_library) for [source dedicated servers](https://developer.valvesoftware.com/wiki/Source_Dedicated_Server). By listening to `SRCDS`'s output in real-time we can monitor its game state, selectively push log messages to external systems (such as Slack and Discord), and automate game flow by sending commands in response to server events.

## Motivation

At [Laclede's LAN](https://lacledeslan.com/) we want fault-tolerant and easily distributable mechanisms to monitor our game servers. We also want to automate as much as possible, particularly in high-demand tournaments with tight time schedules such as Counter-Strike: Global Offensive.

When feasible we run [our game servers](https://github.com/LacledesLAN/README.1ST/tree/master/GameServers) in [Docker containers](https://github.com/LacledesLAN/README.1ST/blob/master/GameServers/DockerAndGameServers.md). By writing sourceseer in golang we can compile it to a single binary and include it inside of our Docker images.

After experimenting and using several existing mechanisms we came to the conclusion it was time for us to develop our own from the ground up. To understand why we came to this conclusion let's review the two most-popular tools in widespread usage for managing CSGO tournament servers:

### Why not `eBot`?

[eBot](https://github.com/deStrO/eBot-CSGO) depends upon the CSGO servers streaming data to and receiving commands from a remote node on the network. While we strive for 100% uptime at our events the dynamic and transient nature of LAN parties greatly increases the risk of unexpected outages. To lose a few log messages would be annoying. But to have automated events get triggered from outdated information or to not happen due to a dropped connection would be unacceptable.

### Why not `WarMod [BFG]`?

[WarMod [BFG]](https://forums.alliedmods.net/showthread.php?t=225474) is a plugin for [sourcemod](https://www.sourcemod.net/) which is itself an addon for [Metamod:Source](https://wiki.alliedmods.net/Metamod:Source). Whenever a source game receives a large-enough update this dependency chain can break rendering WarMod non-functional. While these updates are infrequent and while the community is great to quickly to roll out fixed builds (usually within two weeks) a single, poorly timed CSGO update prior to one of our events could risk our tournament's viability.

Additional concerns are that WarMod is no longer actively maintained, it overwrites/drops some cvar values set from the command line (notably `mp_teamname_1` and `mp_teamname_2`), un-reproducible bugs exist (most alarmingly a bug that causes messages to spam and crash clients), and multiple issues exist that cause unexpected issues after a `changelevel` command has been issued (forcing us to use a new server for each map in a best-of-*n* brackets).

We discussed forking WarMod into our own project but this effort would require investing our time into a limited language ([SourcePawn](https://wiki.alliedmods.net/Introduction_to_SourcePawn_1.7)), wouldn't address the at-risk dependency chain, and ultimately wouldn't be the most effective technology to use for our longer-term automation goals.

### Why Go(lang) was Chosen

When choosing a language for `sourceseer` our critical requirement was to able to compile native-binaries to be added directly to our Docker image without needing to include additional required dependencies. Additionally features we wanted in the language were static types, garbage collection, and a low memory footprint.

Not only did [Go](https://github.com/golang/go) meet all of these requirement but its easy to learn nature and native support for the [CSP model](https://en.wikipedia.org/wiki/Communicating_sequential_processes) made it stand out.

## Overview

`Sourceseer` creates a child-process instance of `SRCDS` grabbing exclusive access to its `standard output`, `standard error`, and `standard in` streams. `Sourcseer` determines (and maintains) the state of the game server by observing all output from the `SRCDS` process and is automate changes by to sending `SRCDS` commands when certain conditions are met.

## Common Definitions

TODO

* Affiliation
* Match = map
* Set = set of matches
