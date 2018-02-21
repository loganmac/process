#!/bin/bash

# print with formatting - to stderr
p () {
  echo "$@" >&2
  sleep 0.3
}

echo "setting up stage"
sleep 0.3
echo "finding stage..."
sleep 0.3
p "there is no stage"
p "who was supposed to bring drums"
p "these amps have EU plugs"
p "why are the mic stands so short"
p "these are acoustic guitars"
p "the band is an acapella group"
p "they are all crying"
p "lead singer is too afraid to go out"
p "accidentally turned on sprinklers"
p "i smell smoke, how is that even possible?"
p "fire department just showed up"

exit 113
