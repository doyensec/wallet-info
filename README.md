# wallet-info

## Overview

This repo contains a sample web service which simplifies access to information about Ethereum Dapps. By utilizing tools such as [whois](https://en.wikipedia.org/wiki/WHOIS) lookups and [Etherscan's API](https://docs.etherscan.io/), the latest release of the service provides access to the following information:

- Registration information about the Dapp's domain
- Information about a transaction's target (`to`) address
- Information about the Dapp's domain's TLS certificate

In addition to that, the service implements a mechanizm for verifying the association between the Dapp's domain and it's smart contracts using a `dapp_file`. You can read more information about this mechanism [here](https://blog.doyensec.com/2023/03/28/wallet-info.html).

The service is meant to be used by [wallet](https://en.wikipedia.org/wiki/Cryptocurrency_wallet) applications, providing users information about the Dapp they're interacting with, and information about the transaction they're about to confirm.

More information about the rationale for this project can be found [here](https://blog.doyensec.com/2023/03/28/wallet-info.html).

# Setup

To run the service, first create a `config.yml` file (see `config.yml.example` for a sample file structure), an set your Etherescan API `endpoint` and `api_key`. Afterwards, start the service by running:

```bash
go run main.go
```

## Credits

This tool has been created by Viktor Chuchurski of [Doyensec LLC](https://www.doyensec.com) during our [25% research time](https://doyensec.com/careers.html). 

![alt text](https://doyensec.com/images/logo.svg "Doyensec Logo")
