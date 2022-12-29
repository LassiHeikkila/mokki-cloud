# Building
When building, run `npm run-script build`.

Re-arrange stuff so everything in `build/` goes into `build/static/`.

Copy `settings.json` file that you have locally into `build/static/`, or you can do that on server side later.

Then `cd` in `build/` and run `zip ../build.zip -r .`.

`scp` the zip onto your server and unzip it there next to your `mokki-server` binary, so all the app files are in `static/`.
