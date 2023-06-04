#!/usr/bin/env bash

domain=$1

total_urls=$(curl -s "https://otx.alienvault.com/api/v1/indicators/domain/$domain/url_list?limit=500&page=1" | jq -r '.full_size')
total_pages=$(expr $total_urls / 500 + 1)

for ((i=1; i<=total_pages; i++))
do
  curl -s "https://otx.alienvault.com/api/v1/indicators/domain/$domain/url_list?limit=500&page=$i" | jq -r '.url_list[].url'
done