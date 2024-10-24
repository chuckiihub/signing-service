# Signature Service - Coding Challenge

## Introduction

Hello, hope you're doing well :) I have to say that I've actually enjoyed doing the code challenge. In my last job I was working with Java and it was in my previous job that I worked with Go for a year. I have to say that I love the language, so, this has been fun and I would love to keep studying and improve with better Go best practices from here on.

Thanks for the comments and good requirement explanation, those have been quite nice to have.

## Solution

The structure of the project is the following, with some explanations:

```
api/                  Represents the HTTP Transport layer.
-- dto/                 Contain the Request and Response declarations
-- validation/          Contains the logic to setup the validation logic and "a" custom
                        validator.

config/               It encapsulates the configuration of the app and provides a way to
                      fetch binding address and log level from env vars.

crypto/               I've treated this directory as a library. I did some modifications
                      to support the addition of new encryption algorithms.

domain/               Device, Signature and Health structures.

errors/apperrors.go   This is a custom error wrapper so when the error travels up in the
                      call stack, the caller knows a bit more about the error (like
                      http.StatusCode)

persistence/          This is basically a repository pattern. I assume we might want to
                      store stuff in different databases so I did not enforce any DB
                      dependency between them.

service/              So to abstract the domain logic from the Transport Layer I build
                      services that could be used later to add a different transport
                      layer (websockets, gRPC)

static/docs/          Docs of the API, so you can also test it quickly. Assumes service
                      is running on port 8081
```

### Addition of new Algorithms

Let's say that you want to add a new algorithm, you're gonna have to:

- Add a new SignatureAlgorithm constant (`crypto/algorithm.go`)
- Implement the Crypto interface (which uses a Generic KeyPair) (`crypto/crypto.go`)
- Implement the Signer interface (`signer.go`) and the KeyGenerator (`generation.go`) interface. I kept Signer interface as the challenge wanted me to implement the Signer interface. TBH I think it would have been better (or easier to extend) if everything would be done in the Crypto interface specified above.
- Watch out with the "NewFromSignatureAlgorithm(algorithm)" (constructors) on Signer, KeyGenerator and Crypto.

The validation logic will grab the new algorithm and use it to validate requests, no need to touch it. The Services will use the new Crypto implementations automatically as well.

## Extras

### Swagger docs

You can access the swagger documentation using the path `/api/v0/docs` endpoint.

### Health endpoint

The health endpoint will call each of the services' CheckHealth method and gather all the information in the Health response. This way tracking what's the problem with the service 
should be easier.

### Verify method

To check that I was doing everything alright I've implemented a verify endpoint that answers 200 if the signature is valid and 429 (I'm a Teapot) if the signature is not valid.
Improvement on the response could be done :) 

### Makefile

There's a simple Makefile where you can run the tests, compile, generate the docs.

### CI/CD

There's a simple CI/CD pipeline that runs the tests.

## Things to note

### Signing data and concurrency

The requirement for the Signing functionality only specified that the Device.SignatureCounter should be handled atomically. I've noteice that we should also lock the device to make a consistent trail of the "lastSignature" field.

### Locking Service

As horizontally scaling is generally needed I did the Locking per device (right before signing) in a separate service. Later on, if we are using an external database, this could be implemented
using Redis or any other service to have a locking functionality across different nodes/pods.

## Things I would have like to have the time to do

### Testing

I am usually in favor of Functional/Endpoint testing so to check how the app reacts to different JSON payloads, testing the error cases and unhappy paths so to make sure that the app always works and to be sure that everything works as expected while modifying and upgrading packages. 

I also like to do a lot of error case / not happy path tests. Due to the time, I've decided to stick with showing knowledge about Mocking, unit testing, etc. As a result, not every functionality is being tested which, I also would have liked to do.
