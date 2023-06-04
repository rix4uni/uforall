import requests
import json
import sys

def get_urls(domain):
    url = f"https://otx.alienvault.com/api/v1/indicators/domain/{domain}/url_list"
    response = requests.get(url, params={"limit": 500, "page": 1})
    
    try:
        response.raise_for_status()
        data = response.json()
        total_urls = data["full_size"]
        total_pages = total_urls // 500 + 1

        for i in range(1, total_pages + 1):
            response = requests.get(url, params={"limit": 500, "page": i})
            response.raise_for_status()
            data = response.json()
            url_list = data["url_list"]
            for url_data in url_list:
                print(url_data["url"])
    
    except (requests.HTTPError, requests.ConnectionError, json.JSONDecodeError) as e:
        pass

if __name__ == "__main__":
    for line in sys.stdin:
        domain = line.strip()
        get_urls(domain)