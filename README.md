# globe
<img src="https://github.com/Kitchen-Kreations/globe/blob/main/img/Globe.PNG?raw=true" data-canonical-src="https://gyazo.com/eb5c5741b6a9a16c692170a41a49c858.png" width="250" height="250" />
<br>
<br>
Port Scanner

## Quick Start
Download the latest release or compile from source
```
$ globe -h
    usage: globe [-h|--help] -i|--ip "<value>" [-p|--ports "<value>"]

                port scanner

    Arguments:

    -h  --help   Print help information
    -i  --ip     What address to scan
    -p  --ports  What ports to scan
```

## Examples
```
# Scan first 1000 ports
globe -i 127.0.0.1
```

```
# Scan all ports
globe -i 127.0.0.1 -p all
```

```
# Scan only certain ports
globe -i 127.0.0.1 -p 80,443
```

```
# Scan a range of ports 80 through 443
glove -i 127.0.0.1 -p 80-443
```
