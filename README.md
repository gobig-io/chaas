# Chaas


Chaas is a bot for multiple messaging platforms (slack, terminal) that runs targets in Makefiles.
While this is a not the original use case for `make`, it does provide a simple,
powerful, cross platform language for running code. Using the Makefile you can
code your actions in any runtime (go, python, ruby, bash, etc)

## Telling Chaas What to Do
Chaas takes a config file [conf.sample.json](conf.sample.json) that defines the name, location of
the actions directory, and the directions to follow. The key are the directions.
The directions should match the targets in the Makefile. Using the directions,
Chaas knows which words will trigger the given target.

When a new message comes into the channel Chaas first parses it for any directions.
If chaas notices a target specified in the message, first the `options` target
is called to find any variables that should be set for the target. Then, the `intro`
target is called to allow a message to be output. Finally, the matched direction
is called on Makefile target to produce the results and stream them to the
messaging platform


## Building
Chaas recommends go > 1.7.

To build a slack bot
  make bin/slack

To build a terminal bot
  make bin/terminal

To cross compile slack for linux
  make bin/slack-linux-amd64
