import sys
import json
import requests
import argparse
from concurrent.futures import ThreadPoolExecutor

# Parse the command-line arguments
parser = argparse.ArgumentParser()
parser.add_argument('--threads', "-t", type=int, default=50, help='Number of threads to use')
args = parser.parse_args()

# Read the domain from sys.stdin
domains = [line.strip() for line in sys.stdin.readlines()]

# Define a function to process a single domain
def process_domain(domain):
    site = f"https://otx.alienvault.com/api/v1/indicators/domain/{domain}/url_list?limit=500"

    # Set the headers
    headers = {
        "Content-Type": "application/json",
        "User-Agent" : "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
    }

    # Send the GET request and get the response
    response = requests.get(site, headers=headers)

    # Parse the response as JSON
    data = response.json()

    # Extract the page and task URLs from the results
    otx_urls = [result['url'] for result in data['url_list']]

    # Merge the page and task URLs and sort them
    urls = list(set(otx_urls))
    urls.sort()

    # Print the URLs one per line
    for url in urls:
        print(url)

# Create a ThreadPoolExecutor with the specified number of threads
with ThreadPoolExecutor(max_workers=args.threads) as executor:
    # Submit the tasks to the executor
    for domain in domains:
        executor.submit(process_domain, domain)
