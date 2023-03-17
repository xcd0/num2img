awk 'BEGIN{ srand('"$RANDOM"'); for (i=0;i<'"$1"';i++){ print( rand() ); } }'
