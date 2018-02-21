#!/bin/bash

# sleep for some amount of time
sleepy () {
  sleep 0.3
}

# print with formatting
p () {
  echo "$@"
  sleepy
}

# print with formatting to stderr
pe () {
  echo "$@" >&2
  sleepy
}

pe "be careful"
pe "no really be careful"
pe "be especially careful"
pe "i'm advising that you be wary"
pe "warning: watch out for warnings"
pe "everything could be bad"
p "just kidding it's fine"
p "actually things are really good"
p "never been better"
p "forget i said anything"
p "s'all good man"