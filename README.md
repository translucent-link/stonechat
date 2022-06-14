# stonechat
Chainlink External Adapter for Google Sheets

## Introduction

Google Sheets is more than just a simple spreadsheet. Is a proving ground for new business ideas, tokenomics, ideation and prototyping.
Stonechat is a Chainlink Adapter that helps node operators build a bridge to Google Sheets, so that web3 development teams can integrate data held in Google Sheets.

Web3 teams can combine Chainlink Oracles and Google Sheets together to rapidly prototype new solutions by generating the right data in Google Sheets and using the adapter to bring the data on-chain.

## Ingredients you will need

A Google Sheet with some meaningfull data on it


## Getting started

1. Clone the Stonechat repo and download the repo.
2. Install the go-dependencies

        go mod tidy

3. Run the server locally

        go run *.go

You should be able to access the following URLs:
- http://localhost:8080/cell - Main endpoint
- http://localhost:8080/health - Simple health=OK endpoint
- http://localhost:8080/metrics - Prometheus-based metrics endpoint

## Building & Shipping

To dockerize and upload the Docker container image run the build script:

    bin/build.sh 0.0.1

The above command will generate a Linux/AMD64-compatible Docker container image and upload it to Docker.io's Hub (it assumes you're logged in with `docker login`). The `0.0.1` parameter above is used to tag the container image.

NOTE: You should customise the `translucentlink/stonechat` references in the `build.sh` to make them work with your Docker credentials, e.g. replace them with `my-acme-corp/price-feed`.


## Using Translucent's pre-built Docker container image

You don't need to build this repo from scratch and deploy it, you can use our container image.

## Deployment

How you deploy the container is up to you (Docker, AWS, Kubernetes, etc.) but nothing is quite as fast & convenient as using [Fly.io](https://fly.io/)

    flyctl launch --image=translucentlink/stonechat:0.0.1

If you haven't got the `flyctl` command installed, check out their [2-minute intro](https://fly.io/docs/getting-started/installing-flyctl/) on installing and logging in.

To access your deployed external adapter

    flyctl open

To deploy an update to your external adapter

    bin/build.sh 0.0.2
    flyctl deploy --image=translucentlink/stonechat:0.0.2

The example container is deployed at https://ethusd-example.fly.dev/.

## Support & Help

Feel free open a [Github Issue](https://github.com/translucent-link/stonechat/issues) or come find us in the [Translucent Discord](https://discord.gg/RgxXeGuz).
