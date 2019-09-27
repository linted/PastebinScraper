# How to build GO-lang version
## Dependencies
You should apt-get install these before beginning.
automake, libtool, libssl-dev

1) `git submodule update --init --recursive`

2) `cd go`

3) `./build.sh`

All done! wasn't that simple? 

# Running it
first make a slack config file that contains your slack endpoint url. It should contain only the slack endpoint url with nothing else.

Next, from the go dir, run:
./main -rule YaraRules/certificates.yar -config slack.conf

Then sit back and watch the pastes come flooding in!

# Special thanks
Thanks to @kevtehhermit for some of the great yara rules for scraping pastebin.
https://github.com/kevthehermit/PasteHunter