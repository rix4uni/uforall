import requests
import json
import regex
import sys
import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed

def get_urls(url, domain):
    param = f"?url=*.{domain}&fl=url&output=json&filter=!=status:404"
    index_url = url + param

    # Send the request to the index
    response = requests.get(index_url)

    # Extract the URLs from the response text
    url_pattern = regex.compile(r'"url": "([^"]+)"')
    urls = url_pattern.findall(response.text)

    return urls

if __name__ == "__main__":
    # Use argparse to specify the number of threads to use
    parser = argparse.ArgumentParser()
    parser.add_argument("--threads", "-t", type=int, default=50, help="Number of threads to use")
    args = parser.parse_args()

    num_threads = args.threads

    # Read the domain name from stdin
    for line in sys.stdin:
        domain = line.strip()

        # Perform an HTTP GET request to the URL
        response = requests.get("https://index.commoncrawl.org/collinfo.json")

        # Parse the JSON data from the response
        data = json.loads(response.text)

        # Create a ThreadPoolExecutor with the specified number of threads
        with ThreadPoolExecutor(max_workers=num_threads) as executor:
            # Create a list of tasks to submit to the executor
            tasks = []
            for item in data:
                url = item['cdx-api']
                task = executor.submit(get_urls, url, domain)
                tasks.append(task)

            # Iterate over the completed tasks and print the results
            for task in as_completed(tasks):
                for url in task.result():
                    print(url)