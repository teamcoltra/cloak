# Cloak

Cloak is designed to help protect (third party) websites that need additional levels of privacy, it's also just a basic reverse proxy if you need it to be.

Using Cloudflare is great but as soon as someone submits a DMCA request they will turn over your host's information and your contact information. Cloak is designed to not store any information about the users who connect to a cloaked website as well as maintaining limited contact information for the hosts themselves. The limited information is currently stored inform them of any situations that will affect the service as well as passing on required legal information. 



## Prerequisites

* **Go** (1.16 or later recommended)
* **Git** (for cloning the repository)

### Optional
* **Make** if you use the makefile

## Installation Instructions

### Automated Installation on Windows

* Download or clone the Cloak repository to your local machine.
```
git clone https://github.com/teamcoltra/cloak.git
```

```
cd cloak
```

* Locate the `install.bat` script in the repository, right-click on it, and select **Run as administrator**. This script will:
  * Check for Go installation.
  * Build the Cloak binary.
  * Move the binary and necessary resources to `C:\Program Files\Cloak`.

### Automated Installation on Linux/MacOS

* Clone the repository:
```
git clone https://github.com/teamcoltra/cloak.git
```

```
cd cloak
```

* Build the project using the provided Makefile:
  * To build the binary:
```
make build
```

  * To clean build artifacts:

```
make clean
```

  * To install the application:

```
make install
```


### Manual Installation Steps

* Open Command Prompt or Terminal as Administrator (Windows) or use your terminal on Linux/macOS.
* Build Cloak:
```
go build -o build/cloak .
```

* Move the binary to your desired location, e.g., `/usr/local/bin` on Linux/macOS or `C:\Program Files` on Windows.

## Configuration

Cloak uses a configuration file, `cloak.cfg`, and command-line flags to customize its behavior. By default, the configuration file is located at `/etc/cloak/cloak.cfg` on Linux/macOS and `C:\Program Files\Cloak` on Windows.

### Configuration Options

* `webDir`: Directory where the `index.html` file is located. Default is `.`.
* `map`: Path to the domain mappings file. Default is `domains.txt`.
* `logDir`: Directory for logs. Default is `/var/log/cloak`.
* `port`: Port to listen on. Default is `8080`.
* `apiKey`: API key for sensitive operations.
* `dictionary`: Dictionary file for Babble. Default is `dictionary.txt`.

### Example Config File
You can view the config in stuff/cloak.cfg

```

# Directory where the index.html file lives
webDir=.
# Path to the domain mappings file
map=domains.txt
# Directory where logs should go
logDir=/var/log/cloak
# What port to listen to
port=8080
# API KEY - This is used to overwrite the domain map
apiKey=
# Dictionary of words to be used by babble
dictionary=dictionary.txt

```


## Resources Installation

The installation process may also involve copying additional resources:

* `stuff/www` to your web directory.
* `stuff/log` to your log directory.
* Dictionary, domains, and configuration files to their respective locations.

## Additional Notes

* Installation of binary to system paths like `/usr/local/bin` or `C:\Program Files` may require elevated privileges.
* For detailed usage and advanced configurations, refer to the project documentation or source code comments.
