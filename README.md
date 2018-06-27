# Sitewatch
 
Like Baywatch, but for websites. This is going to be renamed to 'stalker' soon... 

## If...

... You really consider running this in real live, notice that this documentation
is crap. Raise an issue and write some stuff about you. Do not forget to mention
your use case of Sitewatch. We'll get things going for you...

## What it does

Ever refreshed a website ten times while waiting for some expected update in
order to figure out that nothing happend yet?

Let Sitewatch waist its time to do so: give it a couple of websites and sitewatch
checks if their content change (via md5 sum of its content). You'll only get 
notified if changes occure...

## Install

Since currently no binaries are provided you need to compile Sitewatch by hand.
Go makes this easy: 

1) Install Go (https://golang.org/doc/)
2) Compile: `go get -u github.com/unprofession-al/sitewatch`

That's it. 

## Configure

Create a config.yaml file and provide a list of websites you want to have stalked:

```
---
# check interval in seconds
interval: 300
debug: false
sites:
  localhost:
    url: https://example.com
    recipient: me@example.com
    template: example.com has changed
```

## Run

Run this is the directory where your config lives:

```
SITEWATCH_NOTIFIER_SENDGRID="noreply@stalkingbastard.com SG.yG2dlva4R4KO8-ThisIs MySendGridKey" SITEWATCH_NOTIFIER_STDOUT="YES" sitewatch -config /path/to/config/file.yaml
```
