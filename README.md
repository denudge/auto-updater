Disclaimer: This project is work-in-progress.

# Auto-Updater

Keep your distributed on-premise software up-to-date

## Overview

Distributing on-premise software is quite easy at-first - just put your executable on on your website
and wait for users to come by. But once installed, keeping distributed software up-to-date is often cumbersome,
as most users do not frequently check vendor websites or repositories for updates. Also, users are afraid of
installing updates as they fear to lose data, tweaked configuration or functionality.

If your software is not distributed over any official channel like app stores or package repositories,
chances are high your clients run a wide spread of outdated versions.
These may use deprecated API calls, have functional bugs and - even worse - exploitable security flaws.

This project aims to ease software self-updating.

## Installation

The updater consists of two parts: A catalog server containing all your available versions
and a local updater client. The catalog server also includes CLI commands to manage the versions.

## Running the catalog server

You can simply run the catalog server per-project, containing only the versions for that
particular software product. That way, it does not create cross-team dependencies.

The server consists of only one executable, which can be build using the `make` command
or even just `go build`.

    $ make build-catalog

    // Next we need to initialize the database
    $ edit ".env" file to set "POSTGRES_DSN"
    $ ./bin/catalog db init
    $ ./bin/catalog db migrate

    // then you can run the HTTP server
    $ ./bin/catalog serve

Use the `GOOS` and `GOARCH` environment variables to cross-compile for a different architecture when needed.
