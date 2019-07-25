#!/bin/bash
openssl genrsa -out keys/udpp.rsa 1024
openssl rsa -in keys/udpp.rsa  -pubout >keys/udpp.rsa.pub