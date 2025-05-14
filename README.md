# crtsh-subdomains
Simple tool to find subdomains using crt.sh search

## Description
This tool uses `https://crt.sh` website to find subdomains using domain name search. It parses the received JSON response and extracts the properly formatted subdomains from the `name_value` field in the resulting JSON array.

## Execution
Run the application with the domain name as an argument:
```
crtsh-subdomains example.com
```
It can also be used with the `dig` tool:
```
crtsh-subdomains example.com | while read OUT; do dig $OUT +nostats +nocomments +nocmd; done
```

## Project Support
You can support this project by donating to the following Ethereum wallet:

ethereum:0x0468DcdE81b69b87ea0A546faA6c5aae2F4FE30b

![ethereum:0x0468DcdE81b69b87ea0A546faA6c5aae2F4FE30b](https://github.com/user-attachments/assets/43858c46-3d18-4899-b42b-55df0b2b1eaa)
