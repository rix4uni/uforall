## uforall

**uforall is a fast url crawler this tool crawl all URLs number of different sources**
### Sources 
- [alienvault](https://otx.alienvault.com)
- [WayBackMachine](http://web.archive.org)
- [urlscan](https://urlscan.io)
- [commoncrawl](https://index.commoncrawl.org/)

## Installation
```
go install github.com/rix4uni/uforall@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/uforall/releases/download/v0.0.2/uforall-linux-amd64-0.0.2.tgz
tar -xvzf uforall-linux-amd64-0.0.2.tgz
rm -rf uforall-linux-amd64-0.0.2.tgz
mv uforall ~/go/bin/uforall
```
Or download [binary release](https://github.com/rix4uni/uforall/releases) for your platform.

## Compile from source
```
git clone --depth 1 github.com/rix4uni/uforall.git
cd uforall; go install
```

## Usage
```
Usage of uforall:
  -silent
        silent mode.
  -t string
        Comma-separated list of tools to run: 'otx', 'archive', 'urlscan', 'commoncrawl', or 'all' (default "all")
  -version
        Print the version of the tool and exit.
```

## Usage Examples

Single URL:
```
echo "testphp.vulnweb.com" | uforall
```

Multiple URLs:
```
cat subs.txt | uforall
```

Run specific tools:
```
cat subs.txt | uforall -t otx, urlscan
```

### Use Domain name instead of subdomain list you can save lot of time
```
# cat withoutprotocolsubs.txt
rest.vulnweb.com
testasp.vulnweb.com
testhtml5.vulnweb.com
testaspnet.vulnweb.com
testphp.vulnweb.com
vulnweb.com


# With subdomains
time cat withoutprotocolsubs.txt | uforall -silent -t all 2>/dev/null | unew -el -i -t | wc -l
12953

real    0m30.728s
user    0m0.120s
sys     0m0.013s


# With domain, saved 25 seconds and output is same
time echo "vulnweb.com" | uforall -silent -t all 2>/dev/null | unew -el -i -t | wc -l
12953

real    0m6.447s
user    0m0.060s
sys     0m0.035s
```
