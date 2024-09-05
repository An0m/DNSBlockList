# DNSBlockList

Just another amalgamation of different block lists, aimed at dns filtering

You can find a raw [here](https://raw.githubusercontent.com/An0m/DNSBlockList/main/output/all.txt), updated... whenever, or you generate it on your own by running main.py

This will also chunk the list into different csv files, in case you're crazy enough to load it as different lists in cloudflare zero trust

Because it's getting kinda late, and I for sure am

---

Feel free to add any other hostnames to the manual.txt file before generating

PS: AdGuard and hostname lists SHOULD be supported, for the rest, it's on you

The script ignores any list starting with "!" and "#"s, and runs a regex match to search for the domain

Wild cards are not supported, as the script just defaults to the parent domain (Cloudflare rules, not mine)

---

If a second level domain is inside the list, any subdomains will be skipped.

You can also add any domain to suffix.txt to exclude any subdomains from the final list
(you cannot exclude subdomains OF subdomains in any other way)

