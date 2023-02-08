import json
import requests
import random
import sys
import argparse
from concurrent.futures import ThreadPoolExecutor

# Read the domain from sys.stdin
domains = [line.strip() for line in sys.stdin.readlines()]

# Define a function to process a domain
def process_domain(domain):
    site = f"https://urlscan.io/api/v1/search/?q=domain:{domain}&size=10000"

    # Set the API endpoint and API key
    api_keys = {'key': ['829e49ac-d524-4464-af9b-53a73a859693', 'a6fc703b-b297-4cdb-a383-c12b211a82ba',]}
    api = random.choice(api_keys['key'])
    # Set the headers
    headers = {
        "Content-Type": "application/json",
        "API-Key": f"{api}",
        "User-Agent" : "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
    }

    # Send the GET request and get the response
    response = requests.get(site, headers=headers)

    # Parse the response as JSON
    data = response.json()

    # Extract the page and task URLs from the results
    page_urls = [result['page']['url'] for result in data['results']]
    task_urls = [result['task']['url'] for result in data['results']]

    # Merge the page and task URLs and sort them
    urls = list(set(page_urls + task_urls))
    urls.sort()

    # Print the URLs one per line
    for url in urls:
        print(url)

# Use argparse to parse the number of threads as an argument
parser = argparse.ArgumentParser()
parser.add_argument("--threads", "-t", type=int, default=50, help="Number of threads to use")
args = parser.parse_args()

# Use a ThreadPoolExecutor to process the domains concurrently
with ThreadPoolExecutor(max_workers=args.threads) as executor:
    executor.map(process_domain, domains)
