### Usage

> TODO: change names of modes later

Explicit mode:
```
diagose -- <command>
diagnose -- npm run dev
diagnose -- cmake ..
```

Daemon/Background mode:
```
diagnose start # start diagnose session to listen for errors and automatically give you suggestions

# run all your commands

diagnose end # end diagnose session

# diagnose sessions also end immediately after the exiting the current shell session
```