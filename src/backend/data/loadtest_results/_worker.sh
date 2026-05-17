#!/bin/bash
# Worker: handles full registration lifecycle for one student
# Args: $1=token $2=base_url $3=workshop_id $4=results_file $5=latencies_file
TOKEN="$1"; BASE_URL="$2"; WORKSHOP_ID="$3"; RESULTS="$4"; LATENCIES="$5"
MAX_ATTEMPTS=20; POLL_INTERVAL=3; ATTEMPT=0
START_NS=$(date +%s%N)

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
    ATTEMPT=$((ATTEMPT + 1))
    
    RESP=$(curl -s -w "\n%{http_code}" \
        -X POST "${BASE_URL}/api/v1/registrations" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "{\"workshop_id\":\"${WORKSHOP_ID}\"}" \
        --max-time 15 2>/dev/null || echo -e "\n000")
    
    CODE=$(echo "$RESP" | tail -1)
    
    case "$CODE" in
        202|200)
            END_NS=$(date +%s%N)
            MS=$(( (END_NS - START_NS) / 1000000 ))
            echo "${CODE},${MS}" >> "$RESULTS"
            echo "$MS" >> "$LATENCIES"
            exit 0
            ;;
        400|409)
            END_NS=$(date +%s%N)
            MS=$(( (END_NS - START_NS) / 1000000 ))
            echo "400,${MS}" >> "$RESULTS"
            echo "$MS" >> "$LATENCIES"
            exit 0
            ;;
        429)
            # Poll waiting room until GRANTED
            while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
                sleep $POLL_INTERVAL
                ATTEMPT=$((ATTEMPT + 1))
                
                POLL=$(curl -s \
                    "${BASE_URL}/api/v1/registrations/waiting-room/${WORKSHOP_ID}" \
                    -H "Authorization: Bearer $TOKEN" \
                    --max-time 10 2>/dev/null || echo "")
                
                STATUS=$(echo "$POLL" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    print(d.get('data',{}).get('status_text','UNKNOWN'))
except:
    print('ERROR')
" 2>/dev/null || echo "ERROR")
                
                if [ "$STATUS" = "GRANTED" ] || [ "$STATUS" = "ALREADY_ACTIVE" ]; then
                    break  # Retry registration
                fi
            done
            ;;
        *)
            sleep 1
            ;;
    esac
done

# Timeout
END_NS=$(date +%s%N)
MS=$(( (END_NS - START_NS) / 1000000 ))
echo "TIMEOUT,${MS}" >> "$RESULTS"
echo "$MS" >> "$LATENCIES"
exit 0
