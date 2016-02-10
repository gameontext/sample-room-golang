#!/bin/bash
# This is the startup script for the golang GameOn! room.
#
# The following environment variables are required:
#   CONTAINER_IP  - The public IP address to which your container is bound.
#                   (This is needed for the websocket callback.)
#   GAMEON_ID     - The ID given in Game On!
#   GAMEON_SECRET - The shared secret given in Game On!
#   ROOM_NAME     - Your name for the room.
#
# The following environment variables are optional:
#   GAMEON_ADDR   - The game server address, defaults to game-on.org.
#   GAMEON_PORT   - Our external port, defaults to 3000.
#                   (This is needed for the websocket callback.)
#   GAMEON_DEBUG  - Any non-empty value turns on debug output
#   GAMEON_TIMESHIFT - Use this to adjust our timestamps to match the
#                      GameOn! server in cases where the time is skewed
#                      relative to our room and the server. This is
#                      expressed in milliseconds. The default is 0.
#
#   The following two variables allow for a delay in network connectivity.
#
#   GAMEON_REG_RETRIES - Number of initial registration attempts.
#                        Defaults to 5
#   GAMEON_REG_SECONDS_BETWEEN - Number of seconds between registration
#                                attempts. Defaults to 30

export ROOM_BINARY=$GOPATH/bin/gameon-room-go
# Print an error message ($1) and exit.
function fatal
{
    echo "FAIL: FATAL ERROR, $1"
    exit 1
}

# $1 is var name
# $2 is var content
function assert_var_set
{
    if [ -z "$2" ] ; then
        fatal "$1 was not set."
    fi
}

# Make sure the required env vars are defined with non-empty values
assert_var_set CONTAINER_IP $CONTAINER_IP
assert_var_set GAMEON_ID $GAMEON_ID
assert_var_set GAMEON_SECRET $GAMEON_SECRET
assert_var_set ROOM_NAME $ROOM_NAME

# Make sure any optional env vars are given default values if they are not defined
export GAMEON_ADDR=${GAMEON_ADDR-game-on.org}
export GAMEON_PORT=${GAMEON_PORT-3000}
export GAMEON_REG_RETRIES=${GAMEON_REG_RETRIES-10}
export GAMEON_REG_SECONDS_BETWEEN=${GAMEON_REG_SECONDS_BETWEEN-15}
export GAMEON_TIMESHIFT=${GAMEON_TIMESHIFT-0}
if [ -z "$GAMEON_DEBUG" ] ; then
    DEBUG_FLAG=""
else
    DEBUG_FLAG="-d"
fi

$ROOM_BINARY \
  -c $CONTAINER_IP \
  -g $GAMEON_ADDR \
  -cp 3000 \
  -lp $GAMEON_PORT \
  -r $ROOM_NAME \
  -id "$GAMEON_ID" \
  -secret "$GAMEON_SECRET" \
  -retries $GAMEON_REG_RETRIES \
  -between $GAMEON_REG_SECONDS_BETWEEN \
  -ts $GAMEON_TIMESHIFT \
  $DEBUG_FLAG
