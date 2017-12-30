# Chaas


Chaas is a bot for multiple messaging platforms (slack, terminal) that runs targets in Makefiles.
While this is a not the original use case for `make`, it does provide a simple,
powerful, cross platform language for running code. Using the Makefile you can
code your actions in any runtime (go, python, ruby, bash, etc)

## Telling Chaas What to Do
Chaas takes json file (see [directions.sample.json](directions.sample.json)) that defines the location of the actions directory, and
the directions to follow. The directions should match the targets in the Makefile.
Using the directions, Chaas knows which words will trigger the given target.

When a new message comes into the channel Chaas first parses it for any directions.
If Chaas finds a direction in the message, first the `options` target
is called to find any variables that should be set for the target. Then, the `intro`
target is called to allow a message to be output. Finally, the matched direction
is called on Makefile target to produce the results and stream them to the
messaging platform.

## Example Makefile

  see [example/Makefile](example/Makefile)

## Building
Chaas recommends Go > 1.7.

To build a slack bot

    make bin/slack

To use slack bot

    Usage of ./bin/slack:
    -directions string
      Path to directions.json (default "directions.json")
    -id string
      Slack Bot User ID
    -name string
      Slack Bot Name (default "chaas")
    -token string
      Slack Bot API Key

To build a terminal bot

    make bin/terminal

To cross compile slack for linux

    make bin/slack-linux-amd64
