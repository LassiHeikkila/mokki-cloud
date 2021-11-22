# Building
When building, run `npm run-script build`.

Re-arrange stuff so everything except `index.html` goes into `static/`

Copy `settings.json` file that you have locally into `build/static/`, or you can do that on server side.

Then `cd` in `build/` and run `zip ../build.zip -r .`

`scp` the zip onto your server and unzip it there next to your `mokki-server` binary, so `index.html` remains in same dir and the rest is in `static/`
