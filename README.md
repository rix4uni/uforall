# UForAll

UForAll is a fast url crawler this tool crawl all URLs number of different sources
### Sources 
- [alienvault](https://otx.alienvault.com)
- [WayBackMachine](http://web.archive.org)
- [urlscan](https://urlscan.io)
- [commoncrawl](https://index.commoncrawl.org/)

## Installation
```
git clone https://github.com/rix4uni/UForAll.git
cd UForAll
pip3 install -r requirements.txt
chmod +x ./uforall.sh
```
## Setup Api Key `Important` if you not setup api maybe the tool not work properly
```
# https://urlscan.io/user/signup (Paid/Free)
# open urlscan.py add your api keys
```

## Usage
```
OPTIONS:
	-d, --domain        Single Target domain (domain.com)
	-l, --list          Multiple Target domain (interesting_subs.txt)
	-t, --threads       number of threads to use (default 50)
	-h, --help          Help - Show this help

USAGE EXAMPLES:
    ./uforall.sh -d domain.com -t 100
    ./uforall.sh -l interesting_subs.txt -t 100
```

Note: must use `anew` to filter duplicates

Single URL:
```
./uforall.sh -d testphp.vulnweb.com -t 100 | anew
```

Multiple URLs:
```
./uforall.sh -l interesting_subs.txt -t 100 | anew
```

## If you want to use only one service

Single URL:
```
echo testphp.vulnweb.com | python3 archive.py -t 100 | anew
echo testphp.vulnweb.com | python3 otx.py -t 100 | anew
echo testphp.vulnweb.com | python3 urlscan.py -t 100 | anew
echo testphp.vulnweb.com | python3 commoncrawl.py -t 100 | anew
```

Multiple URLs:
```
cat interesting_subs.txt | python3 archive.py -t 100 | anew
cat interesting_subs.txt | python3 otx.py -t 100 | anew
cat interesting_subs.txt | python3 urlscan.py -t 100 | anew
for url in $(cat interesting_subs.txt);do echo "$url" | python3 commoncrawl.py -t 100 | anew;done
```
