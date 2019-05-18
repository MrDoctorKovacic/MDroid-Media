# MDroid-Media

[![Build Status](https://travis-ci.org/MrDoctorKovacic/MDroid-Media.svg?branch=master)](https://travis-ci.org/MrDoctorKovacic/MDroid-Media) [![Go Report Card](https://goreportcard.com/badge/github.com/MrDoctorKovacic/MDroid-Media)](https://goreportcard.com/report/github.com/MrDoctorKovacic/MDroid-Media)

This is the media controller broken off from the rest of [MDroid-Core](https://github.com/MrDoctorKovacic/MDroid-Core). It handles some low level dbus interactions for a device using bluetooth media. It uses the same routing functions and works in almost the same way as the main program, but without the rest of the unnecessary controllers.

## Requirements
* GoLang, suggested v1.11 at a minimum ([Raspberry Pi Install](https://gist.github.com/kbeflo/9d981573aad107da6fa7ac0603259b3b)) for the rest of MDroid-Core

## Installation 

```go get github.com/MrDoctorKovacic/MDroid-Media/``` 

## Usage

```MDroid-Media --settings-file ./settings.json``` 

## Configuration 

Config is done in the same way as [MDroid-Core](https://github.com/MrDoctorKovacic/MDroid-Core).