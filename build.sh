#!/bin/bash
go build td.go
cp td ~/bin/

# Generate bash autocomplete
cp bash_autocomplete /usr/local/etc/bash_completion.d/td