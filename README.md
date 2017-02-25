# cloudbleed-check
Check a list of URLs/hosts to see if they might be affected by #cloudbleed

This is a small tool to help you match a list of URLs or hostnames against a
list of potential affected site. The list is fetched from this repository,
https://github.com/pirate/sites-using-cloudflare. The actually file that
is fetched is 
https://raw.githubusercontent.com/pirate/sites-using-cloudflare/master/sorted_unique_cf.txt

**Note!** The list that is scanned contains possible affected sites. Just because you get a hit
in the list it doesn't mean that it's affected. Anyway, when was the last time you changed
your password for that site. It doesn't take that much time to change it, just to be sure.

Please also read the disclaimer here 
https://github.com/pirate/sites-using-cloudflare

## Running

The list of hostnames that you want to test is provided on stdin.

    cat mysites.txt | cloudbleed-check

Will result in a list of site that might be affected.

    pingdom.com
    runnable.com

If you want to check sites from your browser history you can export it to a text file and
then pipe it into cloudbleed-check. The most important timeframe is between 17 and 22 of February,
but the bug has been running since last summer, but not in the same scale. 

## Export history and bookmarks from Chrome

There are several extensions to chrome that can export the browser history and bookmarks.

These two are tests:

* [Export History/Bookmarks to JSON/XLS* from http://www.json-xls.com](https://chrome.google.com/webstore/detail/export-historybookmarks-t/dcoegfodcnjofhjfbhegcgjgapeichlf)
    
    Does not do a pretty print output so pipe the file from this extension via ```json_pp```

```
cat chrome_history.json | json_pp | sed 's/.*"url".*: "\(.*\)",/\1/;tx;d;:x' | cloudbleed-check
cat chrome_bookmarks.json | json_pp | sed 's/.*"url".*: "\(.*\)",/\1/;tx;d;:x' | cloudbleed-check
```

* [History export from Quamilek](https://chrome.google.com/webstore/detail/history-export/lpmoaclacdaofhlijejogfldmgkdlglj)

    cat history.json | sed 's/.*"url".*: "\(.*\)",/\1/;tx;d;:x' | cloudbleed-check
    