<p align="center">
  <img src="https://github.com/user-attachments/assets/34a539f0-e8ac-4e8b-b0f9-45c348df25b5" alt="Your Image Description" width="600">
</p>


<h1 align="center"><code>roverctl</code></h1>
<div align="center">
  <a href="https://github.com/VU-ASE/roverctl/releases/latest">Latest release</a>
  <span>&nbsp;&nbsp;•&nbsp;&nbsp;</span>
  <a href="https://ase.vu.nl/docs/category/roverctl">Documentation</a>
  <br />
</div>
<br/>

**`roverctl` is a small terminal user interface (TUI) that provides you with everything you need to develop and work with ASE rover services and pipelines. Interaction is based on an active `roverd` instance, which comes preinstalled on your rover.**

## Installation

### Prebuilt
We provide Linux and macOS users with prebuilt binaries for both amd64 and arm64 systems. You can find and download the latest releases [here](https://github.com/VU-ASE/roverctl/releases/latest). To make installation even easier, you can run our install script, which will detect your system and add the binary to your path automatically:

```bash
curl -fsSL https://raw.githubusercontent.com/VU-ASE/roverctl/main/install.sh | bash
```

### From source

To install the repository from source, you can use our Makefile:
```bash
git clone https://github.com/VU-ASE/roverctl.git
cd roverctl
make build
./bin/roverctl
```
