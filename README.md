# Translation tool

[![GoDoc](https://godoc.org/github.com/dakalab/translate?status.svg)](https://godoc.org/github.com/dakalab/translate)
[![Build Status](https://travis-ci.org/dakalab/translate.svg?branch=master)](https://travis-ci.org/dakalab/translate)
[![](https://goreportcard.com/badge/github.com/dakalab/translate)](https://goreportcard.com/report/github.com/dakalab/translate)
[![codecov](https://codecov.io/gh/dakalab/translate/branch/master/graph/badge.svg)](https://codecov.io/gh/dakalab/translate)
[![Docker Pulls](https://img.shields.io/docker/pulls/dakalab/translate.svg)](https://hub.docker.com/r/dakalab/translate)


This is a translation tool for translating from json-format input file into target language by using google cloud translate.

The supported json format is:

```
{
    "key1": "value1",
    "key2": "value2",
    ...
}
```

## Run in docker (recommend)

1) Print usage

```
docker run dakalab/translate -h
```

2) Get available languages:

```
docker run -e GCLOUD_API_KEY=your-key dakalab/translate -l
```

3) Translate

```
docker run -e GCLOUD_API_KEY=your-key dakalab/translate -i "input-json-file" -o "output-json-file" -s source-language -t target-language
```

Below is a simple example which will translate the demo.json into Chinese and output to stdout:

```
docker run -e GCLOUD_API_KEY=your-key -v $PWD/demo.json:/demo.json dakalab/translate -i "/demo.json" -t zh
```

## Run with local golang

### Install

```
go get -u github.com/dakalab/translate
```

### Usage

1) Print usage

```
translate -h
```

2) Get available languages:

```
GCLOUD_API_KEY=your-key translate -l
```

3) Translate

```
GCLOUD_API_KEY=your-key translate -i "input-json-file" -o "output-json-file" -s source-language -t target-language
```

Below is a simple example which will translate the demo.json into Chinese and output to stdout:

```
GCLOUD_API_KEY=your-key translate -i "./demo.json" -t zh
```
