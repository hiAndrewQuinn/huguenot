# huguenot
One-command translation for Hugo files.

One day in fall 2023, I was [hacking away at a news site archive I was working on](https://github.com/hiAndrewQuinn/selkouutiset-scrape-cleaned) and had a fun thought: What if I translated these articles into _every_ language the Google Translate API supports? I could just run a script and have a bunch of translated articles! The next day I woke up with a bill of almost 200 dollars.

So birthed Project Huguenot, an attempt to distil down the shell scripts I had written during that hacking session to make a tool that takes a list of languages and a (possibly-recursive) stack of Markdown files, and translates them all in-place. The ideal was to make it as easy as possible for *returning* users to use -- that is, running `huguenot` in the root directory should be all you need.

