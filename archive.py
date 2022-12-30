import sys
import requests
import argparse
import concurrent.futures

def process_domain(domain):
    url = f"http://web.archive.org/cdx/search/cdx?url=*.{domain}/*&output=text&fl=original&collapse=urlkey"

    headers = {
        "Content-Type": "application/json",
        "User-Agent" : "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
    }

    response = requests.request("GET", url, headers=headers)

    response_lines = response.text.split()

    # process the response lines
    for line in response_lines:
        print(line)

def main():
    # parse the command-line arguments
    parser = argparse.ArgumentParser()
    parser.add_argument("--threads", "-t", type=int, default=50, help="Number of threads to use")
    args = parser.parse_args()

    # read input from stdin
    domains = [line.strip() for line in sys.stdin.readlines()]

    # process the input domains using a thread pool
    with concurrent.futures.ThreadPoolExecutor(max_workers=args.threads) as executor:
        futures = [executor.submit(process_domain, domain) for domain in domains]
        concurrent.futures.wait(futures)

if __name__ == "__main__":
    main()