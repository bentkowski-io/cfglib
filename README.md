# cfglib

Implementation of multi level configuration.
- Allows to read config value by key.
- Allows to read config value and unmarshal to map, slice or struct

# sourcing from key=value
For reading configuration in format KEY=VALUE use NewWithPropertyProvider
see example of .env file as a provider of k=v

## sourcing from json
For reading configuration from json file use 
