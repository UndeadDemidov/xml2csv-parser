# xml2csv-parser

The tool parses any set of fields from a list of xml files via xpath expressions in a yaml configuration file,
which is passed by -s <path/file_name> flag.

By default the files are filtered for missing data, but you can configure mandatory output in csv result even if the data is missing.

xml2cv-parser saves found data in csv, utf-8, comma delimiter.

The xpath expression must return the result as the single value of a specific field. Array results and Node results are not supported.

## Prerequisites
 - [Download and install go](https://go.dev/doc/install)
 - Clone this repository

## Build
 - for windows: `make win`
 - for mac: `make mac`
 - for linux: `make linux`

## Use
`xml2csv-parser -p ./xml -s ./parsing.yaml`

### Example of yaml config for parsing:
```yaml
includeFilename: true # filename will be added in last field
set:
- messageType: some_message_type
  columns:
  - name: some_column_name
    xpath: //SomeRoot/SomeNode/SomeElement
    optional: true # file will not be dropped if data is missing - csv field will be blank
  - name: another_column_name
    xpath: //SomeRoot/AnotherNode/AnotherElement[/@SomeAttribute='SomeValue']
- messageType: another_message_type
  columns:
  - name: some_column_name
    xpath: //AnotherRoot/SomeNode/SomeElement
  - name: another_column_name
    xpath: //AnotherRoot/AnotherNode/AnotherElement[/@SomeAttribute='SomeValue']
```
