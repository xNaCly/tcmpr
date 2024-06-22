#!/usr/bin/bash

if [ ! -f "./pg61.txt" ]; then
    wget "https://www.gutenberg.org/cache/epub/61/pg61.txt"
fi

OG=$((`wc -c ./pg61.txt | awk '{print $1}'`))
echo "[OG   ] $OG bytes" 

gzip -k --force ./pg61.txt
GZIP=$((`wc -c ./pg61.txt.gz | awk '{print $1}'`))
echo "[GZIP] $GZIP bytes $(echo "$GZIP $OG" | awk '{printf "%.2f% \n", ($1*100)/$2}')" 

go build ../cmd/main.go
./main -v1 ./pg61.txt
V1=$((`wc -c ./pg61.txt.tv1 | awk '{print $1}'`))
echo "[v1  ] $V1 bytes $(echo "$V1 $OG" | awk '{printf "%.2f% \n", ($1*100)/$2}')" 

exit
./main ./pg61.txt
V2=$((`wc -c ./pg61.txt.tv2 | awk '{print $1}'`))
echo "[v2  ] $V2 bytes $(echo "$V2 $OG" | awk '{printf "%.2f% \n", ($1*100)/$2}')" 
