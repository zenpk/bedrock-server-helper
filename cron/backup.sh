#!/bin/bash

. /home/zenpk/mc/vars.sh

tmux send-keys -t $session_panel "say The server will start backup in 5 minutes" Enter
sleep 5m

pid_bedrock_server=$(pidof bedrock_server)
kill -2 $pid_bedrock_server
tmux send-keys -t $session_panel "cd /home/zenpk/mc" Enter
tmux send-keys -t $session_panel ". ./vars.sh" Enter
tmux send-keys -t $session_panel "cd ./backups" Enter

# if already exists, delete it
if_else_block=$(cat << 'EOF'
if [ -d "$date_now" ]; then
  rm -r "$date_now"
  echo "backup $date_now deleted."
fi
EOF
)

tmux send-keys -t $session_panel "$if_else_block" Enter
tmux send-keys -t $session_panel "mkdir $date_now" Enter

# delete the backup older than 7 days
if_else_block_7=$(cat << 'EOF'
if [ -d "$date_7" ]; then
  rm -r "$date_7"
  echo "backup $date_7 deleted."
else
  echo "backup $date_7 does not exist."
fi
EOF
)

tmux send-keys -t $session_panel "$if_else_block_7" Enter
tmux send-keys -t $session_panel "cd .." Enter
tmux send-keys -t $session_panel "cp -r \"./bedrock-server-$THIS_VERSION/worlds/newworld\" \"./backups/$date_now\"" Enter
tmux send-keys -t $session_panel "cd \"./bedrock-server-$THIS_VERSION\"" Enter
echo "$date_now backup finished"
tmux send-keys -t $session_panel ". ./start.sh" Enter
