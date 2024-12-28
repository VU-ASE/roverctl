# Usage

To use, simply run `roverctl` in your shell of choice after installation. You will be greeted with a screen of available options. It is recommended to first set up a connection with a `roverd` instance, but you can already get started writing a service when your Rover is offline.

## Shorthands

### Connect to `roverd`
```bash
roverctl connect
```

Much `roverctl` functionality is based on the public REST API exposed by the `roverd` process, which runs on the Rover. To *connect to a Roverd* you are actually connecting to a *`roverd` instance* instance running on that Rover. Before connecting, make sure that your Rover is powered on and both you and the Rover are connected to the ASE labs WiFi network.

Upon connecting you need to specify the index of your Rover and its username and password. The username and password combination are different than the credentials you might use for SSH. `roverctl` will find the correct IP address for you, but if necessary you can override the IP or hostname manually.

### Initialize a new service
```bash
roverctl service init
```

To quickly get started with the right tools and example code for a service that can run on the Rover, you can initialize a service in your current working directory using `roverctl`. After specifying your service name, the author, source and versioning, you can select which language you want to develop your service in. If there already is a *service.yaml* file present in your current working directory, no new services will be initialized.

### Synchronize/upload an existing service
```bash
roverctl service sync
```

When developing, you want to make sure that the service you are modifying is uploaded to the Rover as well so that it can be enabled in the execution pipeline. `roverctl` will monitor file changes and automatically upload changes to the Rover. You can keep the synchronization process running during development, it will then act as "auto-save".

### View pipeline
```bash
roverctl pipeline
```

Quickly access your pipeline overview, configure, start or stop it. 

### View `roverctl` and `roverd` info
```bash
roverctl info
```

View diagnostics and versioning information about your current `roverctl` installation and the `roverd` instance on your Rover (if connected).