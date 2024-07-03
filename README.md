# What is this?

This is a simple telegram reverse shell bot. It helps managing your remote machines via telegram.

# How to use?

I included a simple [unit file](systemd/tgsh.service) for systemd. You can use it to run the bot as a service. Just copy the file to `/etc/systemd/system/` and run `systemctl enable --now tgsh.service`.

Do not forget to specify your bot token and your telegram id. You can do this by setting the environment variables `TGSH_TOKEN` and `TGSH_USER` respectively or by passing them as arguments to the bot using the `-t` and `-u` flags.

Only the user with the id specified in the `TGSH_USER` environment variable or passed as an argument to the bot can execute commands.

```bash
% tgsh --help 
Usage:
  tgsh [OPTIONS]

Application Options:
  -t, --token= Telegram Bot Token [$TGSH_TOKEN]
  -s, --shell= Shell to execute commands (default: /bin/bash) [$TGSH_SHELL]
  -u, --user=  Telegram User ID allowed to execute commands [$TGSH_USER]
  -d, --debug  Enable debug mode [$TGSH_DEBUG]

Help Options:
  -h, --help   Show this help message
```

# Why did you create this?

I have a couple of servers running in different places and I wanted to have a simple way to manage them without having to expose them to the internet.